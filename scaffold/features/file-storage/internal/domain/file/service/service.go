package service

import (
	"context"
	"fmt"
	"go_platform_template/internal/domain/file/model"
	"go_platform_template/internal/domain/file/repo"
	"go_platform_template/internal/platform/config"
	"io"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// FileService handles file operations including upload, download, and signed URL generation
// It integrates with MinIO for object storage and the database for metadata storage
type FileService struct {
	minioClient *minio.Client
	bucket      string
	repo        repo.FileRepo
	logger      *zap.SugaredLogger
}

// FileServiceConfig defines the configuration required for initializing FileService
type FileServiceConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
}

// NewFileService creates a new instance of FileService with the provided configuration
// It initializes the MinIO client and ensures the bucket exists
//
// Parameters:
//   - repo: File repository for metadata operations
//   - cfg: MinIO configuration from the main application config
//   - logger: Logger for service operations
//
// Returns:
//   - *FileService: Initialized file service instance
//   - error: Any error encountered during MinIO client initialization or bucket creation
func NewFileService(fileRepo repo.FileRepo, cfg *config.Config, logger *zap.SugaredLogger) (*FileService, error) {
	// Extract MinIO configuration from the main config
	minioCfg := FileServiceConfig{
		Endpoint:        cfg.MinIO.MinioEndpoint,
		AccessKeyID:     cfg.MinIO.MinioAccessKey,
		SecretAccessKey: cfg.MinIO.MinioSecretKey,
		Bucket:          cfg.MinIO.MinioBucket,
		UseSSL:          cfg.MinIO.MinioUseSSL,
	}

	// Initialize MinIO client
	minioClient, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKeyID, minioCfg.SecretAccessKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Verify connection and create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, minioCfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		logger.Infof("Creating MinIO bucket: %s", minioCfg.Bucket)
		err = minioClient.MakeBucket(ctx, minioCfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}

		// Set bucket policy for public read access (adjust based on your requirements)
		bucketPolicy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": [
						"s3:GetObject",
						"s3:PutObject",
						"s3:DeleteObject",
						"s3:ListBucket"
					],
					"Resource": [
						"arn:aws:s3:::%s",
						"arn:aws:s3:::%s/*"
					]
				}
			]
		}`, minioCfg.Bucket, minioCfg.Bucket)

		if err := minioClient.SetBucketPolicy(ctx, minioCfg.Bucket, bucketPolicy); err != nil {
			logger.Warnf("Failed to set bucket policy: %v", err)
		}
	} else {
		logger.Infof("Using existing MinIO bucket: %s", minioCfg.Bucket)
	}

	return &FileService{
		minioClient: minioClient,
		bucket:      minioCfg.Bucket,
		repo:        fileRepo,
		logger:      logger,
	}, nil
}

// Upload handles file upload to MinIO storage and saves metadata to database
//
// Parameters:
//   - userID: ID of the user uploading the file
//   - fType: Type of the file (e.g., image, document, video)
//   - fileReader: Reader interface for the file content
//   - objectName: Unique name for the object in storage
//   - size: Size of the file in bytes
//   - contentType: MIME type of the file
//   - originalName: Original filename as uploaded by the user
//
// Returns:
//   - *model.File: File metadata including generated path and ID
//   - error: Any error encountered during upload or metadata save
func (s *FileService) Upload(userID uuid.UUID, fType model.FileType, fileReader io.Reader, objectName string, size int64, contentType string, originalName string) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Upload file to MinIO
	_, err := s.minioClient.PutObject(ctx, s.bucket, objectName, fileReader, size, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"uploaded-by":   userID.String(),
			"file-type":     string(fType),
			"original-name": originalName,
		},
	})
	if err != nil {
		return nil, err
	}

	// Create file metadata using the enhanced File model
	file := &model.File{
		UserID:       userID,
		Path:         objectName,
		Type:         fType,
		Size:         size,
		MimeType:     contentType,
		OriginalName: originalName,
	}

	// Save metadata to database
	if err := s.repo.SaveFileMeta(ctx, file); err != nil {
		// If database save fails, attempt to clean up the uploaded file
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if cleanupErr := s.minioClient.RemoveObject(cleanupCtx, s.bucket, objectName, minio.RemoveObjectOptions{}); cleanupErr != nil {
			s.logger.Warnf("Failed to cleanup file after metadata save failure: %v", cleanupErr)
		}
		return nil, err
	}

	return file, nil
}

// GetSignedURL generates a pre-signed URL for temporary access to a file
// The signed URL can be used to download the file without requiring authentication
// for the specified duration
//
// Parameters:
//   - objectName: Name of the object in storage
//   - expiry: Duration for which the signed URL should be valid
//
// Returns:
//   - string: Pre-signed URL for accessing the file
//   - error: Any error encountered during URL generation
func (s *FileService) GetSignedURL(objectName string, expiry time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqParams := make(url.Values)
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// Delete removes a file from both MinIO storage and the metadata database
//
// Parameters:
//   - objectName: Name of the object to delete
//
// Returns:
//   - error: Any error encountered during deletion
func (s *FileService) Delete(objectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete from MinIO storage
	err := s.minioClient.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	// Delete metadata from database
	if err := s.repo.DeleteFileMeta(context.Background(), objectName); err != nil {
		s.logger.Warnf("Failed to delete file metadata for %s: %v", objectName, err)
		// Don't return error here as the main storage object was deleted successfully
	}

	return nil
}

// FileExists checks if a file exists in MinIO storage
//
// Parameters:
//   - objectName: Name of the object to check
//
// Returns:
//   - bool: true if file exists, false otherwise
//   - error: Any error encountered during the check
func (s *FileService) FileExists(objectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.minioClient.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *FileService) GetFileByPath(ctx context.Context, objectName string) (*model.File, error) {
	return s.repo.GetFileByPath(ctx, objectName)
}

func (s *FileService) GetFilesByUserID(ctx context.Context, userID string) ([]model.File, error) {
	return s.repo.GetFilesByUserID(ctx, userID)
}
