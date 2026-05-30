package geo

import "math/rand"

// IntnSource is the minimal random source needed by the geo package.
// Tests can inject a deterministic source without changing production behavior.
type IntnSource interface {
	Intn(n int) int
}

type globalIntnSource struct{}

func (globalIntnSource) Intn(n int) int {
	return rand.Intn(n)
}

// Direction represents one of 8 compass directions.
type Direction struct {
	Name    string
	Bearing float64 // degrees, 0=North, clockwise
}

// Directions contains all 8 compass directions.
var Directions = [8]Direction{
	{Name: "北", Bearing: 0},
	{Name: "东北", Bearing: 45},
	{Name: "东", Bearing: 90},
	{Name: "东南", Bearing: 135},
	{Name: "南", Bearing: 180},
	{Name: "西南", Bearing: 225},
	{Name: "西", Bearing: 270},
	{Name: "西北", Bearing: 315},
}

// RandomDirection picks a random compass direction.
func RandomDirection() Direction {
	return RandomDirectionWithRand(globalIntnSource{})
}

// RandomDirectionWithRand picks a compass direction using the supplied source.
func RandomDirectionWithRand(rng IntnSource) Direction {
	return Directions[rng.Intn(len(Directions))]
}
