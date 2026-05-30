package geo

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lng1, lat2, lng2 float64
		wantKm, toleranceKm    float64
	}{
		{name: "same point", lat1: 39.9042, lng1: 116.4074, lat2: 39.9042, lng2: 116.4074, wantKm: 0, toleranceKm: 0.000001},
		{name: "beijing to xian", lat1: 39.9042, lng1: 116.4074, lat2: 34.3416, lng2: 108.9398, wantKm: 910, toleranceKm: 20},
		{name: "one degree equator", lat1: 0, lng1: 0, lat2: 0, lng2: 1, wantKm: 111.195, toleranceKm: 0.001},
		{name: "antipodes", lat1: 0, lng1: 0, lat2: 0, lng2: 180, wantKm: math.Pi * earthRadiusKm, toleranceKm: 0.001},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Haversine(tc.lat1, tc.lng1, tc.lat2, tc.lng2)
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("Haversine() = %v, want finite distance", got)
			}
			if delta := math.Abs(got - tc.wantKm); delta > tc.toleranceKm {
				t.Fatalf("Haversine() = %.6fkm, want %.6fkm +/- %.6fkm", got, tc.wantKm, tc.toleranceKm)
			}
		})
	}
}

func TestHaversineIsSymmetric(t *testing.T) {
	forward := Haversine(39.9042, 116.4074, 34.3416, 108.9398)
	backward := Haversine(34.3416, 108.9398, 39.9042, 116.4074)
	if math.Abs(forward-backward) > 0.000001 {
		t.Fatalf("Haversine symmetry difference = %.9fkm", math.Abs(forward-backward))
	}
}
