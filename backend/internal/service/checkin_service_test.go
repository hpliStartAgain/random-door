package service

import (
	"context"
	"errors"
	"testing"

	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
)

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

type fakeCheckinRepo struct {
	create func(context.Context, *model.Checkin) error
}

func (r *fakeCheckinRepo) Create(ctx context.Context, c *model.Checkin) error {
	return r.create(ctx, c)
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

func newTestCheckinService(cityRepo checkinCityRepo, cr checkinRepo, img imageGenerator, st imageStorage) *CheckinService {
	return &CheckinService{cityRepo: cityRepo, checkinRepo: cr, imageClient: img, storage: st}
}

func TestCheckinService_GenerateImage_Success(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		&fakeCheckinRepo{create: func(_ context.Context, _ *model.Checkin) error { return nil }},
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
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		&fakeCheckinRepo{},
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 1, 999, "/tmp/selfie.jpg")
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("expected ErrInvalidParam, got %v", err)
	}
}

func TestCheckinService_GenerateImage_AITimeout(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		&fakeCheckinRepo{},
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
		&fakeCheckinCityRepo{},
		&fakeCheckinRepo{create: func(_ context.Context, c *model.Checkin) error {
			c.ID = 42
			return nil
		}},
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	res, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CheckinID != 42 {
		t.Fatalf("expected checkin_id=42, got %d", res.CheckinID)
	}
}

func TestCheckinService_Create_DBError(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinCityRepo{},
		&fakeCheckinRepo{create: func(_ context.Context, _ *model.Checkin) error { return errors.New("db error") }},
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
