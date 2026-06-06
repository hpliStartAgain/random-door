package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

// Standard error code constants — prevents typos across handlers.
const (
	ErrCodeInvalidParam  = "INVALID_PARAM"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodePermission    = "PERMISSION_DENIED"
	ErrCodeRateLimit     = "RATE_LIMITED"
	ErrCodeAITimeout     = "AI_TIMEOUT"
	ErrCodeAIUpstream    = "AI_UPSTREAM_ERROR"
	ErrCodeFileTooLarge  = "FILE_TOO_LARGE"
	ErrCodeUnsupportedMedia = "UNSUPPORTED_MEDIA"
	ErrCodeInternalError = "INTERNAL_ERROR"
)

// errorResp creates a standard error response per api-contract.md 0.3.
func errorResp(code, message string) gin.H {
	return gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidParam):
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, service.ClientMessage(err)))
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResp(ErrCodeNotFound, service.ClientMessage(err)))
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, errorResp(ErrCodeConflict, service.ClientMessage(err)))
	case errors.Is(err, service.ErrPermission):
		c.JSON(http.StatusForbidden, errorResp(ErrCodePermission, service.ClientMessage(err)))
	default:
		c.JSON(http.StatusInternalServerError, errorResp(ErrCodeInternalError, "internal server error"))
	}
}
