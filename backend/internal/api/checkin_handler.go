package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusRequestEntityTooLarge, errorResp(ErrCodeFileTooLarge, "file exceeds 5MB"))
			return
		}
		c.JSON(http.StatusUnsupportedMediaType, errorResp(ErrCodeUnsupportedMedia, "file type not supported"))
		return
	}

	selfiePath, err := h.storage.Save(file, "selfies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "failed to save file"))
		return
	}

	scenePath := ""
	if sceneFile, sceneErr := c.FormFile("scene_file"); sceneErr == nil {
		if verr := h.validator.Validate(sceneFile); verr != nil {
			if errors.Is(verr, upload.ErrFileTooLarge) {
				c.JSON(http.StatusRequestEntityTooLarge, errorResp(ErrCodeFileTooLarge, "scene_file exceeds 5MB"))
				return
			}
			c.JSON(http.StatusUnsupportedMediaType, errorResp(ErrCodeUnsupportedMedia, "scene_file type not supported"))
			return
		}
		scenePath, err = h.storage.Save(sceneFile, "scenes")
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "failed to save scene file"))
			return
		}
	} else if !errors.Is(sceneErr, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid scene_file"))
		return
	}

	result, err := h.svc.EnqueueImage(c.Request.Context(), userID, cityID, landmarkID, selfiePath, scenePath)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, result)
}

// GetImageTask handles GET /api/checkin/image-tasks/:task_id.
func (h *CheckinHandler) GetImageTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid task_id"))
		return
	}
	userID, err := userIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	result, err := h.svc.GetImageTask(c.Request.Context(), userID, taskID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

type retryImageTaskReq struct {
	UserID int64 `json:"user_id"`
}

// RetryImageTask handles POST /api/checkin/image-tasks/:task_id/retry.
func (h *CheckinHandler) RetryImageTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid task_id"))
		return
	}
	userID, parseErr := userIDFromRequest(c)
	if parseErr != nil {
		var req retryImageTaskReq
		if err := c.ShouldBindJSON(&req); err == nil && req.UserID > 0 {
			userID = req.UserID
		} else {
			c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", parseErr.Error()))
			return
		}
	}
	result, err := h.svc.RetryImageTask(c.Request.Context(), userID, taskID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, result)
}

func userIDFromRequest(c *gin.Context) (int64, error) {
	raw := c.Query("user_id")
	if raw == "" {
		raw = c.GetHeader("X-User-Id")
	}
	if raw == "" {
		return 0, errors.New("user_id is required")
	}
	userID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || userID <= 0 {
		return 0, errors.New("user_id must be a positive integer")
	}
	return userID, nil
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
