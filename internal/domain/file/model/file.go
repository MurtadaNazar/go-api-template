package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileType represents the category or purpose of the uploaded file
type FileType string

const (
	// Profile image file type
	FileTypeProfileImage FileType = "profile_image"

	// CV/document file type
	FileTypeCV FileType = "cv"
)

// File represents a file stored in the system with metadata
// swagger:model File
type File struct {
	// ID is the unique identifier for the file
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	// UserID is the UUID of the user who owns this file
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_files_user_id" json:"user_id"`

	// Path where the file is stored in the system
	// example: /uploads/profile_images/123e4567-e89b-12d3-a456-426614174000.jpg
	Path string `gorm:"type:varchar(1024);not null;uniqueIndex:idx_files_path" json:"path"`

	// Type categorizes the purpose of the file
	// enum: profile_image,cv
	// example: profile_image
	Type FileType `gorm:"type:varchar(50);not null;index:idx_files_type" json:"type"`

	// Size of the file in bytes
	// example: 2048000
	// minimum: 0
	Size int64 `gorm:"type:bigint;not null;default:0" json:"size"`

	// MimeType indicates the media type of the file
	// example: image/jpeg
	MimeType string `gorm:"type:varchar(255);not null" json:"mime_type"`

	// OriginalName is the filename as uploaded by the user
	// example: my_profile_picture.jpg
	// max length: 512
	OriginalName string `gorm:"type:varchar(512);not null" json:"original_name"`

	// UploadedAt indicates when the file was uploaded
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`

	// UpdatedAt shows when the file metadata was last modified
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook that generates a UUID for the file if not already set
func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}

// TableName specifies the custom table name for the File model
func (File) TableName() string {
	return "files"
}

// GetFileTypeFromMIME returns the appropriate FileType based on MIME type
func GetFileTypeFromMIME(mimeType string) FileType {
	switch {
	case mimeType == "application/pdf" ||
		mimeType == "application/msword" ||
		mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return FileTypeCV
	case mimeType == "image/jpeg" ||
		mimeType == "image/png" ||
		mimeType == "image/gif" ||
		mimeType == "image/webp" ||
		mimeType == "image/svg+xml":
		return FileTypeProfileImage
	default:
		if len(mimeType) >= 5 && mimeType[0:5] == "image" {
			return FileTypeProfileImage
		}
		return FileTypeCV
	}
}
