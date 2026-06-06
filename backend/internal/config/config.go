package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	LLM    LLMConfig
	Image  ImageConfig
	Upload UploadConfig
	CORS   CORSConfig
	AI     AIConfig
	Admin  AdminConfig
	Log    LogConfig
}

type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
}

type ServerConfig struct {
	Port      int
	StaticDir string
	UploadDir string
	SeedDir   string
}

type LLMConfig struct {
	APIBase string
	APIKey  string
	Model   string
}

type ImageConfig struct {
	APIBase string
	APIKey  string
	Model   string
}

type UploadConfig struct {
	MaxSizeMB    int
	AllowedTypes []string
}

type CORSConfig struct {
	AllowOrigins []string
}

type AIConfig struct {
	TimeoutSeconds        int
	Timeout               time.Duration
	ChatDailyLimit        int
	ImageDailyLimit       int
	WorkerIntervalSeconds int
	WorkerConcurrency     int
	MaxTaskAttempts       int
}

type AdminConfig struct {
	Token string
}

type LogConfig struct {
	Level string
}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("DB_HOST", "127.0.0.1")
	v.SetDefault("DB_PORT", 3306)
	v.SetDefault("DB_USER", "root")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_NAME", "city_roam")
	v.SetDefault("DB_MAX_OPEN_CONNS", 20)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("STATIC_DIR", "./static")
	v.SetDefault("UPLOAD_DIR", "./uploads")
	v.SetDefault("SEED_DIR", "./data/seed")
	v.SetDefault("UPLOAD_MAX_SIZE_MB", 5)
	v.SetDefault("UPLOAD_ALLOWED_TYPES", "jpg,jpeg,png,webp")
	v.SetDefault("AI_TIMEOUT_SECONDS", 30)
	v.SetDefault("AI_CHAT_DAILY_LIMIT", 30)
	v.SetDefault("AI_IMAGE_DAILY_LIMIT", 5)
	v.SetDefault("AI_WORKER_INTERVAL_SECONDS", 2)
	v.SetDefault("AI_WORKER_CONCURRENCY", 1)
	v.SetDefault("AI_MAX_TASK_ATTEMPTS", 3)
	v.SetDefault("LLM_MODEL", "deepseek-v4-flash")
	v.SetDefault("IMAGE_MODEL", "wan2.7-image-pro")
	v.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost,http://localhost:5173")
	v.SetDefault("LOG_LEVEL", "info")

	timeoutSec := v.GetInt("AI_TIMEOUT_SECONDS")

	cfg := &Config{
		DB: DBConfig{
			Host:         v.GetString("DB_HOST"),
			Port:         v.GetInt("DB_PORT"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			Name:         v.GetString("DB_NAME"),
			MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
		},
		Server: ServerConfig{
			Port:      v.GetInt("SERVER_PORT"),
			StaticDir: v.GetString("STATIC_DIR"),
			UploadDir: v.GetString("UPLOAD_DIR"),
			SeedDir:   v.GetString("SEED_DIR"),
		},
		LLM: LLMConfig{
			APIBase: v.GetString("LLM_API_BASE"),
			APIKey:  v.GetString("LLM_API_KEY"),
			Model:   v.GetString("LLM_MODEL"),
		},
		Image: ImageConfig{
			APIBase: v.GetString("IMAGE_API_BASE"),
			APIKey:  v.GetString("IMAGE_API_KEY"),
			Model:   v.GetString("IMAGE_MODEL"),
		},
		Upload: UploadConfig{
			MaxSizeMB:    v.GetInt("UPLOAD_MAX_SIZE_MB"),
			AllowedTypes: strings.Split(v.GetString("UPLOAD_ALLOWED_TYPES"), ","),
		},
		CORS: CORSConfig{
			AllowOrigins: strings.Split(v.GetString("CORS_ALLOW_ORIGINS"), ","),
		},
		AI: AIConfig{
			TimeoutSeconds:        timeoutSec,
			Timeout:               time.Duration(timeoutSec) * time.Second,
			ChatDailyLimit:        v.GetInt("AI_CHAT_DAILY_LIMIT"),
			ImageDailyLimit:       v.GetInt("AI_IMAGE_DAILY_LIMIT"),
			WorkerIntervalSeconds: v.GetInt("AI_WORKER_INTERVAL_SECONDS"),
			WorkerConcurrency:     v.GetInt("AI_WORKER_CONCURRENCY"),
			MaxTaskAttempts:       v.GetInt("AI_MAX_TASK_ATTEMPTS"),
		},
		Admin: AdminConfig{
			Token: v.GetString("ADMIN_TOKEN"),
		},
		Log: LogConfig{
			Level: v.GetString("LOG_LEVEL"),
		},
	}

	log.Printf("[config] loaded: DB=%s:%d/%s, Server=:%d", cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.Server.Port)
	return cfg
}
