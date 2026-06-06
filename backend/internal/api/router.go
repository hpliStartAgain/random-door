package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/middleware"
)

// Handlers holds all handler instances.
type Handlers struct {
	City        *CityHandler
	Visit       *VisitHandler
	Game        *GameHandler
	Chat        *ChatHandler
	Checkin     *CheckinHandler
	Achievement *AchievementHandler
	Admin       *AdminHandler
}

// NewRouter creates and configures the Gin engine with all routes.
func NewRouter(h Handlers, corsOrigins []string, staticDir, uploadDir string, aiRateLimit float64) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware
	r.Use(middleware.Recover())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(corsOrigins))

	// Static file serving
	r.Static("/static", staticDir)
	r.Static("/uploads", uploadDir)

	// Rate limiter for AI endpoints
	aiLimiter := middleware.RateLimit(aiRateLimit, 5)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// User
		api.POST("/users/anonymous", h.Visit.CreateAnonymousUser)

		// Cities
		api.GET("/cities", h.City.List)
		api.GET("/cities/:city_id", h.City.Detail)

		// Visits
		api.POST("/visits/free", h.Visit.CreateFreeVisit)

		// Game
		api.POST("/game/init", h.Game.Init)
		api.POST("/game/roll", h.Game.Roll)

		// Chat (rate limited)
		api.POST("/chat", aiLimiter, h.Chat.Chat)

		// Checkin (rate limited for image generation)
		api.POST("/checkin/generate-image", aiLimiter, h.Checkin.GenerateImage)
		api.POST("/checkin", h.Checkin.Create)

		// Achievements
		api.GET("/users/:user_id/achievements", h.Achievement.Wall)

		// Admin (protected by ADMIN_TOKEN)
		if h.Admin != nil {
			admin := api.Group("/admin", h.Admin.AuthMiddleware())
			{
				// File upload
				admin.POST("/cities/:city_id/cover-image", h.Admin.UploadCityCover)
				admin.POST("/landmarks/:landmark_id/image", h.Admin.UploadLandmarkImage)
				admin.POST("/characters/:character_id/avatar", h.Admin.UploadCharacterAvatar)
				admin.POST("/foods/:food_id/image", h.Admin.UploadFoodImage)
				// URL bind (paste external link)
				admin.PATCH("/cities/:city_id/cover-image", h.Admin.BindCityCoverURL)
				admin.PATCH("/landmarks/:landmark_id/image", h.Admin.BindLandmarkImageURL)
				admin.PATCH("/characters/:character_id/avatar", h.Admin.BindCharacterAvatarURL)
				admin.PATCH("/foods/:food_id/image", h.Admin.BindFoodImageURL)
			}
		}
	}

	return r
}
