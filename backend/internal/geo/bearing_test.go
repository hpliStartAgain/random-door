package geo

import "testing"

type scriptedIntnSource struct {
	values []int
	index  int
}

func (s *scriptedIntnSource) Intn(n int) int {
	value := s.values[s.index]
	s.index++
	if value < 0 || value >= n {
		panic("scripted random value out of range")
	}
	return value
}

func TestDirections(t *testing.T) {
	want := []Direction{
		{Name: "北", Bearing: 0},
		{Name: "东北", Bearing: 45},
		{Name: "东", Bearing: 90},
		{Name: "东南", Bearing: 135},
		{Name: "南", Bearing: 180},
		{Name: "西南", Bearing: 225},
		{Name: "西", Bearing: 270},
		{Name: "西北", Bearing: 315},
	}
	if len(Directions) != len(want) {
		t.Fatalf("Directions count = %d, want %d", len(Directions), len(want))
	}
	for i := range want {
		if Directions[i] != want[i] {
			t.Fatalf("Directions[%d] = %#v, want %#v", i, Directions[i], want[i])
		}
	}
}

func TestRandomDirectionWithRand(t *testing.T) {
	values := make([]int, len(Directions))
	for i := range values {
		values[i] = i
	}
	rng := &scriptedIntnSource{values: values}

	for i, want := range Directions {
		if got := RandomDirectionWithRand(rng); got != want {
			t.Fatalf("RandomDirectionWithRand() call %d = %#v, want %#v", i, got, want)
		}
	}
}

func TestRandomDirectionReturnsDefinedValue(t *testing.T) {
	defined := make(map[Direction]bool, len(Directions))
	for _, direction := range Directions {
		defined[direction] = true
	}
	for i := 0; i < 1000; i++ {
		if got := RandomDirection(); !defined[got] {
			t.Fatalf("RandomDirection() = %#v, want one of %#v", got, Directions)
		}
	}
}
