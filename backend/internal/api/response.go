package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
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
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", service.ClientMessage(err)))
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResp("NOT_FOUND", service.ClientMessage(err)))
	default:
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
	}
}
