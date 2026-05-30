package service

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

// seqSource returns a fixed queue of Intn results for deterministic rolls.
type seqSource struct {
	values []int
	idx    int
}

func (s *seqSource) Intn(int) int {
	v := s.values[s.idx%len(s.values)]
	s.idx++
	return v
}

type fakeGameUserFinder struct {
	findByID func(context.Context, int64) (*model.User, error)
}

func (r *fakeGameUserFinder) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByID(ctx, id)
}

type fakeGameCityRepo struct {
	listAll  func(context.Context) ([]model.City, error)
	findByID func(context.Context, int64) (*model.City, error)
}

func (r *fakeGameCityRepo) ListAll(ctx context.Context) ([]model.City, error) {
	return r.listAll(ctx)
}

func (r *fakeGameCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}

type fakeGameVisitReader struct {
	listVisited func(context.Context, int64) ([]int64, error)
}

func (r *fakeGameVisitReader) ListVisitedCityIDs(ctx context.Context, userID int64) ([]int64, error) {
	return r.listVisited(ctx, userID)
}

type fakeGameStore struct {
	create func(context.Context, *model.DiceRoll, *model.CityVisit) error
}

func (r *fakeGameStore) CreateRollWithVisit(ctx context.Context, roll *model.DiceRoll, visit *model.CityVisit) error {
	return r.create(ctx, roll, visit)
}

var gameTestCities = []model.City{
	{ID: 1, Name: "北京", Province: "北京", Lat: 39.9042, Lng: 116.4074},
	{ID: 3, Name: "西安", Province: "陕西", Lat: 34.3416, Lng: 108.9398},
	{ID: 6, Name: "广州", Province: "广东", Lat: 23.1291, Lng: 113.2644},
}

func okUser(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil }

func TestGameInit(t *testing.T) {
	t.Run("defaults to Beijing and returns nearest city", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			&fakeGameCityRepo{listAll: func(context.Context) ([]model.City, error) { return gameTestCities, nil }},
			nil, nil,
		)
		got, err := svc.Init(context.Background(), 1, 0, 0)
		if err != nil || got.NearestCity.ID != 1 || got.NearestCity.Name != "北京" {
			t.Fatalf("Init() = %#v, %v; want Beijing", got, err)
		}
	})

	t.Run("matches nearest city to provided coords", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			&fakeGameCityRepo{listAll: func(context.Context) ([]model.City, error) { return gameTestCities, nil }},
			nil, nil,
		)
		got, err := svc.Init(context.Background(), 1, 34.3, 108.9)
		if err != nil || got.NearestCity.ID != 3 {
			t.Fatalf("Init() = %#v, %v; want 西安", got, err)
		}
	})

	for _, tc := range []struct {
		name     string
		userID   int64
		lat, lng float64
	}{
		{name: "invalid user", userID: 0, lat: 39, lng: 116},
		{name: "invalid lat", userID: 1, lat: 200, lng: 116},
		{name: "invalid lng", userID: 1, lat: 39, lng: 999},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewGameService(&fakeGameUserFinder{findByID: okUser}, nil, nil, nil)
			_, err := svc.Init(context.Background(), tc.userID, tc.lat, tc.lng)
			if !errors.Is(err, ErrInvalidParam) {
				t.Fatalf("Init() error = %v, want ErrInvalidParam", err)
			}
		})
	}

	t.Run("classifies missing user", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: func(context.Context, int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound }},
			nil, nil, nil,
		)
		_, err := svc.Init(context.Background(), 1, 39, 116)
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "user not found" {
			t.Fatalf("Init() error = %v, want user not found", err)
		}
	})

	t.Run("empty cities surface as internal error", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			&fakeGameCityRepo{listAll: func(context.Context) ([]model.City, error) { return nil, nil }},
			nil, nil,
		)
		_, err := svc.Init(context.Background(), 1, 39, 116)
		if err == nil || errors.Is(err, ErrInvalidParam) || errors.Is(err, ErrNotFound) {
			t.Fatalf("Init() error = %v, want unclassified internal error", err)
		}
	})
}

func TestGameRoll(t *testing.T) {
	cityRepo := &fakeGameCityRepo{
		listAll:  func(context.Context) ([]model.City, error) { return gameTestCities, nil },
		findByID: func(context.Context, int64) (*model.City, error) { return &model.City{ID: 1}, nil },
	}

	t.Run("persists deterministic roll and game visit", func(t *testing.T) {
		var savedRoll *model.DiceRoll
		var savedVisit *model.CityVisit
		store := &fakeGameStore{create: func(_ context.Context, roll *model.DiceRoll, visit *model.CityVisit) error {
			roll.ID = 2001
			visit.DiceRollID = &roll.ID
			visit.ID = 1002
			savedRoll, savedVisit = roll, visit
			return nil
		}}
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			cityRepo,
			&fakeGameVisitReader{listVisited: func(context.Context, int64) ([]int64, error) { return nil, nil }},
			store,
		).WithRandSource(&seqSource{values: []int{5, 4}}) // 西南, 800km

		got, err := svc.Roll(context.Background(), 1, 1, 39.9042, 116.4074)
		if err != nil {
			t.Fatalf("Roll() error = %v", err)
		}
		if got.Direction != "西南" || got.DistanceKm != 800 {
			t.Fatalf("Roll() dir/dist = %s/%d, want 西南/800", got.Direction, got.DistanceKm)
		}
		if got.VisitID != 1002 || got.DiceRollID != 2001 || got.TargetCity.ID == 1 {
			t.Fatalf("Roll() = %#v, want persisted ids and a city other than the origin", got)
		}
		if savedRoll.ToCityID != got.TargetCity.ID || savedRoll.FromCityID == nil || *savedRoll.FromCityID != 1 {
			t.Fatalf("saved roll = %#v, want linked from/to cities", savedRoll)
		}
		if savedVisit.VisitMode != "game" || savedVisit.Source == nil || *savedVisit.Source != "dice_roll" ||
			savedVisit.DiceRollID == nil || *savedVisit.DiceRollID != 2001 {
			t.Fatalf("saved visit = %#v, want game dice_roll visit linked to roll", savedVisit)
		}
	})

	for _, tc := range []struct {
		name               string
		userID, fromCityID int64
		lat, lng           float64
	}{
		{name: "invalid user", userID: 0, fromCityID: 1, lat: 39, lng: 116},
		{name: "invalid from city", userID: 1, fromCityID: 0, lat: 39, lng: 116},
		{name: "invalid lat", userID: 1, fromCityID: 1, lat: 91, lng: 116},
		{name: "invalid lng", userID: 1, fromCityID: 1, lat: 39, lng: -181},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewGameService(&fakeGameUserFinder{findByID: okUser}, cityRepo, nil, nil)
			_, err := svc.Roll(context.Background(), tc.userID, tc.fromCityID, tc.lat, tc.lng)
			if !errors.Is(err, ErrInvalidParam) {
				t.Fatalf("Roll() error = %v, want ErrInvalidParam", err)
			}
		})
	}

	t.Run("classifies missing user", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: func(context.Context, int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound }},
			cityRepo, nil, nil,
		)
		_, err := svc.Roll(context.Background(), 1, 1, 39, 116)
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "user not found" {
			t.Fatalf("Roll() error = %v, want user not found", err)
		}
	})

	t.Run("classifies missing from city", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			&fakeGameCityRepo{findByID: func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound }},
			nil, nil,
		)
		_, err := svc.Roll(context.Background(), 1, 99, 39, 116)
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "from city not found" {
			t.Fatalf("Roll() error = %v, want from city not found", err)
		}
	})

	t.Run("store failure surfaces as internal error", func(t *testing.T) {
		store := &fakeGameStore{create: func(context.Context, *model.DiceRoll, *model.CityVisit) error {
			return errors.New("mysql password secret")
		}}
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			cityRepo,
			&fakeGameVisitReader{listVisited: func(context.Context, int64) ([]int64, error) { return nil, nil }},
			store,
		).WithRandSource(&seqSource{values: []int{0, 0}})
		_, err := svc.Roll(context.Background(), 1, 1, 39, 116)
		if err == nil || errors.Is(err, ErrInvalidParam) || errors.Is(err, ErrNotFound) {
			t.Fatalf("Roll() error = %v, want unclassified internal error", err)
		}
	})

	t.Run("empty cities surface as internal error", func(t *testing.T) {
		svc := NewGameService(
			&fakeGameUserFinder{findByID: okUser},
			&fakeGameCityRepo{
				listAll:  func(context.Context) ([]model.City, error) { return nil, nil },
				findByID: func(context.Context, int64) (*model.City, error) { return &model.City{ID: 1}, nil },
			},
			&fakeGameVisitReader{listVisited: func(context.Context, int64) ([]int64, error) { return nil, nil }},
			nil,
		).WithRandSource(&seqSource{values: []int{0, 0}})
		_, err := svc.Roll(context.Background(), 1, 1, 39, 116)
		if err == nil || errors.Is(err, ErrInvalidParam) || errors.Is(err, ErrNotFound) {
			t.Fatalf("Roll() error = %v, want unclassified internal error", err)
		}
	})
}
