package dto

import (
	"time"
)

// UploadRequest represents the payload for uploading a file
// swagger:model
type UploadRequest struct {
	// Type of the file to upload
	// Required: true
	// Enum: profile_image,cv
	// Example: profile_image
	Type string `json:"type" form:"type" validate:"required,oneof=profile_image cv"`

	// File to upload
	// Required: true
	// Swagger type: file
	File interface{} `json:"file" form:"file" validate:"required"`
}

// UploadResponse represents the response after successful file upload
// swagger:model
type UploadResponse struct {
	// ID of the uploaded file
	// Example: 550e8400-e29b-41d4-a716-446655440000
	FileID string `json:"file_id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// URL to access the file (signed URL)
	// Example: https://minio.example.com/bucket/path?X-Amz-Algorithm=...
	URL string `json:"url" example:"https://minio.example.com/bucket/path?X-Amz-Algorithm=..."`

	// Path of the file in the storage
	// Example: user-123/profile_20231201-143052.jpg
	Path string `json:"path" example:"user-123/profile_20231201-143052.jpg"`

	// Type of the file
	// Example: profile_image
	Type string `json:"type" example:"profile_image"`

	// Size of the file in bytes
	// Example: 1024576
	Size int64 `json:"size" example:"1024576"`

	// Original name of the file
	// Example: my_profile.jpg
	OriginalName string `json:"original_name" example:"my_profile.jpg"`

	// MIME type of the file
	// Example: image/jpeg
	MimeType string `json:"mime_type" example:"image/jpeg"`

	// UploadedAt is the timestamp when the file was uploaded
	// Example: 2023-12-01T14:30:52Z
	UploadedAt time.Time `json:"uploaded_at"`

	// ExpiresIn is the duration after which the URL expires
	// Example: 15 minutes
	ExpiresIn string `json:"expires_in" example:"15 minutes"`
}

// GetFileResponse represents the response for file access
// swagger:model
type GetFileResponse struct {
	// URL to access the file (signed URL)
	// Example: https://minio.example.com/bucket/path?X-Amz-Algorithm=...
	URL string `json:"url" example:"https://minio.example.com/bucket/path?X-Amz-Algorithm=..."`

	// ExpiresIn is the duration after which the URL expires
	// Example: 15 minutes
	ExpiresIn string `json:"expires_in" example:"15 minutes"`
}

// UserFilesResponse represents the response for listing user files
// swagger:model
type UserFilesResponse struct {
	// Count of files
	// Example: 5
	Count int `json:"count" example:"5"`

	// Files is the list of file metadata
	Files []FileInfo `json:"files"`
}

// FileInfo represents file metadata information
// swagger:model
type FileInfo struct {
	// ID of the file
	// Example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Path of the file in the storage
	// Example: user-123/profile.jpg
	Path string `json:"path" example:"user-123/profile.jpg"`

	// Type of the file
	// Example: profile_image
	Type string `json:"type" example:"profile_image"`

	// Size of the file in bytes
	// Example: 1024576
	Size int64 `json:"size" example:"1024576"`

	// Original name of the file
	// Example: my_profile.jpg
	OriginalName string `json:"original_name" example:"my_profile.jpg"`

	// MIME type of the file
	// Example: image/jpeg
	MimeType string `json:"mime_type" example:"image/jpeg"`

	// UploadedAt is the timestamp when the file was uploaded
	// Example: 2023-12-01T14:30:52Z
	UploadedAt time.Time `json:"uploaded_at"`

	// URL to access the file (signed URL)
	// Example: https://minio.example.com/bucket/path?X-Amz-Algorithm=...
	URL string `json:"url" example:"https://minio.example.com/bucket/path?X-Amz-Algorithm=..."`
}

// DeleteFileRequest represents the payload for deleting a file
// swagger:model
type DeleteFileRequest struct {
	// Path of the file to delete
	// Required: true
	// Example: user-123/profile_20231201-143052.jpg
	Path string `json:"path" validate:"required,min=1"`
}

// SuccessResponse represents a generic success response
// swagger:model
type SuccessResponse struct {
	// Message of the response
	// Example: operation completed successfully
	Message string `json:"message" example:"operation completed successfully"`
}

// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	// Error message
	// Example: error description
	Error string `json:"error" example:"error description"`
}

// FileQueryParams represents query parameters for file operations
// swagger:parameters GetFile DeleteFile
type FileQueryParams struct {
	// Filename path parameter
	// in: path
	// Required: true
	// Example: user-123/profile_20231201-143052.jpg
	Filename string `json:"filename"`
}

// UploadQueryParams represents query parameters for upload operation
// swagger:parameters UploadFile
type UploadQueryParams struct {
	// Type query parameter
	// in: query
	// Required: true
	// Enum: profile_image,cv
	// Example: profile_image
	Type string `json:"type"`
}
