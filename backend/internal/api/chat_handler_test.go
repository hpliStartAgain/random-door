package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestChatHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// Mock handler for testing 400
	r.POST("/api/chat", func(c *gin.Context) {
		var req chatReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "missing params"))
			return
		}
	})

	body := []byte(`{"user_id": 1}`) // Missing other fields
	req, _ := http.NewRequest(http.MethodPost, "/api/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
