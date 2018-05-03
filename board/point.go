package board

// Point - 2D point representation
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (r *Point) addPoint(other Point) {
	r.X += other.X
	r.Y += other.Y
}
