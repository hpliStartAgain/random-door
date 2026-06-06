package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/api"
	"github.com/your-org/city-roam/backend/internal/config"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/repository"
	"github.com/your-org/city-roam/backend/internal/seed"
	"github.com/your-org/city-roam/backend/internal/service"
	"github.com/your-org/city-roam/backend/internal/upload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env if exists
	_ = godotenv.Load("../.env")

	// 1. Load config
	cfg := config.Load()

	// 2. Connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("database connected", "host", cfg.DB.Host, "port", cfg.DB.Port, "db", cfg.DB.Name)

	// 3. AutoMigrate (create tables if not exist)
	if err := db.AutoMigrate(
		&model.User{},
		&model.City{},
		&model.CityTag{},
		&model.Landmark{},
		&model.Food{},
		&model.Character{},
		&model.CityVisit{},
		&model.DiceRoll{},
		&model.Checkin{},
		&model.Achievement{},
		&model.UserAchievement{},
		&model.ChatMessage{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	slog.Info("database migrated")

	// 4. Validate and upsert seed data atomically
	if err := seed.Load(context.Background(), db, cfg.Server.SeedDir); err != nil {
		log.Fatalf("seed data failed: %v", err)
	}
	slog.Info("seed data synchronized")

	// 5. Build dependencies (manual DI)
	userRepo := repository.NewUserRepo(db)
	cityRepo := repository.NewCityRepo(db)
	visitRepo := repository.NewVisitRepo(db)
	diceRepo := repository.NewDiceRepo(db)
	chatRepo := repository.NewChatRepo(db)
	checkinRepo := repository.NewCheckinRepo(db)
	achRepo := repository.NewAchievementRepo(db)

	llmClient := ai.NewLLMClient(cfg.LLM.APIBase, cfg.LLM.APIKey, cfg.LLM.Model, cfg.AI.Timeout)
	imageClient := ai.NewImageClient(cfg.Image.APIBase, cfg.Image.APIKey, cfg.AI.Timeout)
	validator := upload.NewValidator(cfg.Upload.MaxSizeMB, cfg.Upload.AllowedTypes)
	storage := upload.NewStorage(cfg.Server.UploadDir)

	citySvc := service.NewCityService(cityRepo)
	visitSvc := service.NewVisitService(userRepo, cityRepo, visitRepo)
	gameStore := repository.NewGameStore(db, diceRepo, visitRepo)
	gameSvc := service.NewGameService(userRepo, cityRepo, visitRepo, gameStore)
	chatSvc := service.NewChatService(cityRepo, chatRepo, llmClient)
	achSvc := service.NewAchievementService(db, achRepo)
	checkinStore := repository.NewCheckinStore(db, checkinRepo)
	checkinSvc := service.NewCheckinService(userRepo, cityRepo, checkinStore, imageClient, storage)

	handlers := api.Handlers{
		City:        api.NewCityHandler(citySvc),
		Visit:       api.NewVisitHandler(visitSvc),
		Game:        api.NewGameHandler(gameSvc),
		Chat:        api.NewChatHandler(chatSvc),
		Checkin:     api.NewCheckinHandler(checkinSvc, validator, storage),
		Achievement: api.NewAchievementHandler(achSvc),
		Admin:       api.NewAdminHandler(db, storage, validator, cfg.Admin.Token),
	}

	// 6. Create and start Gin router
	router := api.NewRouter(handlers, cfg.CORS.AllowOrigins, cfg.Server.StaticDir, cfg.Server.UploadDir, 2.0)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("server starting", "addr", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
