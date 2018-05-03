package board

import "testing"

func TestAddPositive(t *testing.T) {
	first := Point{1, 2}
	second := Point{3, 4}
	first.addPoint(second)
	if first.X != 4 || first.Y != 6 {
		t.Error("Wrong value after addition")
	}
}

func TestAddNegative(t *testing.T) {
	first := Point{1, 2}
	second := Point{-4, -3}
	first.addPoint(second)
	if first.X != -3 || first.Y != -1 {
		t.Error("Wrong value after addition")
	}
}
