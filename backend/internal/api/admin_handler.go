package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/upload"
	"gorm.io/gorm"
)

// AdminHandler handles admin-only media upload operations.
type AdminHandler struct {
	db        *gorm.DB
	storage   *upload.Storage
	validator *upload.Validator
	token     string
}

func NewAdminHandler(db *gorm.DB, storage *upload.Storage, validator *upload.Validator, token string) *AdminHandler {
	return &AdminHandler{db: db, storage: storage, validator: validator, token: token}
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
	if err := h.db.Model(&model.City{}).Where("id = ?", cityID).
		Update("cover_image_url", urlPath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
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
	if err := h.db.Model(&model.Landmark{}).Where("id = ?", id).
		Update("image_url", urlPath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
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
	if err := h.db.Model(&model.Character{}).Where("id = ?", id).
		Update("avatar_url", urlPath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
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
	if err := h.db.Model(&model.Food{}).Where("id = ?", id).
		Update("image_url", urlPath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": urlPath})
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
	if err := h.db.Model(&model.City{}).Where("id = ?", cityID).
		Update("cover_image_url", req.URL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"cover_image_url": req.URL})
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
	if err := h.db.Model(&model.Landmark{}).Where("id = ?", id).
		Update("image_url", req.URL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": req.URL})
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
	if err := h.db.Model(&model.Character{}).Where("id = ?", id).
		Update("avatar_url", req.URL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar_url": req.URL})
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
	if err := h.db.Model(&model.Food{}).Where("id = ?", id).
		Update("image_url", req.URL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": req.URL})
}
