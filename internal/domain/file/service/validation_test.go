package service

import (
	"go_platform_template/internal/domain/file/model"
	"testing"
)

func TestValidateFileType_Success(t *testing.T) {
	config := DefaultFileValidationConfig()

	tests := []struct {
		name    string
		req     FileValidationRequest
		wantErr bool
	}{
		{
			name: "Valid JPEG image",
			req: FileValidationRequest{
				FileName:    "profile.jpg",
				FileSize:    1024 * 1024, // 1MB
				ContentType: "image/jpeg",
				FileType:    model.FileTypeProfileImage,
			},
			wantErr: false,
		},
		{
			name: "Valid PNG image",
			req: FileValidationRequest{
				FileName:    "photo.png",
				FileSize:    2 * 1024 * 1024, // 2MB
				ContentType: "image/png",
				FileType:    model.FileTypeProfileImage,
			},
			wantErr: false,
		},
		{
			name: "Valid PDF CV",
			req: FileValidationRequest{
				FileName:    "resume.pdf",
				FileSize:    5 * 1024 * 1024, // 5MB
				ContentType: "application/pdf",
				FileType:    model.FileTypeCV,
			},
			wantErr: false,
		},
		{
			name: "Valid DOCX CV",
			req: FileValidationRequest{
				FileName:    "cv.docx",
				FileSize:    3 * 1024 * 1024, // 3MB
				ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				FileType:    model.FileTypeCV,
			},
			wantErr: false,
		},
		{
			name: "Image too large",
			req: FileValidationRequest{
				FileName:    "large.jpg",
				FileSize:    6 * 1024 * 1024, // 6MB (max is 5MB)
				ContentType: "image/jpeg",
				FileType:    model.FileTypeProfileImage,
			},
			wantErr: true,
		},
		{
			name: "CV too large",
			req: FileValidationRequest{
				FileName:    "large_cv.pdf",
				FileSize:    11 * 1024 * 1024, // 11MB (max is 10MB)
				ContentType: "application/pdf",
				FileType:    model.FileTypeCV,
			},
			wantErr: true,
		},
		{
			name: "Wrong MIME type for image file",
			req: FileValidationRequest{
				FileName:    "image.jpg",
				FileSize:    1024 * 1024,
				ContentType: "application/pdf",
				FileType:    model.FileTypeProfileImage,
			},
			wantErr: true,
		},
		{
			name: "Wrong MIME type for CV",
			req: FileValidationRequest{
				FileName:    "cv.jpg",
				FileSize:    1024 * 1024,
				ContentType: "image/jpeg",
				FileType:    model.FileTypeCV,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileType(tt.req, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileExtension_Success(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		contentType string
		wantErr     bool
	}{
		{
			name:        "Valid JPEG extension",
			fileName:    "photo.jpg",
			contentType: "image/jpeg",
			wantErr:     false,
		},
		{
			name:        "Valid JPEG with alternate extension",
			fileName:    "image.jpeg",
			contentType: "image/jpeg",
			wantErr:     false,
		},
		{
			name:        "Valid PNG extension",
			fileName:    "picture.png",
			contentType: "image/png",
			wantErr:     false,
		},
		{
			name:        "Valid PDF extension",
			fileName:    "document.pdf",
			contentType: "application/pdf",
			wantErr:     false,
		},
		{
			name:        "Valid DOCX extension",
			fileName:    "resume.docx",
			contentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			wantErr:     false,
		},
		{
			name:        "Missing extension",
			fileName:    "photo",
			contentType: "image/jpeg",
			wantErr:     true,
		},
		{
			name:        "Wrong extension for MIME type",
			fileName:    "image.pdf",
			contentType: "image/jpeg",
			wantErr:     true,
		},
		{
			name:        "Case-insensitive extension matching",
			fileName:    "photo.JPG",
			contentType: "image/jpeg",
			wantErr:     false,
		},
		{
			name:        "Unknown MIME type",
			fileName:    "file.xyz",
			contentType: "application/unknown",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileExtension(tt.fileName, tt.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFileExtension() error = %v, wantErr %v, fileName=%s, contentType=%s",
					err, tt.wantErr, tt.fileName, tt.contentType)
			}
		})
	}
}

func TestValidateAllowedMimeTypes(t *testing.T) {
	config := DefaultFileValidationConfig()

	tests := []struct {
		name        string
		contentType string
		fileType    model.FileType
		config      FileValidationConfig
		wantErr     bool
	}{
		{
			name:        "Allowed image MIME type",
			contentType: "image/jpeg",
			fileType:    model.FileTypeProfileImage,
			config:      config,
			wantErr:     false,
		},
		{
			name:        "Allowed PDF for CV",
			contentType: "application/pdf",
			fileType:    model.FileTypeCV,
			config:      config,
			wantErr:     false,
		},
		{
			name:        "Allowed DOCX for CV",
			contentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			fileType:    model.FileTypeCV,
			config:      config,
			wantErr:     false,
		},
		{
			name:        "Disallowed MIME type for image",
			contentType: "application/pdf",
			fileType:    model.FileTypeProfileImage,
			config:      config,
			wantErr:     true,
		},
		{
			name:        "Disallowed MIME type for CV",
			contentType: "image/jpeg",
			fileType:    model.FileTypeCV,
			config:      config,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAllowedMimeTypes(tt.contentType, tt.fileType, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAllowedMimeTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultFileValidationConfig(t *testing.T) {
	config := DefaultFileValidationConfig()

	if config.MaxProfileImageSize == 0 {
		t.Error("MaxProfileImageSize should not be 0")
	}

	if config.MaxCVSize == 0 {
		t.Error("MaxCVSize should not be 0")
	}

	if len(config.AllowedImageTypes) == 0 {
		t.Error("AllowedImageTypes should not be empty")
	}

	if len(config.AllowedCVTypes) == 0 {
		t.Error("AllowedCVTypes should not be empty")
	}

	// Verify expected MIME types are present
	imageTypes := make(map[string]bool)
	for _, v := range config.AllowedImageTypes {
		imageTypes[v] = true
	}

	expectedImages := []string{"image/jpeg", "image/png", "image/gif"}
	for _, expected := range expectedImages {
		if !imageTypes[expected] {
			t.Errorf("Expected image type %s not found in config", expected)
		}
	}

	cvTypes := make(map[string]bool)
	for _, v := range config.AllowedCVTypes {
		cvTypes[v] = true
	}

	expectedCVTypes := []string{"application/pdf", "application/msword"}
	for _, expected := range expectedCVTypes {
		if !cvTypes[expected] {
			t.Errorf("Expected CV type %s not found in config", expected)
		}
	}
}
