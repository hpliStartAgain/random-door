package geo

import (
	"errors"
	"fmt"
)

// CityPoint represents a city with coordinates for matching.
type CityPoint struct {
	ID  int64
	Lat float64
	Lng float64
}

// MatchOptions controls city matching behavior.
type MatchOptions struct {
	ExcludeCityID  int64
	VisitedCityIDs []int64
}

var (
	ErrNoCandidates       = errors.New("no candidate cities available")
	ErrInvalidCoordinates = errors.New("invalid coordinates")
)

// MatchNearestCity finds the best city to visit given a target point.
// Priority: unvisited city nearest to target > fallback strategies.
func MatchNearestCity(cities []CityPoint, targetLat, targetLng float64, opt MatchOptions) (CityPoint, error) {
	if len(cities) == 0 {
		return CityPoint{}, ErrNoCandidates
	}
	if !validCoordinates(targetLat, targetLng) {
		return CityPoint{}, fmt.Errorf("%w: target point", ErrInvalidCoordinates)
	}

	visited := make(map[int64]bool)
	for _, id := range opt.VisitedCityIDs {
		visited[id] = true
	}

	// Separate candidates (exclude current city)
	var unvisited, allCandidates []CityPoint
	for _, c := range cities {
		if !validCoordinates(c.Lat, c.Lng) {
			return CityPoint{}, fmt.Errorf("%w: city %d", ErrInvalidCoordinates, c.ID)
		}
		if c.ID == opt.ExcludeCityID {
			continue
		}
		allCandidates = append(allCandidates, c)
		if !visited[c.ID] {
			unvisited = append(unvisited, c)
		}
	}

	if len(allCandidates) == 0 {
		return CityPoint{}, ErrNoCandidates
	}

	// Normal path: nearest unvisited city to target
	if len(unvisited) > 0 {
		return findNearest(unvisited, targetLat, targetLng), nil
	}

	// Fallback: all cities were visited, so allow a repeat while still honoring
	// the random target point generated from this roll.
	return findNearest(allCandidates, targetLat, targetLng), nil
}

func findNearest(cities []CityPoint, lat, lng float64) CityPoint {
	var nearest CityPoint
	var nearestDistance float64
	found := false
	for _, c := range cities {
		d := Haversine(lat, lng, c.Lat, c.Lng)
		if !found || d < nearestDistance || (d == nearestDistance && c.ID < nearest.ID) {
			nearest = c
			nearestDistance = d
			found = true
		}
	}
	return nearest
}
