package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type CityRepository interface {
	ListAll(ctx context.Context) ([]model.City, error)
	FindByID(ctx context.Context, id int64) (*model.City, error)
	ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
	ListFoods(ctx context.Context, cityID int64) ([]model.Food, error)
	ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error)
	ListAllCounts(ctx context.Context) (map[int64]model.CityCounts, error)
}

type CityService struct {
	cityRepo CityRepository
}

func NewCityService(cityRepo CityRepository) *CityService {
	return &CityService{cityRepo: cityRepo}
}

// CityListItem is a lightweight city representation for the list endpoint.
type CityListItem struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Province       string   `json:"province"`
	Lat            float64  `json:"lat"`
	Lng            float64  `json:"lng"`
	CoverImageURL  *string  `json:"cover_image_url,omitempty"`
	Tags           []string `json:"tags"`
	LandmarkCount  int      `json:"landmark_count"`
	FoodCount      int      `json:"food_count"`
	CharacterCount int      `json:"character_count"`
}

// CityDetail is the full city response with related entities.
type CityDetail struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	Province           string          `json:"province"`
	Lat                float64         `json:"lat"`
	Lng                float64         `json:"lng"`
	Intro              *string         `json:"intro,omitempty"`
	CoverImageURL      *string         `json:"cover_image_url,omitempty"`
	DialectSample      *string         `json:"dialect_sample,omitempty"`
	DialectExplanation *string         `json:"dialect_explanation,omitempty"`
	Tags               []string        `json:"tags"`
	Landmarks          []LandmarkItem  `json:"landmarks"`
	Foods              []FoodItem      `json:"foods"`
	Characters         []CharacterItem `json:"characters"`
}

type LandmarkItem struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	ImageURL      *string `json:"image_url,omitempty"`
	Description   *string `json:"description,omitempty"`
	SoundscapeURL *string `json:"soundscape_url,omitempty"`
}

type FoodItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	ImageURL    *string `json:"image_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CharacterItem struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	CharacterType string  `json:"character_type"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	DialectStyle  *string `json:"dialect_style,omitempty"`
	RoleTitle     *string `json:"role_title,omitempty"`
	LifeSpan      *string `json:"life_span,omitempty"`
	IntroQuote    *string `json:"intro_quote,omitempty"`
}

func (s *CityService) List(ctx context.Context) ([]CityListItem, error) {
	cities, err := s.cityRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	counts, err := s.cityRepo.ListAllCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list city counts: %w", err)
	}

	result := make([]CityListItem, 0, len(cities))
	for _, c := range cities {
		tags, err := s.cityRepo.ListTags(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("list tags for city %d: %w", c.ID, err)
		}
		tagNames := make([]string, 0, len(tags))
		for _, t := range tags {
			tagNames = append(tagNames, t.Tag)
		}
		cnt := counts[c.ID]
		result = append(result, CityListItem{
			ID:             c.ID,
			Name:           c.Name,
			Province:       c.Province,
			Lat:            c.Lat,
			Lng:            c.Lng,
			CoverImageURL:  c.CoverImageURL,
			Tags:           tagNames,
			LandmarkCount:  cnt.LandmarkCount,
			FoodCount:      cnt.FoodCount,
			CharacterCount: cnt.CharacterCount,
		})
	}
	return result, nil
}

func (s *CityService) Detail(ctx context.Context, cityID int64) (*CityDetail, error) {
	if cityID <= 0 {
		return nil, invalidParam("city_id must be a positive integer")
	}

	city, err := s.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	tags, err := s.cityRepo.ListTags(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("list tags for city %d: %w", cityID, err)
	}
	tagNames := make([]string, 0, len(tags))
	for _, t := range tags {
		tagNames = append(tagNames, t.Tag)
	}

	landmarks, err := s.cityRepo.ListLandmarks(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("list landmarks for city %d: %w", cityID, err)
	}
	lms := make([]LandmarkItem, 0, len(landmarks))
	for _, l := range landmarks {
		lms = append(lms, LandmarkItem{
			ID: l.ID, Name: l.Name, ImageURL: l.ImageURL, Description: l.Description, SoundscapeURL: l.SoundscapeURL,
		})
	}

	foods, err := s.cityRepo.ListFoods(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("list foods for city %d: %w", cityID, err)
	}
	fs := make([]FoodItem, 0, len(foods))
	for _, f := range foods {
		fs = append(fs, FoodItem{
			ID: f.ID, Name: f.Name, ImageURL: f.ImageURL, Description: f.Description,
		})
	}

	chars, err := s.cityRepo.ListCharacters(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("list characters for city %d: %w", cityID, err)
	}
	cs := make([]CharacterItem, 0, len(chars))
	for _, ch := range chars {
		// Persona and Prompt are NOT included (json:"-" on model)
		cs = append(cs, CharacterItem{
			ID: ch.ID, Name: ch.Name, CharacterType: ch.CharacterType,
			AvatarURL: ch.AvatarURL, DialectStyle: ch.DialectStyle,
			RoleTitle: ch.RoleTitle, LifeSpan: ch.LifeSpan, IntroQuote: ch.IntroQuote,
		})
	}

	return &CityDetail{
		ID: city.ID, Name: city.Name, Province: city.Province,
		Lat: city.Lat, Lng: city.Lng,
		Intro: city.Intro, CoverImageURL: city.CoverImageURL,
		DialectSample: city.DialectSample, DialectExplanation: city.DialectExplanation,
		Tags: tagNames, Landmarks: lms, Foods: fs, Characters: cs,
	}, nil
}
