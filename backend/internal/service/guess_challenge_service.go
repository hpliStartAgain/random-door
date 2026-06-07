package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

const (
	guessChallengeTTL     = 7 * 24 * time.Hour
	maxGuessImageBytes    = 5 * 1024 * 1024
	guessChallengeCodeLen = 8
)

type guessChallengeRepo interface {
	CreateChallenge(ctx context.Context, challenge *model.GuessChallenge) error
	FindChallengeByCode(ctx context.Context, code string) (*model.GuessChallenge, error)
	CreateAnswer(ctx context.Context, answer *model.GuessAnswer) error
}

type guessChallengeCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
}

type guessImageStorage interface {
	SaveBytes(data []byte, subDir, ext string) (string, error)
}

type GuessChallengeService struct {
	repo     guessChallengeRepo
	cityRepo guessChallengeCityRepo
	storage  guessImageStorage
	now      func() time.Time
}

func NewGuessChallengeService(repo guessChallengeRepo, cityRepo guessChallengeCityRepo, storage guessImageStorage) *GuessChallengeService {
	return &GuessChallengeService{
		repo:     repo,
		cityRepo: cityRepo,
		storage:  storage,
		now:      time.Now,
	}
}

type CreateGuessChallengeRequest struct {
	UserID       *int64 `json:"user_id,omitempty"`
	CityID       int64  `json:"city_id"`
	TargetName   string `json:"target_name"`
	ImageURL     string `json:"image_url"`
	ImageDataURL string `json:"image_data_url"`
	Caption      string `json:"caption"`
}

type GuessChallengeResponse struct {
	Code       string    `json:"code"`
	ShareURL   string    `json:"share_url,omitempty"`
	CityID     int64     `json:"city_id"`
	CityName   string    `json:"city_name,omitempty"`
	TargetName *string   `json:"target_name,omitempty"`
	ImageURL   *string   `json:"image_url,omitempty"`
	Caption    *string   `json:"caption,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type GuessAnswerResponse struct {
	AnswerID   int64  `json:"answer_id"`
	IsCorrect  bool   `json:"is_correct"`
	AnswerText string `json:"answer_text"`
	CityName   string `json:"city_name"`
	TargetName string `json:"target_name,omitempty"`
	Message    string `json:"message"`
}

func (s *GuessChallengeService) Create(ctx context.Context, req CreateGuessChallengeRequest) (*GuessChallengeResponse, error) {
	if req.CityID <= 0 {
		return nil, invalidParam("city_id must be a positive integer")
	}
	if req.UserID != nil && *req.UserID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	city, err := s.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	targetName, err := optionalText("target_name", req.TargetName, 128)
	if err != nil {
		return nil, err
	}
	caption, err := optionalText("caption", req.Caption, 300)
	if err != nil {
		return nil, err
	}
	imageURL, err := s.resolveChallengeImage(req.ImageDataURL, req.ImageURL, city.CoverImageURL)
	if err != nil {
		return nil, err
	}

	now := s.now()
	challenge := &model.GuessChallenge{
		Code:       newGuessCode(),
		UserID:     req.UserID,
		CityID:     req.CityID,
		TargetName: targetName,
		ImageURL:   imageURL,
		Caption:    caption,
		ExpiresAt:  now.Add(guessChallengeTTL),
	}
	for i := 0; i < 3; i++ {
		if err := s.repo.CreateChallenge(ctx, challenge); err == nil {
			return s.challengeResponse(challenge, city), nil
		} else if i == 2 {
			return nil, fmt.Errorf("create guess challenge: %w", err)
		}
		challenge.Code = newGuessCode()
	}
	return nil, fmt.Errorf("create guess challenge: exhausted retries")
}

func (s *GuessChallengeService) Get(ctx context.Context, code string) (*GuessChallengeResponse, error) {
	challenge, city, err := s.loadActiveChallenge(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.challengeResponse(challenge, city), nil
}

func (s *GuessChallengeService) Answer(ctx context.Context, code, answerText string) (*GuessAnswerResponse, error) {
	challenge, city, err := s.loadActiveChallenge(ctx, code)
	if err != nil {
		return nil, err
	}
	answerText = strings.TrimSpace(answerText)
	if answerText == "" {
		return nil, invalidParam("answer_text must not be empty")
	}
	if utf8.RuneCountInString(answerText) > 64 {
		return nil, invalidParam("answer_text must be at most 64 characters")
	}
	target := deref(challenge.TargetName)
	isCorrect := guessAnswerMatches(answerText, city.Name, target)
	answer := &model.GuessAnswer{
		ChallengeCode: challenge.Code,
		AnswerText:    answerText,
		IsCorrect:     isCorrect,
	}
	if err := s.repo.CreateAnswer(ctx, answer); err != nil {
		return nil, fmt.Errorf("create guess answer: %w", err)
	}
	message := fmt.Sprintf("还差一点，这里是%s。", city.Name)
	if isCorrect {
		message = fmt.Sprintf("猜对了，这里是%s。", city.Name)
	}
	return &GuessAnswerResponse{
		AnswerID:   answer.ID,
		IsCorrect:  isCorrect,
		AnswerText: answer.AnswerText,
		CityName:   city.Name,
		TargetName: target,
		Message:    message,
	}, nil
}

func (s *GuessChallengeService) loadActiveChallenge(ctx context.Context, code string) (*model.GuessChallenge, *model.City, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" || utf8.RuneCountInString(code) > 16 {
		return nil, nil, invalidParam("challenge code is invalid")
	}
	challenge, err := s.repo.FindChallengeByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, notFound("challenge not found")
		}
		return nil, nil, fmt.Errorf("find guess challenge: %w", err)
	}
	if !challenge.ExpiresAt.After(s.now()) {
		return nil, nil, notFound("challenge expired")
	}
	city, err := s.cityRepo.FindByID(ctx, challenge.CityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, notFound("city not found")
		}
		return nil, nil, fmt.Errorf("find city: %w", err)
	}
	return challenge, city, nil
}

func (s *GuessChallengeService) resolveChallengeImage(dataURL, imageURL string, fallback *string) (*string, error) {
	dataURL = strings.TrimSpace(dataURL)
	if dataURL != "" {
		if s.storage == nil {
			return nil, invalidParam("image_data_url is not supported")
		}
		data, ext, err := parseImageDataURL(dataURL)
		if err != nil {
			return nil, err
		}
		url, err := s.storage.SaveBytes(data, "guess", ext)
		if err != nil {
			return nil, fmt.Errorf("save guess image: %w", err)
		}
		return &url, nil
	}

	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" && fallback != nil {
		imageURL = *fallback
	}
	if imageURL == "" {
		return nil, nil
	}
	if !strings.HasPrefix(imageURL, "/static/") && !strings.HasPrefix(imageURL, "/uploads/") {
		return nil, invalidParam("image_url must start with /static/ or /uploads/")
	}
	return &imageURL, nil
}

func (s *GuessChallengeService) challengeResponse(challenge *model.GuessChallenge, city *model.City) *GuessChallengeResponse {
	return &GuessChallengeResponse{
		Code:       challenge.Code,
		ShareURL:   "/?guess=" + challenge.Code,
		CityID:     challenge.CityID,
		CityName:   city.Name,
		TargetName: challenge.TargetName,
		ImageURL:   challenge.ImageURL,
		Caption:    challenge.Caption,
		CreatedAt:  challenge.CreatedAt,
		ExpiresAt:  challenge.ExpiresAt,
	}
}

func optionalText(field, value string, max int) (*string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(value) > max {
		return nil, invalidParam(fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return &value, nil
}

func parseImageDataURL(value string) ([]byte, string, error) {
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 || !strings.HasSuffix(parts[0], ";base64") {
		return nil, "", invalidParam("image_data_url must be a base64 data URL")
	}
	mimeType := strings.TrimPrefix(strings.TrimSuffix(parts[0], ";base64"), "data:")
	ext := ""
	switch mimeType {
	case "image/png":
		ext = ".png"
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/webp":
		ext = ".webp"
	default:
		return nil, "", invalidParam("image_data_url must be png, jpeg, or webp")
	}
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", invalidParam("image_data_url is invalid base64")
	}
	if len(data) == 0 || len(data) > maxGuessImageBytes {
		return nil, "", invalidParam("image_data_url must be between 1 byte and 5MB")
	}
	return data, ext, nil
}

func newGuessCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, guessChallengeCodeLen)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%08X", time.Now().UnixNano())
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return string(buf)
}

func guessAnswerMatches(answer, cityName, targetName string) bool {
	answer = strings.TrimSpace(strings.ToLower(answer))
	cityName = strings.TrimSpace(strings.ToLower(cityName))
	targetName = strings.TrimSpace(strings.ToLower(targetName))
	if answer == "" {
		return false
	}
	return answer == cityName ||
		(targetName != "" && answer == targetName) ||
		strings.Contains(answer, cityName) ||
		(targetName != "" && strings.Contains(answer, targetName))
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
