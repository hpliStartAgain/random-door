package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/your-org/city-roam/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProfileUserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	UpdateProfile(ctx context.Context, userID int64, fields map[string]any) error
	UpdateAccount(ctx context.Context, userID int64, fields map[string]any) error
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
	Username      *string `json:"username,omitempty"`
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

type RegisterRequest struct {
	UserID   *int64  `json:"user_id,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Nickname *string `json:"nickname,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResult struct {
	UserID        int64   `json:"user_id"`
	AnonymousID   string  `json:"anonymous_id"`
	Username      string  `json:"username"`
	Nickname      *string `json:"nickname,omitempty"`
	CurrentCityID *int64  `json:"current_city_id"`
}

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,32}$`)

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

func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*AuthResult, error) {
	username, err := normalizeUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}
	if existing, err := s.userRepo.FindByUsername(ctx, username); err == nil && existing != nil {
		return nil, conflict("username already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find username: %w", err)
	}

	var nickname *string
	if req.Nickname != nil {
		value, err := normalizeProfileText("nickname", *req.Nickname, 64)
		if err != nil {
			return nil, err
		}
		nickname = &value
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	passwordHash := string(hashBytes)

	if req.UserID != nil && *req.UserID > 0 {
		user, err := s.userRepo.FindByID(ctx, *req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, notFound("user not found")
			}
			return nil, fmt.Errorf("find user: %w", err)
		}
		if user.Username != nil && *user.Username != "" {
			return nil, conflict("user already has an account")
		}
		fields := map[string]any{
			"username":      username,
			"password_hash": passwordHash,
		}
		user.Username = &username
		user.PasswordHash = &passwordHash
		if nickname != nil {
			fields["nickname"] = *nickname
			user.Nickname = nickname
		}
		if err := s.userRepo.UpdateAccount(ctx, user.ID, fields); err != nil {
			return nil, classifyAccountCreateError(err)
		}
		return authResult(user), nil
	}

	user := &model.User{
		AnonymousID:  uuid.NewString(),
		Username:     &username,
		PasswordHash: &passwordHash,
		Nickname:     nickname,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, classifyAccountCreateError(err)
	}
	return authResult(user), nil
}

func (s *UserService) Login(ctx context.Context, req LoginRequest) (*AuthResult, error) {
	username, err := normalizeUsername(req.Username)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, permissionDenied("username or password is incorrect")
		}
		return nil, fmt.Errorf("find username: %w", err)
	}
	if user.PasswordHash == nil || bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)) != nil {
		return nil, permissionDenied("username or password is incorrect")
	}
	return authResult(user), nil
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

func normalizeUsername(value string) (string, error) {
	value = strings.TrimSpace(value)
	if !usernamePattern.MatchString(value) {
		return "", invalidParam("username must be 3-32 letters, numbers, or underscores")
	}
	return strings.ToLower(value), nil
}

func validatePassword(value string) error {
	if len(value) < 6 || len(value) > 72 {
		return invalidParam("password must be 6-72 bytes")
	}
	return nil
}

func classifyAccountCreateError(err error) error {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return conflict("username already exists")
	}
	return fmt.Errorf("save account: %w", err)
}

func publicProfile(user *model.User) *UserProfile {
	return &UserProfile{
		UserID:        user.ID,
		AnonymousID:   user.AnonymousID,
		Username:      user.Username,
		Nickname:      user.Nickname,
		AvatarURL:     user.AvatarURL,
		Age:           user.Age,
		HomeRegion:    user.HomeRegion,
		CurrentCityID: user.CurrentCityID,
	}
}

func authResult(user *model.User) *AuthResult {
	username := ""
	if user.Username != nil {
		username = *user.Username
	}
	return &AuthResult{
		UserID:        user.ID,
		AnonymousID:   user.AnonymousID,
		Username:      username,
		Nickname:      user.Nickname,
		CurrentCityID: user.CurrentCityID,
	}
}
