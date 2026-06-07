package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeGuessChallengeRepo struct {
	challenge *model.GuessChallenge
	answer    *model.GuessAnswer
}

func (r *fakeGuessChallengeRepo) CreateChallenge(_ context.Context, challenge *model.GuessChallenge) error {
	challenge.ID = 1
	challenge.CreatedAt = time.Date(2026, 6, 7, 9, 0, 0, 0, time.UTC)
	r.challenge = challenge
	return nil
}

func (r *fakeGuessChallengeRepo) FindChallengeByCode(_ context.Context, code string) (*model.GuessChallenge, error) {
	if r.challenge == nil || r.challenge.Code != code {
		return nil, gorm.ErrRecordNotFound
	}
	return r.challenge, nil
}

func (r *fakeGuessChallengeRepo) CreateAnswer(_ context.Context, answer *model.GuessAnswer) error {
	answer.ID = 2
	r.answer = answer
	return nil
}

type fakeGuessChallengeCityRepo struct {
	city *model.City
	err  error
}

func (r *fakeGuessChallengeCityRepo) FindByID(context.Context, int64) (*model.City, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.city, nil
}

type fakeGuessStorage struct {
	url string
}

func (s *fakeGuessStorage) SaveBytes([]byte, string, string) (string, error) {
	return s.url, nil
}

func TestGuessChallengeServiceCreateAndAnswer(t *testing.T) {
	cover := "/static/landmarks/xian.png"
	repo := &fakeGuessChallengeRepo{}
	svc := NewGuessChallengeService(repo, &fakeGuessChallengeCityRepo{
		city: &model.City{ID: 3, Name: "西安", CoverImageURL: &cover},
	}, &fakeGuessStorage{url: "/uploads/guess/shot.png"})
	svc.now = func() time.Time { return time.Date(2026, 6, 7, 9, 0, 0, 0, time.UTC) }

	got, err := svc.Create(context.Background(), CreateGuessChallengeRequest{
		CityID:       3,
		TargetName:   "兵马俑",
		ImageDataURL: "data:image/png;base64,AAAA",
		Caption:      "猜猜我在哪",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.Code == "" || got.ShareURL == "" || got.ImageURL == nil || *got.ImageURL != "/uploads/guess/shot.png" {
		t.Fatalf("Create() = %#v, want share metadata", got)
	}

	answer, err := svc.Answer(context.Background(), repo.challenge.Code, "西安")
	if err != nil {
		t.Fatalf("Answer() error = %v", err)
	}
	if !answer.IsCorrect || repo.answer == nil || !repo.answer.IsCorrect {
		t.Fatalf("Answer() = %#v, stored=%#v, want correct", answer, repo.answer)
	}
}

func TestGuessChallengeServiceExpired(t *testing.T) {
	repo := &fakeGuessChallengeRepo{challenge: &model.GuessChallenge{
		Code:      "ABCDEFGH",
		CityID:    3,
		ExpiresAt: time.Date(2026, 6, 7, 8, 0, 0, 0, time.UTC),
	}}
	svc := NewGuessChallengeService(repo, &fakeGuessChallengeCityRepo{city: &model.City{ID: 3, Name: "西安"}}, nil)
	svc.now = func() time.Time { return time.Date(2026, 6, 7, 9, 0, 0, 0, time.UTC) }

	_, err := svc.Get(context.Background(), "ABCDEFGH")
	if !errors.Is(err, ErrNotFound) || !strings.Contains(ClientMessage(err), "expired") {
		t.Fatalf("Get() error = %v, want expired not found", err)
	}
}

func TestGuessChallengeServiceInvalidImageURL(t *testing.T) {
	svc := NewGuessChallengeService(&fakeGuessChallengeRepo{}, &fakeGuessChallengeCityRepo{
		city: &model.City{ID: 3, Name: "西安"},
	}, nil)
	_, err := svc.Create(context.Background(), CreateGuessChallengeRequest{CityID: 3, ImageURL: "https://example.com/a.png"})
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("Create() error = %v, want ErrInvalidParam", err)
	}
}

func TestGuessAnswerMatchesRequiresMeaningfulHit(t *testing.T) {
	if !guessAnswerMatches("我猜这里是西安", "西安", "兵马俑") {
		t.Fatal("expected phrase containing city name to match")
	}
	if !guessAnswerMatches("兵马俑", "西安", "兵马俑") {
		t.Fatal("expected exact target name to match")
	}
	if guessAnswerMatches("西", "西安", "兵马俑") {
		t.Fatal("single-character partial city answer should not match")
	}
}
