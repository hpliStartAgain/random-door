package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type GuessHandler struct {
	svc          *service.GuessService
	challengeSvc *service.GuessChallengeService
}

func NewGuessHandler(svc *service.GuessService) *GuessHandler {
	return &GuessHandler{svc: svc}
}

func (h *GuessHandler) WithChallengeService(svc *service.GuessChallengeService) *GuessHandler {
	h.challengeSvc = svc
	return h
}

type guessCaptionReq struct {
	UserID     *int64 `json:"user_id"`
	CityID     int64  `json:"city_id" binding:"required"`
	TargetName string `json:"target_name"`
	SceneHint  string `json:"scene_hint"`
}

// Caption handles POST /api/guess/caption.
func (h *GuessHandler) Caption(c *gin.Context) {
	var req guessCaptionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "city_id is required"))
		return
	}
	result, err := h.svc.GenerateCaption(c.Request.Context(), service.GuessCaptionRequest{
		UserID:     req.UserID,
		CityID:     req.CityID,
		TargetName: strings.TrimSpace(req.TargetName),
		SceneHint:  strings.TrimSpace(req.SceneHint),
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
