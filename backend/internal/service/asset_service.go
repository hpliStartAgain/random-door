package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type assetVisitRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]model.CityVisit, error)
}

type assetCheckinRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]model.Checkin, error)
}

type AssetService struct {
	userRepo       checkinUserFinder
	cityRepo       checkinCityRepo
	visitRepo      assetVisitRepo
	checkinRepo    assetCheckinRepo
	achievementSvc *AchievementService
}

func NewAssetService(
	userRepo checkinUserFinder,
	cityRepo checkinCityRepo,
	visitRepo assetVisitRepo,
	checkinRepo assetCheckinRepo,
	achievementSvc *AchievementService,
) *AssetService {
	return &AssetService{
		userRepo:       userRepo,
		cityRepo:       cityRepo,
		visitRepo:      visitRepo,
		checkinRepo:    checkinRepo,
		achievementSvc: achievementSvc,
	}
}

type AssetResult struct {
	VisitedCities       []VisitedCity  `json:"visited_cities"`
	Posters             []PosterAsset  `json:"posters"`
	AchievementProgress []ProgressItem `json:"achievement_progress"`
}

type VisitedCity struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Province  string `json:"province"`
	VisitedAt string `json:"visited_at"`
}

type PosterAsset struct {
	CheckinID         int64   `json:"checkin_id"`
	CityID            int64   `json:"city_id"`
	CityName          string  `json:"city_name"`
	LandmarkName      *string `json:"landmark_name,omitempty"`
	GeneratedImageURL string  `json:"generated_image_url"`
	CreatedAt         string  `json:"created_at"`
}

func (s *AssetService) Assets(ctx context.Context, userID int64) (*AssetResult, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	visits, err := s.visitRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list visits: %w", err)
	}
	visited := make([]VisitedCity, 0)
	seen := make(map[int64]bool)
	for _, visit := range visits {
		if seen[visit.CityID] {
			continue
		}
		seen[visit.CityID] = true
		city, err := s.cityRepo.FindByID(ctx, visit.CityID)
		if err != nil {
			continue
		}
		visited = append(visited, VisitedCity{
			ID:        city.ID,
			Name:      city.Name,
			Province:  city.Province,
			VisitedAt: formatAPITime(visit.CreatedAt),
		})
	}

	checkins, err := s.checkinRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list checkins: %w", err)
	}
	posters := make([]PosterAsset, 0)
	for _, checkin := range checkins {
		if checkin.GeneratedImageURL == nil || *checkin.GeneratedImageURL == "" {
			continue
		}
		city, err := s.cityRepo.FindByID(ctx, checkin.CityID)
		if err != nil {
			continue
		}
		var landmarkName *string
		if checkin.LandmarkID != nil {
			landmarks, err := s.cityRepo.ListLandmarks(ctx, checkin.CityID)
			if err == nil {
				for _, landmark := range landmarks {
					if landmark.ID == *checkin.LandmarkID {
						name := landmark.Name
						landmarkName = &name
						break
					}
				}
			}
		}
		posters = append(posters, PosterAsset{
			CheckinID:         checkin.ID,
			CityID:            checkin.CityID,
			CityName:          city.Name,
			LandmarkName:      landmarkName,
			GeneratedImageURL: *checkin.GeneratedImageURL,
			CreatedAt:         formatAPITime(checkin.CreatedAt),
		})
	}

	wall, err := s.achievementSvc.Wall(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &AssetResult{
		VisitedCities:       visited,
		Posters:             posters,
		AchievementProgress: wall.Progress,
	}, nil
}

func formatAPITime(t time.Time) string {
	return t.Format(time.RFC3339)
}
