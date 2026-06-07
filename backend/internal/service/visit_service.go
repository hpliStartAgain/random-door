package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByAnonymousID(ctx context.Context, anonymousID string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type CityFinder interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
}

type VisitRepository interface {
	Create(ctx context.Context, visit *model.CityVisit) error
}

type visitAchievementEvaluator interface {
	UnlockForUser(ctx context.Context, userID int64) ([]model.Achievement, error)
}

type VisitService struct {
	userRepo     UserRepository
	cityRepo     CityFinder
	visitRepo    VisitRepository
	achievements visitAchievementEvaluator
}

func NewVisitService(userRepo UserRepository, cityRepo CityFinder, visitRepo VisitRepository) *VisitService {
	return &VisitService{userRepo: userRepo, cityRepo: cityRepo, visitRepo: visitRepo}
}

func (s *VisitService) WithAchievementEvaluator(evaluator visitAchievementEvaluator) *VisitService {
	s.achievements = evaluator
	return s
}

// CreateAnonymousUser finds or creates a user by anonymous_id.
func (s *VisitService) CreateAnonymousUser(ctx context.Context, anonymousID string) (*model.User, error) {
	normalizedID, err := normalizeAnonymousID(anonymousID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByAnonymousID(ctx, normalizedID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find user: %w", err)
	}

	newUser := &model.User{AnonymousID: normalizedID}
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		// Another request may have inserted the same browser UUID concurrently.
		user, findErr := s.userRepo.FindByAnonymousID(ctx, normalizedID)
		if findErr == nil {
			return user, nil
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return newUser, nil
}

// CreateFreeVisit records a free-mode city visit.
type FreeVisitResult struct {
	Visit                *model.CityVisit
	UnlockedAchievements []AchievementBrief
}

func (s *VisitService) CreateFreeVisit(ctx context.Context, userID, cityID int64, source string) (*FreeVisitResult, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	if cityID <= 0 {
		return nil, invalidParam("city_id must be a positive integer")
	}
	if source == "" {
		source = "map_click"
	}
	if source != "map_click" && source != "search" {
		return nil, invalidParam("source must be map_click or search")
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	if _, err := s.cityRepo.FindByID(ctx, cityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	visit := &model.CityVisit{
		UserID:    userID,
		CityID:    cityID,
		VisitMode: "free",
		Source:    &source,
	}
	if err := s.visitRepo.Create(ctx, visit); err != nil {
		return nil, fmt.Errorf("create visit: %w", err)
	}
	var unlocked []AchievementBrief
	if s.achievements != nil {
		achievements, err := s.achievements.UnlockForUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("evaluate achievements: %w", err)
		}
		unlocked = briefAchievements(achievements)
	}
	return &FreeVisitResult{Visit: visit, UnlockedAchievements: unlocked}, nil
}

func normalizeAnonymousID(anonymousID string) (string, error) {
	if len(anonymousID) != 36 || strings.TrimSpace(anonymousID) != anonymousID {
		return "", invalidParam("anonymous_id must be a UUID")
	}
	parsed, err := uuid.Parse(anonymousID)
	if err != nil || parsed.String() != strings.ToLower(anonymousID) {
		return "", invalidParam("anonymous_id must be a UUID")
	}
	return parsed.String(), nil
}
