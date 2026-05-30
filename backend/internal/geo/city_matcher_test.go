package geo

import (
	"errors"
	"math"
	"testing"
)

func TestMatchNearestCity(t *testing.T) {
	cities := []CityPoint{
		{ID: 1, Lat: 0, Lng: 0},
		{ID: 2, Lat: 0, Lng: 1},
		{ID: 3, Lat: 0, Lng: 2},
	}

	tests := []struct {
		name                 string
		cities               []CityPoint
		targetLat, targetLng float64
		options              MatchOptions
		wantID               int64
		wantErr              error
	}{
		{
			name: "nearest city without exclusion", cities: cities,
			targetLat: 0, targetLng: 0.9, wantID: 2,
		},
		{
			name: "excludes current city", cities: cities,
			targetLat: 0, targetLng: 1, options: MatchOptions{ExcludeCityID: 2}, wantID: 1,
		},
		{
			name: "prefers farther unvisited city", cities: cities,
			targetLat: 0, targetLng: 1.1, options: MatchOptions{VisitedCityIDs: []int64{2}}, wantID: 3,
		},
		{
			name: "all visited allows nearest repeat", cities: cities,
			targetLat: 0, targetLng: 1.1, options: MatchOptions{ExcludeCityID: 1, VisitedCityIDs: []int64{1, 2, 3}}, wantID: 2,
		},
		{
			name: "stable tie break uses lower id", cities: []CityPoint{
				{ID: 9, Lat: 0, Lng: 1},
				{ID: 4, Lat: 0, Lng: -1},
			},
			targetLat: 0, targetLng: 0, wantID: 4,
		},
		{
			name: "empty city list", cities: nil,
			targetLat: 0, targetLng: 0, wantErr: ErrNoCandidates,
		},
		{
			name: "only current city", cities: []CityPoint{{ID: 1, Lat: 0, Lng: 0}},
			targetLat: 0, targetLng: 0, options: MatchOptions{ExcludeCityID: 1}, wantErr: ErrNoCandidates,
		},
		{
			name: "invalid target", cities: cities,
			targetLat: math.NaN(), targetLng: 0, wantErr: ErrInvalidCoordinates,
		},
		{
			name: "invalid city", cities: []CityPoint{{ID: 1, Lat: 91, Lng: 0}},
			targetLat: 0, targetLng: 0, wantErr: ErrInvalidCoordinates,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MatchNearestCity(tc.cities, tc.targetLat, tc.targetLng, tc.options)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("MatchNearestCity() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && got.ID != tc.wantID {
				t.Fatalf("MatchNearestCity() ID = %d, want %d", got.ID, tc.wantID)
			}
		})
	}
}

func TestMatchNearestCityTieBreakIsInputOrderIndependent(t *testing.T) {
	for _, cities := range [][]CityPoint{
		{{ID: 9, Lat: 0, Lng: 1}, {ID: 4, Lat: 0, Lng: -1}},
		{{ID: 4, Lat: 0, Lng: -1}, {ID: 9, Lat: 0, Lng: 1}},
	} {
		got, err := MatchNearestCity(cities, 0, 0, MatchOptions{})
		if err != nil {
			t.Fatalf("MatchNearestCity() error = %v", err)
		}
		if got.ID != 4 {
			t.Fatalf("MatchNearestCity() ID = %d, want stable lower ID 4", got.ID)
		}
	}
}
