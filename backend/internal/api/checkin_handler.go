package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/service"
	"github.com/your-org/city-roam/backend/internal/upload"
)

type CheckinHandler struct {
	svc       *service.CheckinService
	validator *upload.Validator
	storage   *upload.Storage
}

func NewCheckinHandler(svc *service.CheckinService, validator *upload.Validator, storage *upload.Storage) *CheckinHandler {
	return &CheckinHandler{svc: svc, validator: validator, storage: storage}
}

// GenerateImage handles POST /api/checkin/generate-image (multipart/form-data)
func (h *CheckinHandler) GenerateImage(c *gin.Context) {
	userID, err := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id is required"))
		return
	}
	cityID, err := strconv.ParseInt(c.PostForm("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "city_id is required"))
		return
	}
	landmarkID, err := strconv.ParseInt(c.PostForm("landmark_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "landmark_id is required"))
		return
	}

	file, err := c.FormFile("selfie_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "selfie_file is required"))
		return
	}

	// Validate file
	if verr := h.validator.Validate(file); verr != nil {
		if errors.Is(verr, upload.ErrFileTooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, errorResp("FILE_TOO_LARGE", "file exceeds 5MB"))
			return
		}
		c.JSON(http.StatusUnsupportedMediaType, errorResp("UNSUPPORTED_MEDIA", "file type not supported"))
		return
	}

	// Save selfie
	selfiePath, err := h.storage.Save(file, "selfies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "failed to save file"))
		return
	}

	result, err := h.svc.GenerateImage(c.Request.Context(), userID, cityID, landmarkID, selfiePath)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrInvalidParam) {
			writeServiceError(c, err)
			return
		}
		if errors.Is(err, aiPkg.ErrAITimeout) {
			c.JSON(http.StatusGatewayTimeout, errorResp("AI_TIMEOUT", "image generation timeout"))
			return
		}
		if errors.Is(err, aiPkg.ErrAIUpstream) {
			c.JSON(http.StatusBadGateway, errorResp("AI_UPSTREAM_ERROR", "image generation failed"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}

	c.JSON(http.StatusOK, result)
}

type checkinReq struct {
	UserID            int64   `json:"user_id" binding:"required"`
	CityID            int64   `json:"city_id" binding:"required"`
	LandmarkID        *int64  `json:"landmark_id"`
	VisitID           *int64  `json:"visit_id"`
	GeneratedImageURL *string `json:"generated_image_url"`
}

// Create handles POST /api/checkin
func (h *CheckinHandler) Create(c *gin.Context) {
	var req checkinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id and city_id are required"))
		return
	}

	result, err := h.svc.Create(c.Request.Context(), service.CreateCheckinRequest{
		UserID:            req.UserID,
		CityID:            req.CityID,
		LandmarkID:        req.LandmarkID,
		VisitID:           req.VisitID,
		GeneratedImageURL: req.GeneratedImageURL,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
