package utils

// Modulo is mathematical modulo
func Modulo(n, m int) int {
	val := n % m
	if val >= 0 {
		return val
	}
	return val + m
}
