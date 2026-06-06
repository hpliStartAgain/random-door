package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// List handles GET /api/comments.
func (h *CommentHandler) List(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID, err := strconv.ParseInt(c.Query("target_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "target_id must be a positive integer"))
		return
	}
	limit := 0
	if raw := c.Query("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "limit must be between 1 and 100"))
			return
		}
	}

	result, err := h.svc.List(c.Request.Context(), targetType, targetID, limit)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

type createCommentReq struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   int64  `json:"target_id" binding:"required"`
	UserID     *int64 `json:"user_id"`
	Nickname   string `json:"nickname"`
	Content    string `json:"content" binding:"required"`
}

// Create handles POST /api/comments.
func (h *CommentHandler) Create(c *gin.Context) {
	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "target_type, target_id, and content are required"))
		return
	}
	result, err := h.svc.Create(c.Request.Context(), service.CreateCommentRequest{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		UserID:     req.UserID,
		Nickname:   req.Nickname,
		Content:    req.Content,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
