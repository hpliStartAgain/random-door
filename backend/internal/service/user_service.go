package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type ProfileUserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
	UpdateProfile(ctx context.Context, userID int64, fields map[string]any) error
}

type UserService struct {
	userRepo ProfileUserRepository
}

func NewUserService(userRepo ProfileUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type UserProfile struct {
	UserID        int64   `json:"user_id"`
	AnonymousID   string  `json:"anonymous_id"`
	Nickname      *string `json:"nickname,omitempty"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	Age           *int    `json:"age,omitempty"`
	HomeRegion    *string `json:"home_region,omitempty"`
	CurrentCityID *int64  `json:"current_city_id"`
}

type UpdateUserProfileRequest struct {
	Nickname   *string `json:"nickname,omitempty"`
	Age        *int    `json:"age,omitempty"`
	HomeRegion *string `json:"home_region,omitempty"`
}

func (s *UserService) Profile(ctx context.Context, userID int64) (*UserProfile, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return publicProfile(user), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req UpdateUserProfileRequest) (*UserProfile, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	fields := make(map[string]any)
	if req.Nickname != nil {
		value, err := normalizeProfileText("nickname", *req.Nickname, 64)
		if err != nil {
			return nil, err
		}
		fields["nickname"] = value
		user.Nickname = &value
	}
	if req.Age != nil {
		if *req.Age < 1 || *req.Age > 120 {
			return nil, invalidParam("age must be between 1 and 120")
		}
		fields["age"] = *req.Age
		user.Age = req.Age
	}
	if req.HomeRegion != nil {
		value, err := normalizeProfileText("home_region", *req.HomeRegion, 64)
		if err != nil {
			return nil, err
		}
		fields["home_region"] = value
		user.HomeRegion = &value
	}

	if err := s.userRepo.UpdateProfile(ctx, userID, fields); err != nil {
		return nil, fmt.Errorf("update user profile: %w", err)
	}
	return publicProfile(user), nil
}

func normalizeProfileText(field, value string, max int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", invalidParam(field + " must not be empty")
	}
	if utf8.RuneCountInString(value) > max {
		return "", invalidParam(fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return value, nil
}

func publicProfile(user *model.User) *UserProfile {
	return &UserProfile{
		UserID:        user.ID,
		AnonymousID:   user.AnonymousID,
		Nickname:      user.Nickname,
		AvatarURL:     user.AvatarURL,
		Age:           user.Age,
		HomeRegion:    user.HomeRegion,
		CurrentCityID: user.CurrentCityID,
	}
}
