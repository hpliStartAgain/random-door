package service

import (
	"context"
	"errors"
	"testing"

	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeCheckinUserFinder struct {
	findByID func(context.Context, int64) (*model.User, error)
}

func (r *fakeCheckinUserFinder) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByID(ctx, id)
}

type fakeCheckinCityRepo struct {
	findByID      func(context.Context, int64) (*model.City, error)
	listLandmarks func(context.Context, int64) ([]model.Landmark, error)
}

func (r *fakeCheckinCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}
func (r *fakeCheckinCityRepo) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	return r.listLandmarks(ctx, cityID)
}

type fakeCheckinStore struct {
	createAndEvaluate func(context.Context, *model.Checkin) ([]model.Achievement, error)
}

func (r *fakeCheckinStore) CreateAndEvaluate(ctx context.Context, c *model.Checkin) ([]model.Achievement, error) {
	return r.createAndEvaluate(ctx, c)
}

type fakeImageGen struct {
	generate func(context.Context, string, string, string) ([]byte, error)
}

func (f *fakeImageGen) Generate(ctx context.Context, selfie, ref, prompt string) ([]byte, error) {
	return f.generate(ctx, selfie, ref, prompt)
}

type fakeImageStorage struct {
	saveBytes func([]byte, string, string) (string, error)
}

func (f *fakeImageStorage) SaveBytes(data []byte, subDir, ext string) (string, error) {
	return f.saveBytes(data, subDir, ext)
}

var lm1ID int64 = 10
var lm1 = model.Landmark{ID: lm1ID, CityID: 1, Name: "大雁塔"}

func okCheckinUser(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil }

func newTestCheckinService(userRepo checkinUserFinder, cityRepo checkinCityRepo, store checkinStore, img imageGenerator, st imageStorage) *CheckinService {
	return &CheckinService{userRepo: userRepo, cityRepo: cityRepo, store: store, imageClient: img, storage: st}
}

func TestCheckinService_GenerateImage_Success(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{generate: func(_ context.Context, _, _, _ string) ([]byte, error) { return []byte("img"), nil }},
		&fakeImageStorage{saveBytes: func(_ []byte, _, _ string) (string, error) { return "/uploads/generated/x.png", nil }},
	)
	res, err := svc.GenerateImage(context.Background(), 1, 1, lm1ID, "/tmp/selfie.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GeneratedImageURL != "/uploads/generated/x.png" {
		t.Fatalf("unexpected url: %s", res.GeneratedImageURL)
	}
}

func TestCheckinService_GenerateImage_LandmarkNotFound(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 1, 999, "/tmp/selfie.jpg")
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("expected ErrInvalidParam, got %v", err)
	}
}

func TestCheckinService_GenerateImage_CityNotFound(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID: func(_ context.Context, _ int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
		},
		nil,
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 999, lm1ID, "/tmp/selfie.jpg")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCheckinService_GenerateImage_AITimeout(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{generate: func(_ context.Context, _, _, _ string) ([]byte, error) { return nil, aiPkg.ErrAITimeout }},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 1, lm1ID, "/tmp/selfie.jpg")
	if !errors.Is(err, aiPkg.ErrAITimeout) {
		t.Fatalf("expected ErrAITimeout, got %v", err)
	}
}

func TestCheckinService_Create_Success(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, c *model.Checkin) ([]model.Achievement, error) {
			c.ID = 42
			return nil, nil
		}},
		nil, nil,
	)
	res, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CheckinID != 42 {
		t.Fatalf("expected checkin_id=42, got %d", res.CheckinID)
	}
}

func TestCheckinService_Create_StoreError(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, _ *model.Checkin) ([]model.Achievement, error) {
			return nil, errors.New("db error")
		}},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckinService_Create_UserNotFound(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: func(_ context.Context, _ int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound }},
		&fakeCheckinCityRepo{},
		&fakeCheckinStore{},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 999, CityID: 1})
	if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestCheckinService_Create_CityNotFound(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound }},
		&fakeCheckinStore{},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 999})
	if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "city not found" {
		t.Fatalf("expected city not found, got %v", err)
	}
}
