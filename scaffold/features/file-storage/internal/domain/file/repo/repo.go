package repo

import (
	"context"
	"go_platform_template/internal/domain/file/model"

	"gorm.io/gorm"
)

type FileRepo interface {
	SaveFileMeta(ctx context.Context, file *model.File) error
	GetFileByID(ctx context.Context, id string) (*model.File, error)
	DeleteFileMeta(ctx context.Context, objectPath string) error
	GetFileByPath(ctx context.Context, objectPath string) (*model.File, error)
	GetFilesByUserID(ctx context.Context, userID string) ([]model.File, error)
}

type fileRepo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) FileRepo {
	return &fileRepo{db: db}
}

func (r *fileRepo) SaveFileMeta(ctx context.Context, file *model.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepo) GetFileByID(ctx context.Context, id string) (*model.File, error) {
	var file model.File
	err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFileMeta deletes file metadata by object path (soft delete)
func (r *fileRepo) DeleteFileMeta(ctx context.Context, objectPath string) error {
	return r.db.WithContext(ctx).Where("path = ?", objectPath).Delete(&model.File{}).Error
}

// GetFileByPath retrieves a file by its object path
func (r *fileRepo) GetFileByPath(ctx context.Context, objectPath string) (*model.File, error) {
	var file model.File
	err := r.db.WithContext(ctx).First(&file, "path = ?", objectPath).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFilesByUserID retrieves all files for a specific user
func (r *fileRepo) GetFilesByUserID(ctx context.Context, userID string) ([]model.File, error) {
	var files []model.File
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
