package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/service"
)

type handlerGuessChallengeRepo struct {
	challenge *model.GuessChallenge
}

func (r *handlerGuessChallengeRepo) CreateChallenge(_ context.Context, challenge *model.GuessChallenge) error {
	challenge.ID = 1
	challenge.CreatedAt = time.Date(2026, 6, 7, 9, 0, 0, 0, time.UTC)
	r.challenge = challenge
	return nil
}

func (r *handlerGuessChallengeRepo) FindChallengeByCode(_ context.Context, code string) (*model.GuessChallenge, error) {
	return r.challenge, nil
}

func (r *handlerGuessChallengeRepo) CreateAnswer(_ context.Context, answer *model.GuessAnswer) error {
	answer.ID = 2
	return nil
}

type handlerGuessCityRepo struct{}

func (handlerGuessCityRepo) FindByID(context.Context, int64) (*model.City, error) {
	cover := "/static/landmarks/xian.png"
	return &model.City{ID: 3, Name: "西安", CoverImageURL: &cover}, nil
}

func (handlerGuessCityRepo) ListLandmarks(context.Context, int64) ([]model.Landmark, error) {
	return nil, nil
}

func (handlerGuessCityRepo) ListFoods(context.Context, int64) ([]model.Food, error) {
	return nil, nil
}

func (handlerGuessCityRepo) ListCharacters(context.Context, int64) ([]model.Character, error) {
	return nil, nil
}

type handlerGuessStorage struct{}

func (handlerGuessStorage) SaveBytes([]byte, string, string) (string, error) {
	return "/uploads/guess/shot.png", nil
}

func TestGuessChallengeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerGuessChallengeRepo{}
	challengeSvc := service.NewGuessChallengeService(repo, handlerGuessCityRepo{}, handlerGuessStorage{})
	guessHandler := NewGuessHandler(service.NewGuessService(handlerGuessCityRepo{}, nil)).WithChallengeService(challengeSvc)
	router := gin.New()
	router.POST("/api/guess/challenges", guessHandler.CreateChallenge)
	router.GET("/api/guess/challenges/:code", guessHandler.GetChallenge)
	router.POST("/api/guess/challenges/:code/answers", guessHandler.AnswerChallenge)

	created := requestJSON(t, router, http.MethodPost, "/api/guess/challenges", `{"city_id":3,"target_name":"兵马俑","image_data_url":"data:image/png;base64,AAAA","caption":"猜猜我在哪"}`, http.StatusCreated)
	code, ok := created["code"].(string)
	if !ok || code == "" || created["share_url"] == "" {
		t.Fatalf("create challenge = %#v, want code/share_url", created)
	}

	got := requestJSON(t, router, http.MethodGet, "/api/guess/challenges/"+code, "", http.StatusOK)
	if got["city_name"] != "西安" {
		t.Fatalf("get challenge = %#v, want city name", got)
	}

	answer := requestJSON(t, router, http.MethodPost, "/api/guess/challenges/"+code+"/answers", `{"answer_text":"西安"}`, http.StatusCreated)
	if answer["is_correct"] != true {
		t.Fatalf("answer = %#v, want correct", answer)
	}
}
