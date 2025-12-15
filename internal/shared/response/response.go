package response

import "time"

// SuccessResponse is the standard success response wrapper
type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorResponse is the standard error response wrapper
type ErrorResponse struct {
	Error     string    `json:"error"`
	Type      string    `json:"type"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// PaginatedResponse is used for paginated endpoints
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Offset    int         `json:"offset"`
	Limit     int         `json:"limit"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, requestID string) *SuccessResponse {
	return &SuccessResponse{
		Data:      data,
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message, errType, details, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Type:      errType,
		Details:   details,
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, offset, limit int, requestID string) *PaginatedResponse {
	return &PaginatedResponse{
		Data:      data,
		Total:     total,
		Offset:    offset,
		Limit:     limit,
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	}
}
