package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/service"
)

type ChatHandler struct {
	svc *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

type chatReq struct {
	UserID      int64  `json:"user_id" binding:"required"`
	CityID      int64  `json:"city_id" binding:"required"`
	CharacterID int64  `json:"character_id" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

// Chat handles POST /api/chat
func (h *ChatHandler) Chat(c *gin.Context) {
	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id, city_id, character_id, message are required"))
		return
	}

	reply, err := h.svc.Chat(c.Request.Context(), req.UserID, req.CityID, req.CharacterID, req.Message)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, errorResp("NOT_FOUND", "resource not found"))
			return
		}
		if strings.Contains(err.Error(), "ai timeout") || errors.Is(err, aiPkg.ErrAITimeout) {
			c.JSON(http.StatusGatewayTimeout, errorResp("AI_TIMEOUT", "AI service timeout"))
			return
		}
		c.JSON(http.StatusBadGateway, errorResp("AI_UPSTREAM_ERROR", "AI service error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
