package api

import (
	"go_platform_template/internal/domain/auth/model"
	"go_platform_template/internal/domain/auth/service"
	apperrors "go_platform_template/internal/shared/errors"
	"go_platform_template/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service *service.AuthService
	logger  *zap.SugaredLogger
}

func NewAuthHandler(s *service.AuthService, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{service: s, logger: logger}
}

// MeResponse represents the response for /me endpoint
type MeResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// Login godoc
// @Summary User login
// @Description Authenticates a user with email and password, returns access and refresh tokens
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("invalid login request", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid request payload",
			err.Error(),
		))
		return
	}

	access, refresh, err := h.service.Login(c.Request.Context(), req.EmailOrUsername, req.Password)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Login failed"))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(model.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, requestID))
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Rotates refresh token and returns new access & refresh tokens
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param refresh body model.RefreshRequest true "Refresh token"
// @Success 200 {object} model.RefreshResponse
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("invalid refresh request", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid request payload",
			err.Error(),
		))
		return
	}

	access, newRefresh, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Token refresh failed"))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(model.RefreshResponse{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, requestID))
}

// Logout godoc
// @Summary Logout user
// @Description Revokes a refresh token (blacklist) to log out the user
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param logout body model.RefreshRequest true "Refresh token to revoke"
// @Success 200 {object} response.SuccessResponse "Logged out successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Router /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("invalid logout request", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid request payload",
			err.Error(),
		))
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Logout failed"))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "logged out successfully"}, requestID))
}

// Me godoc
// @Summary Get current logged-in user info
// @Description Returns user ID and role from access token
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{
		"user_id": userID,
		"role":    role,
	}, requestID))
}
