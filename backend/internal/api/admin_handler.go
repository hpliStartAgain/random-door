package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
	"github.com/your-org/city-roam/backend/internal/upload"
)

// AdminHandler handles admin-only media upload operations.
type AdminHandler struct {
	svc       *service.AdminService
	storage   *upload.Storage
	validator *upload.Validator
	token     string
}

func NewAdminHandler(svc *service.AdminService, storage *upload.Storage, validator *upload.Validator, token string) *AdminHandler {
	return &AdminHandler{svc: svc, storage: storage, validator: validator, token: token}
}

// AuthMiddleware checks X-Admin-Token or Authorization: Bearer <token>.
func (h *AdminHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.token == "" {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable,
				errorResp("ADMIN_DISABLED", "admin API is not configured (ADMIN_TOKEN not set)"))
			return
		}
		tok := c.GetHeader("X-Admin-Token")
		if tok == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tok = auth[7:]
			}
		}
		if tok != h.token {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				errorResp("UNAUTHORIZED", "invalid or missing admin token"))
			return
		}
		c.Next()
	}
}

// Coverage handles GET /api/admin/catalog/coverage.
func (h *AdminHandler) Coverage(c *gin.Context) {
	result, err := h.svc.Coverage(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) ListTags(c *gin.Context) {
	result, err := h.svc.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) RenameTag(c *gin.Context) {
	var req service.RenameTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "tag is required"))
		return
	}
	if err := h.svc.RenameTag(c.Request.Context(), c.Param("tag"), req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AdminHandler) DeleteTag(c *gin.Context) {
	if err := h.svc.DeleteTag(c.Request.Context(), c.Param("tag")); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminHandler) ListAchievements(c *gin.Context) {
	result, err := h.svc.ListAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) CreateAchievement(c *gin.Context) {
	var req service.CreateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement payload"))
		return
	}
	result, err := h.svc.CreateAchievement(c.Request.Context(), req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *AdminHandler) UpdateAchievement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("achievement_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement_id"))
		return
	}
	var req service.UpdateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement payload"))
		return
	}
	if err := h.svc.UpdateAchievement(c.Request.Context(), id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AdminHandler) DeleteAchievement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("achievement_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement_id"))
		return
	}
	if err := h.svc.DeleteAchievement(c.Request.Context(), id); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UpdateCity handles PATCH /api/admin/cities/:city_id.
func (h *AdminHandler) UpdateCity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	var req service.UpdateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid request body"))
		return
	}
	if err := h.svc.UpdateCity(c.Request.Context(), id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// UpdateLandmark handles PATCH /api/admin/landmarks/:landmark_id.
func (h *AdminHandler) UpdateLandmark(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("landmark_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid landmark_id"))
		return
	}
	var req service.UpdatePOIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid request body"))
		return
	}
	if err := h.svc.UpdateLandmark(c.Request.Context(), id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// UpdateFood handles PATCH /api/admin/foods/:food_id.
func (h *AdminHandler) UpdateFood(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid food_id"))
		return
	}
	var req service.UpdatePOIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid request body"))
		return
	}
	if err := h.svc.UpdateFood(c.Request.Context(), id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// UpdateCharacter handles PATCH /api/admin/characters/:character_id.
func (h *AdminHandler) UpdateCharacter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("character_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid character_id"))
		return
	}
	var req service.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid request body"))
		return
	}
	if err := h.svc.UpdateCharacter(c.Request.Context(), id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// CreateLandmark handles POST /api/admin/cities/:city_id/landmarks.
func (h *AdminHandler) CreateLandmark(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	var req service.CreatePOIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "name is required"))
		return
	}
	result, err := h.svc.CreateLandmark(c.Request.Context(), cityID, req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// CreateFood handles POST /api/admin/cities/:city_id/foods.
func (h *AdminHandler) CreateFood(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	var req service.CreatePOIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "name is required"))
		return
	}
	result, err := h.svc.CreateFood(c.Request.Context(), cityID, req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// CreateCharacter handles POST /api/admin/cities/:city_id/characters.
func (h *AdminHandler) CreateCharacter(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	var req service.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "name is required"))
		return
	}
	result, err := h.svc.CreateCharacter(c.Request.Context(), cityID, req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// DeleteLandmark handles DELETE /api/admin/landmarks/:landmark_id.
func (h *AdminHandler) DeleteLandmark(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("landmark_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid landmark_id"))
		return
	}
	if err := h.svc.DeleteLandmark(c.Request.Context(), id); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// DeleteFood handles DELETE /api/admin/foods/:food_id.
func (h *AdminHandler) DeleteFood(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid food_id"))
		return
	}
	if err := h.svc.DeleteFood(c.Request.Context(), id); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// DeleteCharacter handles DELETE /api/admin/characters/:character_id.
func (h *AdminHandler) DeleteCharacter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("character_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid character_id"))
		return
	}
	if err := h.svc.DeleteCharacter(c.Request.Context(), id); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UploadCityCover handles POST /api/admin/cities/:city_id/cover-image
func (h *AdminHandler) UploadCityCover(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	urlPath, err := h.saveImage(c, "cities")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	if err := h.svc.UpdateCity(c.Request.Context(), cityID, service.UpdateCityRequest{CoverImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"cover_image_url": urlPath})
}

// UploadLandmarkImage handles POST /api/admin/landmarks/:landmark_id/image
func (h *AdminHandler) UploadLandmarkImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("landmark_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid landmark_id"))
		return
	}
	urlPath, err := h.saveImage(c, "landmarks")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	if err := h.svc.UpdateLandmark(c.Request.Context(), id, service.UpdatePOIRequest{ImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": urlPath})
}

// UploadCharacterAvatar handles POST /api/admin/characters/:character_id/avatar
func (h *AdminHandler) UploadCharacterAvatar(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("character_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid character_id"))
		return
	}
	urlPath, err := h.saveImage(c, "avatars")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	if err := h.svc.UpdateCharacter(c.Request.Context(), id, service.UpdateCharacterRequest{AvatarURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar_url": urlPath})
}

// UploadFoodImage handles POST /api/admin/foods/:food_id/image
func (h *AdminHandler) UploadFoodImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid food_id"))
		return
	}
	urlPath, err := h.saveImage(c, "foods")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	if err := h.svc.UpdateFood(c.Request.Context(), id, service.UpdatePOIRequest{ImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": urlPath})
}

// UploadAchievementBadge handles POST /api/admin/achievements/:achievement_id/badge
func (h *AdminHandler) UploadAchievementBadge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("achievement_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement_id"))
		return
	}
	urlPath, err := h.saveImage(c, "badges")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", err.Error()))
		return
	}
	if err := h.svc.UpdateAchievement(c.Request.Context(), id, service.UpdateAchievementRequest{BadgeURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"badge_url": urlPath})
}

func (h *AdminHandler) saveImage(c *gin.Context, subDir string) (string, error) {
	fh, err := c.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("file is required")
	}
	if err := h.validator.Validate(fh); err != nil {
		return "", fmt.Errorf("file validation failed: %s", err.Error())
	}
	return h.storage.Save(fh, subDir)
}

// ---- URL-bind endpoints (PATCH) — accept JSON {"url":"https://..."} ----

type urlBindReq struct {
	URL string `json:"url" binding:"required"`
}

// BindCityCoverURL handles PATCH /api/admin/cities/:city_id/cover-image
func (h *AdminHandler) BindCityCoverURL(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid city_id"))
		return
	}
	var req urlBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "url is required"))
		return
	}
	urlPath, err := h.svc.ImportImageURL(c.Request.Context(), req.URL)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if err := h.svc.UpdateCity(c.Request.Context(), cityID, service.UpdateCityRequest{CoverImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"cover_image_url": urlPath})
}

// BindLandmarkImageURL handles PATCH /api/admin/landmarks/:landmark_id/image
func (h *AdminHandler) BindLandmarkImageURL(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("landmark_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid landmark_id"))
		return
	}
	var req urlBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "url is required"))
		return
	}
	urlPath, err := h.svc.ImportImageURL(c.Request.Context(), req.URL)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if err := h.svc.UpdateLandmark(c.Request.Context(), id, service.UpdatePOIRequest{ImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": urlPath})
}

// BindCharacterAvatarURL handles PATCH /api/admin/characters/:character_id/avatar
func (h *AdminHandler) BindCharacterAvatarURL(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("character_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid character_id"))
		return
	}
	var req urlBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "url is required"))
		return
	}
	urlPath, err := h.svc.ImportImageURL(c.Request.Context(), req.URL)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if err := h.svc.UpdateCharacter(c.Request.Context(), id, service.UpdateCharacterRequest{AvatarURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar_url": urlPath})
}

// BindFoodImageURL handles PATCH /api/admin/foods/:food_id/image
func (h *AdminHandler) BindFoodImageURL(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid food_id"))
		return
	}
	var req urlBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "url is required"))
		return
	}
	urlPath, err := h.svc.ImportImageURL(c.Request.Context(), req.URL)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if err := h.svc.UpdateFood(c.Request.Context(), id, service.UpdatePOIRequest{ImageURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": urlPath})
}

// BindAchievementBadgeURL handles PATCH /api/admin/achievements/:achievement_id/badge
func (h *AdminHandler) BindAchievementBadgeURL(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("achievement_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "invalid achievement_id"))
		return
	}
	var req urlBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "url is required"))
		return
	}
	urlPath, err := h.svc.ImportImageURL(c.Request.Context(), req.URL)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if err := h.svc.UpdateAchievement(c.Request.Context(), id, service.UpdateAchievementRequest{BadgeURL: &urlPath}); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"badge_url": urlPath})
}
