package service

import (
	"fmt"
	"go_platform_template/internal/domain/file/model"
	"mime"
	"path/filepath"
	"strings"
)

// FileValidationConfig holds configuration for file validation
type FileValidationConfig struct {
	MaxProfileImageSize int64
	MaxCVSize           int64
	AllowedImageTypes   []string
	AllowedCVTypes      []string
}

// DefaultFileValidationConfig returns the default file validation configuration
func DefaultFileValidationConfig() FileValidationConfig {
	return FileValidationConfig{
		MaxProfileImageSize: 5 << 20,  // 5MB
		MaxCVSize:           10 << 20, // 10MB
		AllowedImageTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"image/svg+xml",
		},
		AllowedCVTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.oasis.opendocument.text",
		},
	}
}

// FileValidationRequest represents the input for file validation
type FileValidationRequest struct {
	FileName    string
	FileSize    int64
	ContentType string
	FileType    model.FileType
}

// ValidateFileType validates the uploaded file against the expected file type and constraints
func ValidateFileType(req FileValidationRequest, config FileValidationConfig) error {
	// Determine actual file type from MIME type
	actualType := model.GetFileTypeFromMIME(req.ContentType)

	if req.FileType != actualType {
		return fmt.Errorf("invalid file type for %s. Expected %s file, got %s",
			string(req.FileType), string(req.FileType), string(actualType))
	}

	// Validate file extension matches content type
	if err := validateFileExtension(req.FileName, req.ContentType); err != nil {
		return err
	}

	// Validate against allowed MIME types
	if err := validateAllowedMimeTypes(req.ContentType, req.FileType, config); err != nil {
		return err
	}

	// Validate file size based on type
	switch req.FileType {
	case model.FileTypeProfileImage:
		if req.FileSize > config.MaxProfileImageSize {
			return fmt.Errorf("profile image too large. Maximum size is %dMB", config.MaxProfileImageSize>>20)
		}
	case model.FileTypeCV:
		if req.FileSize > config.MaxCVSize {
			return fmt.Errorf("CV file too large. Maximum size is %dMB", config.MaxCVSize>>20)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	return nil
}

// validateFileExtension checks if the file extension matches the content type
func validateFileExtension(fileName, contentType string) error {
	ext := strings.ToLower(filepath.Ext(fileName))

	if ext == "" {
		return fmt.Errorf("file must have a valid extension")
	}

	// Get expected extensions for the content type
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		// Fallback validation: check common extension-to-MIME mappings
		return validateExtensionAgainstCommonMimes(ext, contentType)
	}

	// Check if the file extension matches any of the expected extensions
	for _, expectedExt := range exts {
		if strings.EqualFold(ext, expectedExt) {
			return nil
		}
	}

	return fmt.Errorf("file extension '%s' does not match content type '%s'", ext, contentType)
}

// validateExtensionAgainstCommonMimes performs fallback validation for common MIME types
func validateExtensionAgainstCommonMimes(ext, contentType string) error {
	// Common mappings for image types
	imageExts := map[string][]string{
		"image/jpeg":    {".jpg", ".jpeg"},
		"image/png":     {".png"},
		"image/gif":     {".gif"},
		"image/webp":    {".webp"},
		"image/svg+xml": {".svg"},
	}

	// Common mappings for document types
	docExts := map[string][]string{
		"application/pdf":    {".pdf"},
		"application/msword": {".doc"},
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {".docx"},
		"application/vnd.oasis.opendocument.text":                                 {".odt"},
	}

	// Combine both mappings
	allExts := make(map[string][]string)
	for k, v := range imageExts {
		allExts[k] = v
	}
	for k, v := range docExts {
		allExts[k] = v
	}

	if expectedExts, ok := allExts[contentType]; ok {
		for _, expectedExt := range expectedExts {
			if strings.EqualFold(ext, expectedExt) {
				return nil
			}
		}
		return fmt.Errorf("file extension '%s' does not match content type '%s'", ext, contentType)
	}

	// Unknown MIME type - be strict and reject it
	return fmt.Errorf("unknown content type '%s'", contentType)
}

// validateAllowedMimeTypes checks if the MIME type is allowed for the file type
func validateAllowedMimeTypes(contentType string, fileType model.FileType, config FileValidationConfig) error {
	var allowedTypes []string

	switch fileType {
	case model.FileTypeProfileImage:
		allowedTypes = config.AllowedImageTypes
	case model.FileTypeCV:
		allowedTypes = config.AllowedCVTypes
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return nil
		}
	}

	return fmt.Errorf("content type '%s' is not allowed for %s files", contentType, fileType)
}

// Add this method to your FileService to use the validation
func (s *FileService) ValidateUpload(fileName string, fileSize int64, contentType string, fileType model.FileType) error {
	config := DefaultFileValidationConfig()
	req := FileValidationRequest{
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
		FileType:    fileType,
	}
	return ValidateFileType(req, config)
}
