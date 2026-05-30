package geo

import (
	"math"
	"testing"
)

func TestDistanceLevels(t *testing.T) {
	want := []int{100, 200, 300, 500, 800, 1200}
	if len(DistanceLevels) != len(want) {
		t.Fatalf("DistanceLevels count = %d, want %d", len(DistanceLevels), len(want))
	}
	for i := range want {
		if DistanceLevels[i] != want[i] {
			t.Fatalf("DistanceLevels[%d] = %d, want %d", i, DistanceLevels[i], want[i])
		}
	}
}

func TestRandomDistanceWithRand(t *testing.T) {
	values := make([]int, len(DistanceLevels))
	for i := range values {
		values[i] = i
	}
	rng := &scriptedIntnSource{values: values}

	for i, want := range DistanceLevels {
		if got := RandomDistanceWithRand(rng); got != want {
			t.Fatalf("RandomDistanceWithRand() call %d = %d, want %d", i, got, want)
		}
	}
}

func TestRandomDistanceReturnsDefinedLevel(t *testing.T) {
	defined := make(map[int]bool, len(DistanceLevels))
	for _, distance := range DistanceLevels {
		defined[distance] = true
	}
	for i := 0; i < 1000; i++ {
		if got := RandomDistance(); !defined[got] {
			t.Fatalf("RandomDistance() = %d, want one of %#v", got, DistanceLevels)
		}
	}
}

func TestTargetPoint(t *testing.T) {
	tests := []struct {
		name                               string
		lat, lng, bearing                  float64
		distanceKm                         int
		wantLat, wantLng, toleranceDegrees float64
	}{
		{name: "north from equator", lat: 0, lng: 0, bearing: 0, distanceKm: 100, wantLat: 0.899322, wantLng: 0, toleranceDegrees: 0.000001},
		{name: "east from equator", lat: 0, lng: 0, bearing: 90, distanceKm: 100, wantLat: 0, wantLng: 0.899322, toleranceDegrees: 0.000001},
		{name: "cross dateline east", lat: 0, lng: 179, bearing: 90, distanceKm: 300, wantLat: 0, wantLng: -178.302035, toleranceDegrees: 0.000001},
		{name: "normalize large negative longitude", lat: 0, lng: -720, bearing: 0, distanceKm: 0, wantLat: 0, wantLng: 0, toleranceDegrees: 0.000001},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLat, gotLng := TargetPoint(tc.lat, tc.lng, tc.bearing, tc.distanceKm)
			if math.Abs(gotLat-tc.wantLat) > tc.toleranceDegrees || math.Abs(gotLng-tc.wantLng) > tc.toleranceDegrees {
				t.Fatalf("TargetPoint() = (%.9f, %.9f), want (%.9f, %.9f)", gotLat, gotLng, tc.wantLat, tc.wantLng)
			}
			if gotLng < -180 || gotLng >= 180 {
				t.Fatalf("TargetPoint() longitude = %.9f, want [-180, 180)", gotLng)
			}
		})
	}
}

func TestTargetPointRoundTripDistance(t *testing.T) {
	tests := []struct {
		lat, lng, bearing float64
		distanceKm        int
	}{
		{lat: 39.9042, lng: 116.4074, bearing: 225, distanceKm: 800},
		{lat: 0, lng: 179, bearing: 90, distanceKm: 1200},
		{lat: -33.8688, lng: 151.2093, bearing: 315, distanceKm: 500},
	}

	for _, tc := range tests {
		targetLat, targetLng := TargetPoint(tc.lat, tc.lng, tc.bearing, tc.distanceKm)
		gotDistance := Haversine(tc.lat, tc.lng, targetLat, targetLng)
		if math.Abs(gotDistance-float64(tc.distanceKm)) > 0.000001 {
			t.Fatalf("round trip distance = %.9fkm, want %dkm", gotDistance, tc.distanceKm)
		}
	}
}
