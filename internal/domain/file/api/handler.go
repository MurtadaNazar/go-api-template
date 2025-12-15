package api

import (
	"go_platform_template/internal/domain/file/dto"
	"go_platform_template/internal/domain/file/model"
	"go_platform_template/internal/domain/file/service"
	apperrors "go_platform_template/internal/shared/errors"
	"go_platform_template/internal/shared/response"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileHandler struct {
	service *service.FileService
	logger  *zap.SugaredLogger
}

func NewFileHandler(s *service.FileService, logger *zap.SugaredLogger) *FileHandler {
	return &FileHandler{service: s, logger: logger}
}

// Upload godoc
// @Summary Upload a file
// @Description Upload a file (profile image or CV) for the authenticated user
// @Tags files
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param type query string true "File type" Enums(profile_image, cv)
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Success 200 {object} dto.UploadResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 413 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /files/upload [post]

func (h *FileHandler) Upload(c *gin.Context) {
	requestID, _ := c.Get("RequestID")
	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		h.logger.Warnw("upload attempt without authentication", "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "User authentication required"))
		return
	}

	// Parse userID string to uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warnw("invalid user ID format", "user_id", userIDStr, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.BadRequestError, "Invalid user ID"))
		return
	}

	fType := c.Query("type")
	if fType != string(model.FileTypeProfileImage) && fType != string(model.FileTypeCV) {
		h.logger.Warnw("invalid file type", "file_type", fType, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid file type",
			"Must be 'profile_image' or 'cv'",
		))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Warnw("file not provided in upload", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.BadRequestError, "File not provided"))
		return
	}

	// Validate file using service validation
	contentType := file.Header.Get("Content-Type")
	if err := h.service.ValidateUpload(file.Filename, file.Size, contentType, model.FileType(fType)); err != nil {
		h.logger.Warnw("file validation failed", "filename", file.Filename, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"File validation failed",
			err.Error(),
		))
		return
	}

	src, err := file.Open()
	if err != nil {
		h.logger.Errorw("failed to open uploaded file", "filename", file.Filename, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to open file"))
		return
	}
	defer src.Close()

	// Generate secure object name with timestamp to prevent collisions
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().Format("20060102-150405")
	objectName := userID.String() + "/" + baseName + "_" + timestamp + ext

	// Upload file
	uploaded, err := h.service.Upload(
		userID, // pass uuid.UUID instead of string
		model.FileType(fType),
		src,
		objectName,
		file.Size,
		contentType,
		file.Filename,
	)
	if err != nil {
		h.logger.Errorw("failed to upload file", "user_id", userID, "filename", file.Filename, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to upload file"))
		return
	}

	// Generate signed URL
	url, err := h.service.GetSignedURL(uploaded.Path, 15*time.Minute)
	if err != nil {
		h.logger.Errorw("failed to generate signed URL", "file_path", uploaded.Path, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to generate access URL"))
		return
	}

	h.logger.Infow("file uploaded successfully", "user_id", userID, "file_id", uploaded.ID, "request_id", requestID)
	requestIDStr, ok := requestID.(string)
	if !ok {
		requestIDStr = "unknown"
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(dto.UploadResponse{
		FileID:       uploaded.ID.String(),
		URL:          url,
		Path:         uploaded.Path,
		Type:         string(uploaded.Type),
		Size:         uploaded.Size,
		OriginalName: uploaded.OriginalName,
		MimeType:     uploaded.MimeType,
		UploadedAt:   uploaded.UploadedAt,
		ExpiresIn:    "15 minutes",
	}, requestIDStr))
}

// GetFile godoc
// @Summary Get file by filename
// @Description Get a temporary signed URL to access a file
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param filename path string true "File path/name"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /files/{filename} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	requestID, _ := c.Get("RequestID")
	objectName := c.Param("filename")
	if objectName == "" {
		h.logger.Warnw("get file without filename", "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.BadRequestError, "Filename is required"))
		return
	}

	// Verify file exists before generating URL
	exists, err := h.service.FileExists(objectName)
	if err != nil {
		h.logger.Errorw("failed to check file existence", "filename", objectName, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to check file existence"))
		return
	}
	if !exists {
		h.logger.Warnw("file not found", "filename", objectName, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.NotFoundError, "File not found"))
		return
	}

	url, err := h.service.GetSignedURL(objectName, 15*time.Minute)
	if err != nil {
		h.logger.Errorw("failed to generate signed URL", "filename", objectName, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to generate access URL"))
		return
	}

	h.logger.Infow("file signed URL generated", "filename", objectName, "request_id", requestID)
	requestIDStr, ok := requestID.(string)
	if !ok {
		requestIDStr = "unknown"
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(dto.GetFileResponse{
		URL:       url,
		ExpiresIn: "15 minutes",
	}, requestIDStr))
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file and its metadata
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param filename path string true "File path/name"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /files/{filename} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	requestID, _ := c.Get("RequestID")
	userIDStr := c.GetString("userID")
	objectName := c.Param("filename")

	if userIDStr == "" {
		h.logger.Warnw("delete file attempt without authentication", "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "User authentication required"))
		return
	}

	// Convert string userID to uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warnw("invalid user ID format on file delete", "user_id", userIDStr, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.BadRequestError, "Invalid user ID"))
		return
	}

	// Verify the file belongs to the user
	file, err := h.service.GetFileByPath(c.Request.Context(), objectName)
	if err != nil {
		h.logger.Warnw("file not found for deletion", "filename", objectName, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.NotFoundError, "File not found"))
		return
	}

	if file.UserID != userID {
		h.logger.Warnw("unauthorized file delete attempt", "user_id", userID, "file_owner", file.UserID, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.ForbiddenError, "You do not have permission to delete this file"))
		return
	}

	if err := h.service.Delete(objectName); err != nil {
		h.logger.Errorw("failed to delete file", "filename", objectName, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to delete file"))
		return
	}

	h.logger.Infow("file deleted successfully", "filename", objectName, "user_id", userID, "request_id", requestID)
	requestIDStr, ok := requestID.(string)
	if !ok {
		requestIDStr = "unknown"
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "file deleted successfully"}, requestIDStr))
}

// GetUserFiles godoc
// @Summary Get user's files
// @Description Get all files uploaded by the authenticated user
// @Tags files
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /files [get]
func (h *FileHandler) GetUserFiles(c *gin.Context) {
	requestID, _ := c.Get("RequestID")
	userID := c.GetString("userID")
	if userID == "" {
		h.logger.Warnw("get user files without authentication", "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "User authentication required"))
		return
	}

	files, err := h.service.GetFilesByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorw("failed to retrieve user files", "user_id", userID, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to retrieve files"))
		return
	}

	responseData := dto.UserFilesResponse{
		Count: len(files),
		Files: make([]dto.FileInfo, len(files)),
	}

	for i, file := range files {
		url, err := h.service.GetSignedURL(file.Path, 15*time.Minute)
		if err != nil {
			h.logger.Warnw("failed to generate signed URL for file", "file_id", file.ID, "error", err, "request_id", requestID)
			url = ""
		}
		responseData.Files[i] = dto.FileInfo{
			ID:           file.ID.String(),
			Path:         file.Path,
			Type:         string(file.Type),
			Size:         file.Size,
			OriginalName: file.OriginalName,
			MimeType:     file.MimeType,
			UploadedAt:   file.UploadedAt,
			URL:          url,
		}
	}

	h.logger.Infow("user files retrieved", "user_id", userID, "count", len(files), "request_id", requestID)
	requestIDStr, ok := requestID.(string)
	if !ok {
		requestIDStr = "unknown"
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(responseData, requestIDStr))
}
