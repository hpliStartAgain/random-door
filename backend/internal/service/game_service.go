package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/your-org/city-roam/backend/internal/geo"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/repository"
	"gorm.io/gorm"
)

type GameService struct {
	db        *gorm.DB
	cityRepo  *repository.CityRepo
	visitRepo *repository.VisitRepo
	diceRepo  *repository.DiceRepo
}

func NewGameService(db *gorm.DB, cityRepo *repository.CityRepo, visitRepo *repository.VisitRepo, diceRepo *repository.DiceRepo) *GameService {
	return &GameService{db: db, cityRepo: cityRepo, visitRepo: visitRepo, diceRepo: diceRepo}
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
	// Default to Beijing if no coords provided
	if lat == 0 && lng == 0 {
		lat, lng = 39.9042, 116.4074
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
	VisitID     int64  `json:"visit_id"`
	DiceRollID  int64  `json:"dice_roll_id"`
	Direction   string `json:"direction"`
	DistanceKm  int    `json:"distance_km"`
	TargetPoint struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"target_point"`
	TargetCity CityBrief `json:"target_city"`
}

// Roll executes a dice roll: random direction + distance -> target point -> nearest city.
func (s *GameService) Roll(ctx context.Context, userID, fromCityID int64, lat, lng float64) (*RollResult, error) {
	// 1. Random direction & distance
	dir := geo.RandomDirection()
	dist := geo.RandomDistance()

	// 2. Calculate target point
	tLat, tLng := geo.TargetPoint(lat, lng, dir.Bearing, dist)

	slog.Info("dice roll", "user_id", userID, "direction", dir.Name, "distance_km", dist,
		"target_lat", tLat, "target_lng", tLng)

	// 3. Load cities and find visited
	cities, err := s.cityRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	visitedIDs, _ := s.visitRepo.ListVisitedCityIDs(ctx, userID)

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
		DirectionDeg:   dir.Bearing,
	})
	if err != nil {
		return nil, fmt.Errorf("match city: %w", err)
	}

	targetCity := cityMap[matched.ID]

	// 5. Transaction: write dice_rolls + city_visits
	var diceRoll model.DiceRoll
	var visit model.CityVisit

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		diceRoll = model.DiceRoll{
			UserID:     userID,
			FromCityID: &fromCityID,
			ToCityID:   targetCity.ID,
			Direction:  dir.Name,
			DistanceKm: dist,
			TargetLat:  &tLat,
			TargetLng:  &tLng,
		}
		if err := s.diceRepo.CreateTx(tx, &diceRoll); err != nil {
			return fmt.Errorf("create dice_roll: %w", err)
		}

		source := "dice_roll"
		visit = model.CityVisit{
			UserID:     userID,
			CityID:     targetCity.ID,
			VisitMode:  "game",
			Source:     &source,
			FromCityID: &fromCityID,
			DiceRollID: &diceRoll.ID,
		}
		if err := s.visitRepo.CreateTx(tx, &visit); err != nil {
			return fmt.Errorf("create visit: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	slog.Info("dice roll completed", "user_id", userID, "to_city", targetCity.Name,
		"dice_roll_id", diceRoll.ID, "visit_id", visit.ID)

	return &RollResult{
		VisitID:    visit.ID,
		DiceRollID: diceRoll.ID,
		Direction:  dir.Name,
		DistanceKm: dist,
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
