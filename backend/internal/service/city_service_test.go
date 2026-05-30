package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeCityRepository struct {
	listAll        func(context.Context) ([]model.City, error)
	findByID       func(context.Context, int64) (*model.City, error)
	listTags       func(context.Context, int64) ([]model.CityTag, error)
	listLandmarks  func(context.Context, int64) ([]model.Landmark, error)
	listFoods      func(context.Context, int64) ([]model.Food, error)
	listCharacters func(context.Context, int64) ([]model.Character, error)
}

func (r *fakeCityRepository) ListAll(ctx context.Context) ([]model.City, error) {
	return r.listAll(ctx)
}

func (r *fakeCityRepository) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}

func (r *fakeCityRepository) ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error) {
	return r.listTags(ctx, cityID)
}

func (r *fakeCityRepository) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	return r.listLandmarks(ctx, cityID)
}

func (r *fakeCityRepository) ListFoods(ctx context.Context, cityID int64) ([]model.Food, error) {
	return r.listFoods(ctx, cityID)
}

func (r *fakeCityRepository) ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error) {
	return r.listCharacters(ctx, cityID)
}

func TestCityServiceList(t *testing.T) {
	repo := detailCityRepository()
	repo.listAll = func(context.Context) ([]model.City, error) {
		return []model.City{{ID: 1, Name: "北京", Province: "北京", Lat: 39.9042, Lng: 116.4074}}, nil
	}

	cities, err := NewCityService(repo).List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(cities) != 1 || cities[0].Name != "北京" || len(cities[0].Tags) != 1 || cities[0].Tags[0] != "ancient_capital" {
		t.Fatalf("List() = %#v, want city with tags", cities)
	}
}

func TestCityServiceListPropagatesTagError(t *testing.T) {
	repo := detailCityRepository()
	repo.listAll = func(context.Context) ([]model.City, error) {
		return []model.City{{ID: 1}}, nil
	}
	repo.listTags = func(context.Context, int64) ([]model.CityTag, error) {
		return nil, errors.New("database unavailable")
	}

	_, err := NewCityService(repo).List(context.Background())
	if err == nil || !strings.Contains(err.Error(), "list tags for city 1") {
		t.Fatalf("List() error = %v, want wrapped tag error", err)
	}
}

func TestCityServiceDetail(t *testing.T) {
	t.Run("success omits sensitive character fields", func(t *testing.T) {
		detail, err := NewCityService(detailCityRepository()).Detail(context.Background(), 1)
		if err != nil {
			t.Fatalf("Detail() error = %v", err)
		}
		payload, err := json.Marshal(detail)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}
		if strings.Contains(string(payload), "sensitive") {
			t.Fatalf("Detail() leaked sensitive character data: %s", payload)
		}
		if len(detail.Landmarks) != 1 || len(detail.Foods) != 1 || len(detail.Characters) != 1 {
			t.Fatalf("Detail() = %#v, want all related content", detail)
		}
	})

	t.Run("rejects non-positive id", func(t *testing.T) {
		_, err := NewCityService(detailCityRepository()).Detail(context.Background(), 0)
		if !errors.Is(err, ErrInvalidParam) {
			t.Fatalf("Detail() error = %v, want ErrInvalidParam", err)
		}
	})

	t.Run("classifies missing city", func(t *testing.T) {
		repo := detailCityRepository()
		repo.findByID = func(context.Context, int64) (*model.City, error) {
			return nil, gorm.ErrRecordNotFound
		}
		_, err := NewCityService(repo).Detail(context.Background(), 404)
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "city not found" {
			t.Fatalf("Detail() error = %v, want city not found", err)
		}
	})

	for _, relation := range []string{"tags", "landmarks", "foods", "characters"} {
		t.Run("propagates "+relation+" error", func(t *testing.T) {
			repo := detailCityRepository()
			wantErr := errors.New("database unavailable")
			switch relation {
			case "tags":
				repo.listTags = func(context.Context, int64) ([]model.CityTag, error) { return nil, wantErr }
			case "landmarks":
				repo.listLandmarks = func(context.Context, int64) ([]model.Landmark, error) { return nil, wantErr }
			case "foods":
				repo.listFoods = func(context.Context, int64) ([]model.Food, error) { return nil, wantErr }
			case "characters":
				repo.listCharacters = func(context.Context, int64) ([]model.Character, error) { return nil, wantErr }
			}
			_, err := NewCityService(repo).Detail(context.Background(), 1)
			if !errors.Is(err, wantErr) {
				t.Fatalf("Detail() error = %v, want wrapped relation error", err)
			}
		})
	}
}

func detailCityRepository() *fakeCityRepository {
	imageURL := "/static/image.jpg"
	description := "description"
	style := "京腔"
	return &fakeCityRepository{
		listAll: func(context.Context) ([]model.City, error) {
			return nil, nil
		},
		findByID: func(context.Context, int64) (*model.City, error) {
			return &model.City{ID: 1, Name: "北京", Province: "北京", Lat: 39.9042, Lng: 116.4074}, nil
		},
		listTags: func(context.Context, int64) ([]model.CityTag, error) {
			return []model.CityTag{{CityID: 1, Tag: "ancient_capital"}}, nil
		},
		listLandmarks: func(context.Context, int64) ([]model.Landmark, error) {
			return []model.Landmark{{ID: 1, CityID: 1, Name: "故宫", ImageURL: &imageURL, Description: &description}}, nil
		},
		listFoods: func(context.Context, int64) ([]model.Food, error) {
			return []model.Food{{ID: 1, CityID: 1, Name: "烤鸭", ImageURL: &imageURL, Description: &description}}, nil
		},
		listCharacters: func(context.Context, int64) ([]model.Character, error) {
			return []model.Character{{
				ID: 1, CityID: 1, Name: "老舍风格向导", CharacterType: "culture",
				AvatarURL: &imageURL, Persona: "sensitive persona", DialectStyle: &style, Prompt: "sensitive prompt",
			}}, nil
		},
	}
}
