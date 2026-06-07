package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/middleware"
)

// Handlers holds all handler instances.
type Handlers struct {
	City        *CityHandler
	User        *UserHandler
	Visit       *VisitHandler
	Game        *GameHandler
	Chat        *ChatHandler
	Comment     *CommentHandler
	Guess       *GuessHandler
	Checkin     *CheckinHandler
	Achievement *AchievementHandler
	Asset       *AssetHandler
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
		if h.User != nil {
			api.GET("/users/:user_id/profile", h.User.Profile)
			api.PATCH("/users/:user_id/profile", h.User.UpdateProfile)
		}

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

		// Comments / danmaku
		api.GET("/comments", h.Comment.List)
		api.POST("/comments", h.Comment.Create)

		// Panorama guessing captions (rate limited)
		api.POST("/guess/caption", aiLimiter, h.Guess.Caption)
		api.POST("/guess/challenges", h.Guess.CreateChallenge)
		api.GET("/guess/challenges/:code", h.Guess.GetChallenge)
		api.POST("/guess/challenges/:code/answers", h.Guess.AnswerChallenge)

		// Checkin (rate limited for image generation)
		api.POST("/checkin/generate-image", aiLimiter, h.Checkin.GenerateImage)
		api.GET("/checkin/image-tasks/:task_id", h.Checkin.GetImageTask)
		api.POST("/checkin/image-tasks/:task_id/retry", aiLimiter, h.Checkin.RetryImageTask)
		api.POST("/checkin", h.Checkin.Create)

		// Achievements
		api.GET("/users/:user_id/achievements", h.Achievement.Wall)
		api.GET("/users/:user_id/assets", h.Asset.Assets)

		// Admin (protected by ADMIN_TOKEN)
		if h.Admin != nil {
			admin := api.Group("/admin", h.Admin.AuthMiddleware())
			{
				admin.GET("/catalog/coverage", h.Admin.Coverage)
				admin.PATCH("/cities/:city_id", h.Admin.UpdateCity)
				admin.PATCH("/landmarks/:landmark_id", h.Admin.UpdateLandmark)
				admin.PATCH("/characters/:character_id", h.Admin.UpdateCharacter)
				admin.PATCH("/foods/:food_id", h.Admin.UpdateFood)
				admin.POST("/cities/:city_id/landmarks", h.Admin.CreateLandmark)
				admin.POST("/cities/:city_id/characters", h.Admin.CreateCharacter)
				admin.POST("/cities/:city_id/foods", h.Admin.CreateFood)
				admin.DELETE("/landmarks/:landmark_id", h.Admin.DeleteLandmark)
				admin.DELETE("/characters/:character_id", h.Admin.DeleteCharacter)
				admin.DELETE("/foods/:food_id", h.Admin.DeleteFood)
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
