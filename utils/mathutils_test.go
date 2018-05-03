package utils

import "testing"

func TestModuloPositive(t *testing.T) {
	res := Modulo(4, 10)
	if res != 4 {
		t.Error("Modulo fail")
	}
	res = Modulo(15, 10)
	if res != 5 {
		t.Error("Modulo fail")
	}
}

func TestModuloNegative(t *testing.T) {
	res := Modulo(-1, 10)
	if res != 9 {
		t.Error("Modulo fail")
	}
	res = Modulo(-39, 10)
	if res != 1 {
		t.Error("Modulo fail")
	}
}
