package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

const maxSceneHintLen = 200

type guessCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
	ListFoods(ctx context.Context, cityID int64) ([]model.Food, error)
	ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error)
}

type guessLLM interface {
	Chat(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

type GuessService struct {
	cityRepo guessCityRepo
	llm      guessLLM
}

func NewGuessService(cityRepo guessCityRepo, llm guessLLM) *GuessService {
	return &GuessService{cityRepo: cityRepo, llm: llm}
}

type GuessCaptionRequest struct {
	UserID     *int64 `json:"user_id,omitempty"`
	CityID     int64  `json:"city_id"`
	TargetName string `json:"target_name"`
	SceneHint  string `json:"scene_hint"`
}

type GuessCaptionResult struct {
	Weibo    string   `json:"weibo"`
	Moments  string   `json:"moments"`
	Hashtags []string `json:"hashtags"`
}

func (s *GuessService) GenerateCaption(ctx context.Context, req GuessCaptionRequest) (*GuessCaptionResult, error) {
	if req.CityID <= 0 {
		return nil, invalidParam("city_id must be a positive integer")
	}
	if req.UserID != nil && *req.UserID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	req.TargetName = strings.TrimSpace(req.TargetName)
	req.SceneHint = strings.TrimSpace(req.SceneHint)
	if len([]rune(req.SceneHint)) > maxSceneHintLen {
		return nil, invalidParam("scene_hint too long (max 200 characters)")
	}

	city, err := s.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}
	if req.TargetName == "" {
		req.TargetName = city.Name
	}
	if req.SceneHint == "" {
		req.SceneHint = s.defaultSceneHint(ctx, city.ID)
	}

	result := fallbackGuessCaption(city.Name, req.TargetName)
	if s.llm == nil {
		return result, nil
	}

	if weibo, err := s.generateOne(ctx, city.Name, req.TargetName, req.SceneHint, "微博"); err == nil && strings.TrimSpace(weibo) != "" {
		result.Weibo = clampRunes(strings.TrimSpace(weibo), 120)
	} else if err != nil {
		slog.Warn("guess weibo caption fallback", "error", err, "city_id", city.ID)
	}
	if moments, err := s.generateOne(ctx, city.Name, req.TargetName, req.SceneHint, "微信朋友圈"); err == nil && strings.TrimSpace(moments) != "" {
		result.Moments = clampRunes(strings.TrimSpace(moments), 90)
	} else if err != nil {
		slog.Warn("guess moments caption fallback", "error", err, "city_id", city.ID)
	}
	return result, nil
}

func (s *GuessService) generateOne(ctx context.Context, cityName, targetName, sceneHint, platform string) (string, error) {
	systemPrompt, userPrompt := ai.BuildGuessCaptionPrompt(cityName, targetName, sceneHint, platform)
	return s.llm.Chat(ctx, systemPrompt, userPrompt)
}

func (s *GuessService) defaultSceneHint(ctx context.Context, cityID int64) string {
	landmarks, _ := s.cityRepo.ListLandmarks(ctx, cityID)
	foods, _ := s.cityRepo.ListFoods(ctx, cityID)
	characters, _ := s.cityRepo.ListCharacters(ctx, cityID)
	parts := make([]string, 0, 3)
	if len(landmarks) > 0 {
		parts = append(parts, "地标："+landmarks[0].Name)
	}
	if len(foods) > 0 {
		parts = append(parts, "美食："+foods[0].Name)
	}
	if len(characters) > 0 {
		parts = append(parts, "人物："+characters[0].Name)
	}
	if len(parts) == 0 {
		return "全景视角里的城市文化场景。"
	}
	return strings.Join(parts, "；")
}

func fallbackGuessCaption(cityName, targetName string) *GuessCaptionResult {
	if targetName == "" {
		targetName = cityName
	}
	return &GuessCaptionResult{
		Weibo:    fmt.Sprintf("猜猜我在哪？镜头停在%s的一角，%s的风物已经露出线索。#任意门漫游# #%s#", targetName, cityName, cityName),
		Moments:  fmt.Sprintf("把视角停在%s，像是闯进%s的一段现场。猜猜这是哪里？", targetName, cityName),
		Hashtags: []string{"任意门漫游", cityName},
	}
}

func clampRunes(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
