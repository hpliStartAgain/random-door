package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type createGuessChallengeReq struct {
	UserID       *int64 `json:"user_id"`
	CityID       int64  `json:"city_id" binding:"required"`
	TargetName   string `json:"target_name"`
	ImageURL     string `json:"image_url"`
	ImageDataURL string `json:"image_data_url"`
	Caption      string `json:"caption"`
}

type answerGuessChallengeReq struct {
	AnswerText string `json:"answer_text" binding:"required"`
}

func (h *GuessHandler) CreateChallenge(c *gin.Context) {
	if h.challengeSvc == nil {
		c.JSON(http.StatusNotFound, errorResp(ErrCodeNotFound, "challenge service not available"))
		return
	}
	var req createGuessChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "city_id is required"))
		return
	}
	result, err := h.challengeSvc.Create(c.Request.Context(), service.CreateGuessChallengeRequest{
		UserID:       req.UserID,
		CityID:       req.CityID,
		TargetName:   strings.TrimSpace(req.TargetName),
		ImageURL:     strings.TrimSpace(req.ImageURL),
		ImageDataURL: strings.TrimSpace(req.ImageDataURL),
		Caption:      strings.TrimSpace(req.Caption),
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *GuessHandler) GetChallenge(c *gin.Context) {
	if h.challengeSvc == nil {
		c.JSON(http.StatusNotFound, errorResp(ErrCodeNotFound, "challenge service not available"))
		return
	}
	result, err := h.challengeSvc.Get(c.Request.Context(), c.Param("code"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *GuessHandler) AnswerChallenge(c *gin.Context) {
	if h.challengeSvc == nil {
		c.JSON(http.StatusNotFound, errorResp(ErrCodeNotFound, "challenge service not available"))
		return
	}
	var req answerGuessChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "answer_text is required"))
		return
	}
	result, err := h.challengeSvc.Answer(c.Request.Context(), c.Param("code"), req.AnswerText)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
