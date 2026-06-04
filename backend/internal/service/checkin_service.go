package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type checkinUserFinder interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type checkinCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
}

// checkinStore writes a checkin and evaluates achievements in a single transaction.
type checkinStore interface {
	CreateAndEvaluate(ctx context.Context, checkin *model.Checkin) ([]model.Achievement, error)
}

type imageGenerator interface {
	Generate(ctx context.Context, selfiePath, refImagePath, prompt string) ([]byte, error)
}

type imageStorage interface {
	SaveBytes(data []byte, subDir, ext string) (string, error)
}

type CheckinService struct {
	userRepo    checkinUserFinder
	cityRepo    checkinCityRepo
	store       checkinStore
	imageClient imageGenerator
	storage     imageStorage
}

func NewCheckinService(
	userRepo checkinUserFinder,
	cityRepo checkinCityRepo,
	store checkinStore,
	imageClient imageGenerator,
	storage imageStorage,
) *CheckinService {
	return &CheckinService{
		userRepo: userRepo, cityRepo: cityRepo, store: store,
		imageClient: imageClient, storage: storage,
	}
}

// GenerateImageResult is the response from image generation.
type GenerateImageResult struct {
	Status            string `json:"status"`
	GeneratedImageURL string `json:"generated_image_url"`
}

// GenerateImage processes selfie + landmark reference to create a cyber check-in photo.
func (s *CheckinService) GenerateImage(ctx context.Context, userID, cityID, landmarkID int64, selfiePath string) (*GenerateImageResult, error) {
	// Get city and landmark info for prompt
	city, err := s.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	landmarks, _ := s.cityRepo.ListLandmarks(ctx, cityID)
	var landmarkName string
	var refImagePath string
	for _, l := range landmarks {
		if l.ID == landmarkID {
			landmarkName = l.Name
			if l.ImageURL != nil {
				refImagePath = *l.ImageURL
			}
			break
		}
	}
	if landmarkName == "" {
		return nil, invalidParam("landmark not found")
	}

	// Build image generation prompt
	prompt := ai.BuildImagePrompt(city.Name, landmarkName)

	// Call image generation API
	imgBytes, err := s.imageClient.Generate(ctx, selfiePath, refImagePath, prompt)
	if err != nil {
		slog.Error("image generation failed", "error", err, "user_id", userID)
		return nil, err
	}

	// Save generated image
	generatedURL, err := s.storage.SaveBytes(imgBytes, "generated", ".png")
	if err != nil {
		return nil, fmt.Errorf("save generated image: %w", err)
	}

	slog.Info("image generated", "user_id", userID, "url", generatedURL)
	return &GenerateImageResult{
		Status:            "success",
		GeneratedImageURL: generatedURL,
	}, nil
}

// CreateCheckinRequest is the input for creating a check-in.
type CreateCheckinRequest struct {
	UserID            int64   `json:"user_id"`
	CityID            int64   `json:"city_id"`
	LandmarkID        *int64  `json:"landmark_id,omitempty"`
	VisitID           *int64  `json:"visit_id,omitempty"`
	GeneratedImageURL *string `json:"generated_image_url,omitempty"`
}

// CreateCheckinResult is the response from creating a check-in.
type CreateCheckinResult struct {
	CheckinID            int64              `json:"checkin_id"`
	UnlockedAchievements []AchievementBrief `json:"unlocked_achievements"`
}

type AchievementBrief struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// Create records a check-in and evaluates achievements within a single transaction.
func (s *CheckinService) Create(ctx context.Context, req CreateCheckinRequest) (*CreateCheckinResult, error) {
	// checkins has no FK constraints, so verify references explicitly.
	if _, err := s.userRepo.FindByID(ctx, req.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	if _, err := s.cityRepo.FindByID(ctx, req.CityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	checkin := &model.Checkin{
		UserID:            req.UserID,
		CityID:            req.CityID,
		LandmarkID:        req.LandmarkID,
		VisitID:           req.VisitID,
		GeneratedImageURL: req.GeneratedImageURL,
	}

	newAchs, err := s.store.CreateAndEvaluate(ctx, checkin)
	if err != nil {
		return nil, fmt.Errorf("create checkin: %w", err)
	}

	achBriefs := make([]AchievementBrief, 0, len(newAchs))
	for _, a := range newAchs {
		achBriefs = append(achBriefs, AchievementBrief{
			Code: a.Code, Name: a.Name, Description: a.Description,
		})
	}

	slog.Info("checkin created", "user_id", req.UserID, "city_id", req.CityID,
		"checkin_id", checkin.ID, "new_achievements", len(newAchs))

	return &CreateCheckinResult{
		CheckinID:            checkin.ID,
		UnlockedAchievements: achBriefs,
	}, nil
}
