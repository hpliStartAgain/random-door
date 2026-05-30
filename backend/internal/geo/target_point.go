package geo

import (
	"math"
)

// DistanceLevels defines the 6 possible travel distances in km.
var DistanceLevels = []int{100, 200, 300, 500, 800, 1200}

// RandomDistance picks one of the 6 distance levels.
func RandomDistance() int {
	return RandomDistanceWithRand(globalIntnSource{})
}

// RandomDistanceWithRand picks a travel distance using the supplied source.
func RandomDistanceWithRand(rng IntnSource) int {
	return DistanceLevels[rng.Intn(len(DistanceLevels))]
}

// TargetPoint computes the destination point given a start point,
// bearing (degrees clockwise from north), and distance (km).
// Uses spherical geometry (destination point formula).
func TargetPoint(lat, lng, bearingDeg float64, distKm int) (lat2, lng2 float64) {
	delta := float64(distKm) / earthRadiusKm // angular distance in radians
	theta := toRad(bearingDeg)
	phi1 := toRad(lat)
	lambda1 := toRad(lng)

	phi2 := math.Asin(math.Sin(phi1)*math.Cos(delta) +
		math.Cos(phi1)*math.Sin(delta)*math.Cos(theta))
	lambda2 := lambda1 + math.Atan2(
		math.Sin(theta)*math.Sin(delta)*math.Cos(phi1),
		math.Cos(delta)-math.Sin(phi1)*math.Sin(phi2))

	lat2 = toDeg(phi2)
	lng2 = normalizeLongitude(toDeg(lambda2))
	return
}

func normalizeLongitude(lng float64) float64 {
	normalized := math.Mod(lng+180, 360)
	if normalized < 0 {
		normalized += 360
	}
	return normalized - 180
}
