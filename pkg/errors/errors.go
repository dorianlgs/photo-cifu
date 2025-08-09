package errors

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Cause   error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Common error constructors
func BadRequest(message string, cause error) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
		Cause:   cause,
	}
}

func InternalError(message string, cause error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
		Cause:   cause,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func ValidationError(message string, cause error) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Status:  http.StatusUnprocessableEntity,
		Cause:   cause,
	}
}

// HandleError converts AppError to PocketBase response
func HandleError(e *core.RequestEvent, err error) error {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Status {
		case http.StatusBadRequest:
			return e.BadRequestError(appErr.Message, appErr.Cause)
		case http.StatusNotFound:
			return e.NotFoundError(appErr.Message, appErr.Cause)
		case http.StatusUnprocessableEntity:
			return e.BadRequestError(appErr.Message, appErr.Cause)
		default:
			return e.InternalServerError(appErr.Message, appErr.Cause)
		}
	}
	return e.InternalServerError("Internal server error", err)
}