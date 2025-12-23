package api

import (
	"fmt"
	"go_platform_template/internal/domain/user/dto"
	"go_platform_template/internal/domain/user/service"
	apperrors "go_platform_template/internal/shared/errors"
	"go_platform_template/internal/shared/response"
	"go_platform_template/internal/platform/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service   service.UserService
	validator *validation.Validator
	logger    *zap.SugaredLogger
}

func NewUserHandler(s service.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		service:   s,
		validator: validation.New(),
		logger:    logger,
	}
}

// ListUsers godoc
// @Summary List users with pagination, filters, and sorting
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by email"
// @Param user_type query string false "Filter by user type"
// @Param sort_by query string false "Sort by field (created_at or username)"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	offset := 0
	limit := 20

	if v := c.Query("offset"); v != "" {
		if _, err := fmt.Sscan(v, &offset); err != nil {
			h.logger.Warnw("invalid offset value", "offset", v, "request_id", requestID)
			_ = c.Error(apperrors.NewAppErrorWithDetails(
				apperrors.BadRequestError,
				"Invalid offset value",
				err.Error(),
			))
			return
		}
	}
	if v := c.Query("limit"); v != "" {
		if _, err := fmt.Sscan(v, &limit); err != nil {
			h.logger.Warnw("invalid limit value", "limit", v, "request_id", requestID)
			_ = c.Error(apperrors.NewAppErrorWithDetails(
				apperrors.BadRequestError,
				"Invalid limit value",
				err.Error(),
			))
			return
		}
	}

	filters := make(map[string]interface{})
	if v := c.Query("username"); v != "" {
		filters["username"] = v
	}
	if v := c.Query("email"); v != "" {
		filters["email"] = v
	}
	if v := c.Query("user_type"); v != "" {
		filters["user_type"] = v
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	users, err := h.service.List(c.Request.Context(), offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		h.logger.Errorw("failed to list users", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to fetch users"))
		return
	}

	for _, u := range users {
		u.Password = ""
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(users, requestID))
}

// GetUser godoc
// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	id := c.Param("id")
	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Failed to fetch user"))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, response.NewSuccessResponse(user, requestID))
}

// Register godoc
// @Summary Register a new user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "User to create"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/ [post]
func (h *UserHandler) Register(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("invalid register request", "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid request payload",
			err.Error(),
		))
		return
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		h.logger.Warnw("validation error on register", "error", err, "request_id", requestID)
		_ = c.Error(err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Registration failed"))
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, response.NewSuccessResponse(user, requestID))
}

// UpdateUser godoc
// @Summary Update an existing user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UserUpdateRequest true "Updated user data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	id := c.Param("id")
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("invalid update request", "user_id", id, "error", err, "request_id", requestID)
		_ = c.Error(apperrors.NewAppErrorWithDetails(
			apperrors.BadRequestError,
			"Invalid request payload",
			err.Error(),
		))
		return
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		h.logger.Warnw("validation error on update", "user_id", id, "error", err, "request_id", requestID)
		_ = c.Error(err)
		return
	}

	updated, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Update failed"))
		return
	}

	updated.Password = ""
	c.JSON(http.StatusOK, response.NewSuccessResponse(updated, requestID))
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	requestIDVal, _ := c.Get("RequestID")
	requestID, ok := requestIDVal.(string)
	if !ok {
		requestID = "unknown"
	}

	id := c.Param("id")
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := apperrors.IsAppError(err); ok {
			_ = c.Error(appErr)
			return
		}
		_ = c.Error(apperrors.NewAppError(apperrors.InternalError, "Deletion failed"))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "user deleted successfully"}, requestID))
}
