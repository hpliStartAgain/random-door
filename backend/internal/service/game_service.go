package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/your-org/city-roam/backend/internal/geo"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

// GameCityRepository exposes the city reads the game flow needs.
type GameCityRepository interface {
	ListAll(ctx context.Context) ([]model.City, error)
	FindByID(ctx context.Context, id int64) (*model.City, error)
}

// GameUserFinder verifies that the acting user exists.
type GameUserFinder interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

// GameVisitReader reports which cities a user has already visited.
type GameVisitReader interface {
	ListVisitedCityIDs(ctx context.Context, userID int64) ([]int64, error)
}

// GameStore persists a dice roll and its resulting visit atomically.
type GameStore interface {
	CreateRollWithVisit(ctx context.Context, roll *model.DiceRoll, visit *model.CityVisit) error
}

type GameService struct {
	userRepo     GameUserFinder
	cityRepo     GameCityRepository
	visitRepo    GameVisitReader
	store        GameStore
	achievements visitAchievementEvaluator
	rng          geo.IntnSource
}

func NewGameService(userRepo GameUserFinder, cityRepo GameCityRepository, visitRepo GameVisitReader, store GameStore) *GameService {
	return &GameService{userRepo: userRepo, cityRepo: cityRepo, visitRepo: visitRepo, store: store}
}

func (s *GameService) WithAchievementEvaluator(evaluator visitAchievementEvaluator) *GameService {
	s.achievements = evaluator
	return s
}

// WithRandSource overrides the random source; used by tests for deterministic rolls.
func (s *GameService) WithRandSource(rng geo.IntnSource) *GameService {
	s.rng = rng
	return s
}

func (s *GameService) directionSource() geo.IntnSource {
	if s.rng != nil {
		return s.rng
	}
	return nil
}

func validateCoords(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return invalidParam("lat must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return invalidParam("lng must be between -180 and 180")
	}
	return nil
}

func (s *GameService) ensureUser(ctx context.Context, userID int64) error {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFound("user not found")
		}
		return fmt.Errorf("find user: %w", err)
	}
	return nil
}

// InitResult contains the nearest city for game start.
type InitResult struct {
	NearestCity CityBrief `json:"nearest_city"`
}

type CityBrief struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Province string  `json:"province"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

// Init finds the nearest city to the given coordinates as the game starting point.
func (s *GameService) Init(ctx context.Context, userID int64, lat, lng float64) (*InitResult, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	// Default to Beijing if no coords provided (lat/lng are optional per contract).
	if lat == 0 && lng == 0 {
		lat, lng = 39.9042, 116.4074
	} else if err := validateCoords(lat, lng); err != nil {
		return nil, err
	}

	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}

	cities, err := s.cityRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	points := make([]geo.CityPoint, 0, len(cities))
	cityMap := make(map[int64]model.City)
	for _, c := range cities {
		points = append(points, geo.CityPoint{ID: c.ID, Lat: c.Lat, Lng: c.Lng})
		cityMap[c.ID] = c
	}

	matched, err := geo.MatchNearestCity(points, lat, lng, geo.MatchOptions{})
	if err != nil {
		return nil, fmt.Errorf("match city: %w", err)
	}

	city := cityMap[matched.ID]
	return &InitResult{
		NearestCity: CityBrief{
			ID: city.ID, Name: city.Name, Province: city.Province,
			Lat: city.Lat, Lng: city.Lng,
		},
	}, nil
}

// RollResult contains all data from a dice roll.
type RollResult struct {
	VisitID              int64              `json:"visit_id"`
	DiceRollID           int64              `json:"dice_roll_id"`
	Direction            string             `json:"direction"`
	DistanceKm           int                `json:"distance_km"`
	UnlockedAchievements []AchievementBrief `json:"unlocked_achievements"`
	TargetPoint          struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"target_point"`
	TargetCity CityBrief `json:"target_city"`
}

// Roll executes a dice roll: random direction + distance -> target point -> nearest city.
func (s *GameService) Roll(ctx context.Context, userID, fromCityID int64, lat, lng float64) (*RollResult, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	if fromCityID <= 0 {
		return nil, invalidParam("from_city_id must be a positive integer")
	}
	if err := validateCoords(lat, lng); err != nil {
		return nil, err
	}

	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := s.cityRepo.FindByID(ctx, fromCityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("from city not found")
		}
		return nil, fmt.Errorf("find from city: %w", err)
	}

	// 1. Random direction & distance (deterministic when a source is injected for tests).
	var dir geo.Direction
	var dist int
	if src := s.directionSource(); src != nil {
		dir = geo.RandomDirectionWithRand(src)
		dist = geo.RandomDistanceWithRand(src)
	} else {
		dir = geo.RandomDirection()
		dist = geo.RandomDistance()
	}

	// 2. Calculate target point
	tLat, tLng := geo.TargetPoint(lat, lng, dir.Bearing, dist)

	slog.Info("dice roll", "user_id", userID, "direction", dir.Name, "distance_km", dist,
		"target_lat", tLat, "target_lng", tLng)

	// 3. Load cities and find visited
	cities, err := s.cityRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	visitedIDs, err := s.visitRepo.ListVisitedCityIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list visited cities: %w", err)
	}

	points := make([]geo.CityPoint, 0, len(cities))
	cityMap := make(map[int64]model.City)
	for _, c := range cities {
		points = append(points, geo.CityPoint{ID: c.ID, Lat: c.Lat, Lng: c.Lng})
		cityMap[c.ID] = c
	}

	// 4. Match nearest city (with fallback)
	matched, err := geo.MatchNearestCity(points, tLat, tLng, geo.MatchOptions{
		ExcludeCityID:  fromCityID,
		VisitedCityIDs: visitedIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("match city: %w", err)
	}

	targetCity := cityMap[matched.ID]

	// 5. Transaction: write dice_rolls + city_visits atomically.
	source := "dice_roll"
	diceRoll := model.DiceRoll{
		UserID:     userID,
		FromCityID: &fromCityID,
		ToCityID:   targetCity.ID,
		Direction:  dir.Name,
		DistanceKm: dist,
		TargetLat:  &tLat,
		TargetLng:  &tLng,
	}
	visit := model.CityVisit{
		UserID:     userID,
		CityID:     targetCity.ID,
		VisitMode:  "game",
		Source:     &source,
		FromCityID: &fromCityID,
	}
	if err := s.store.CreateRollWithVisit(ctx, &diceRoll, &visit); err != nil {
		return nil, fmt.Errorf("persist dice roll: %w", err)
	}
	var unlocked []AchievementBrief
	if s.achievements != nil {
		achievements, err := s.achievements.UnlockForUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("evaluate achievements: %w", err)
		}
		unlocked = briefAchievements(achievements)
	}

	slog.Info("dice roll completed", "user_id", userID, "to_city", targetCity.Name,
		"dice_roll_id", diceRoll.ID, "visit_id", visit.ID)

	return &RollResult{
		VisitID:              visit.ID,
		DiceRollID:           diceRoll.ID,
		Direction:            dir.Name,
		DistanceKm:           dist,
		UnlockedAchievements: unlocked,
		TargetPoint: struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{Lat: tLat, Lng: tLng},
		TargetCity: CityBrief{
			ID: targetCity.ID, Name: targetCity.Name, Province: targetCity.Province,
			Lat: targetCity.Lat, Lng: targetCity.Lng,
		},
	}, nil
}
