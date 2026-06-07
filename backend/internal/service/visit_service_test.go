package service

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

const testAnonymousID = "550e8400-e29b-41d4-a716-446655440000"

type fakeUserRepository struct {
	findByAnonymousID func(context.Context, string) (*model.User, error)
	findByID          func(context.Context, int64) (*model.User, error)
	create            func(context.Context, *model.User) error
}

func (r *fakeUserRepository) FindByAnonymousID(ctx context.Context, anonymousID string) (*model.User, error) {
	return r.findByAnonymousID(ctx, anonymousID)
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByID(ctx, id)
}

func (r *fakeUserRepository) Create(ctx context.Context, user *model.User) error {
	return r.create(ctx, user)
}

type fakeCityFinder struct {
	findByID func(context.Context, int64) (*model.City, error)
}

func (r *fakeCityFinder) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}

type fakeVisitRepository struct {
	create func(context.Context, *model.CityVisit) error
}

func (r *fakeVisitRepository) Create(ctx context.Context, visit *model.CityVisit) error {
	return r.create(ctx, visit)
}

func TestCreateAnonymousUser(t *testing.T) {
	t.Run("restores existing user", func(t *testing.T) {
		want := &model.User{ID: 7, AnonymousID: testAnonymousID}
		svc := NewVisitService(&fakeUserRepository{
			findByAnonymousID: func(context.Context, string) (*model.User, error) { return want, nil },
			findByID:          unusedFindUserByID,
			create:            unusedCreateUser,
		}, unusedCityFinder(), unusedVisitRepository())

		got, err := svc.CreateAnonymousUser(context.Background(), testAnonymousID)
		if err != nil || got != want {
			t.Fatalf("CreateAnonymousUser() = %#v, %v; want existing user", got, err)
		}
	})

	t.Run("creates normalized user", func(t *testing.T) {
		var created *model.User
		svc := NewVisitService(&fakeUserRepository{
			findByAnonymousID: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
			findByID:          unusedFindUserByID,
			create: func(_ context.Context, user *model.User) error {
				created = user
				user.ID = 8
				return nil
			},
		}, unusedCityFinder(), unusedVisitRepository())

		got, err := svc.CreateAnonymousUser(context.Background(), "550E8400-E29B-41D4-A716-446655440000")
		if err != nil || got.ID != 8 || created.AnonymousID != testAnonymousID {
			t.Fatalf("CreateAnonymousUser() = %#v, %v; want normalized created user", got, err)
		}
	})

	t.Run("recovers concurrent insert", func(t *testing.T) {
		findCalls := 0
		want := &model.User{ID: 9, AnonymousID: testAnonymousID}
		svc := NewVisitService(&fakeUserRepository{
			findByAnonymousID: func(context.Context, string) (*model.User, error) {
				findCalls++
				if findCalls == 1 {
					return nil, gorm.ErrRecordNotFound
				}
				return want, nil
			},
			findByID: unusedFindUserByID,
			create:   func(context.Context, *model.User) error { return errors.New("duplicate key") },
		}, unusedCityFinder(), unusedVisitRepository())

		got, err := svc.CreateAnonymousUser(context.Background(), testAnonymousID)
		if err != nil || got != want || findCalls != 2 {
			t.Fatalf("CreateAnonymousUser() = %#v, %v with %d finds; want concurrent result", got, err, findCalls)
		}
	})

	for _, anonymousID := range []string{"", "not-a-uuid", " 550e8400-e29b-41d4-a716-446655440000", "{550e8400-e29b-41d4-a716-446655440000}"} {
		t.Run("rejects "+anonymousID, func(t *testing.T) {
			svc := NewVisitService(nil, nil, nil)
			_, err := svc.CreateAnonymousUser(context.Background(), anonymousID)
			if !errors.Is(err, ErrInvalidParam) {
				t.Fatalf("CreateAnonymousUser(%q) error = %v, want ErrInvalidParam", anonymousID, err)
			}
		})
	}
}

func TestCreateFreeVisit(t *testing.T) {
	t.Run("creates default map click visit", func(t *testing.T) {
		var created *model.CityVisit
		svc := newFreeVisitService(
			func(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil },
			func(context.Context, int64) (*model.City, error) { return &model.City{ID: 3}, nil },
			func(_ context.Context, visit *model.CityVisit) error {
				created = visit
				visit.ID = 1001
				return nil
			},
		)

		got, err := svc.CreateFreeVisit(context.Background(), 1, 3, "")
		if err != nil || got.Visit.ID != 1001 || created.VisitMode != "free" || created.Source == nil || *created.Source != "map_click" || created.DiceRollID != nil {
			t.Fatalf("CreateFreeVisit() = %#v, %v; want free map_click visit", got, err)
		}
	})

	for _, tc := range []struct {
		name   string
		userID int64
		cityID int64
		source string
	}{
		{name: "invalid user", userID: 0, cityID: 3},
		{name: "invalid city", userID: 1, cityID: -1},
		{name: "invalid source", userID: 1, cityID: 3, source: "dice_roll"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewVisitService(nil, nil, nil)
			_, err := svc.CreateFreeVisit(context.Background(), tc.userID, tc.cityID, tc.source)
			if !errors.Is(err, ErrInvalidParam) {
				t.Fatalf("CreateFreeVisit() error = %v, want ErrInvalidParam", err)
			}
		})
	}

	t.Run("classifies missing user", func(t *testing.T) {
		svc := newFreeVisitService(
			func(context.Context, int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
			unusedFindCityByID,
			unusedCreateVisit,
		)
		_, err := svc.CreateFreeVisit(context.Background(), 1, 3, "search")
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "user not found" {
			t.Fatalf("CreateFreeVisit() error = %v, want user not found", err)
		}
	})

	t.Run("classifies missing city", func(t *testing.T) {
		svc := newFreeVisitService(
			func(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil },
			func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
			unusedCreateVisit,
		)
		_, err := svc.CreateFreeVisit(context.Background(), 1, 3, "search")
		if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "city not found" {
			t.Fatalf("CreateFreeVisit() error = %v, want city not found", err)
		}
	})
}

func newFreeVisitService(
	findUserByID func(context.Context, int64) (*model.User, error),
	findCityByID func(context.Context, int64) (*model.City, error),
	createVisit func(context.Context, *model.CityVisit) error,
) *VisitService {
	return NewVisitService(&fakeUserRepository{
		findByAnonymousID: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		findByID:          findUserByID,
		create:            unusedCreateUser,
	}, &fakeCityFinder{findByID: findCityByID}, &fakeVisitRepository{create: createVisit})
}

func unusedFindUserByID(context.Context, int64) (*model.User, error) {
	panic("unexpected FindByID call")
}

func unusedCreateUser(context.Context, *model.User) error {
	panic("unexpected Create call")
}

func unusedCityFinder() *fakeCityFinder {
	return &fakeCityFinder{findByID: unusedFindCityByID}
}

func unusedFindCityByID(context.Context, int64) (*model.City, error) {
	panic("unexpected FindByID call")
}

func unusedVisitRepository() *fakeVisitRepository {
	return &fakeVisitRepository{create: unusedCreateVisit}
}

func unusedCreateVisit(context.Context, *model.CityVisit) error {
	panic("unexpected Create call")
}
