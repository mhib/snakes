package main

//Snake - snake's state
type Snake struct {
	Body          []Point `json:"body"`
	Points        int     `json:"points"`
	Direction     int     `json:"-"`
	PrevDirection int     `json:"-"`
	Lost          bool    `json:"-"`
	Eaten         int     `json:"-"`
	Name          string  `json:"name"`
	Color         string  `json:"color"`
	ID            string  `json:"id"`
}

//Possible directions
const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

//IsOpposite returns true if direction are opposite
func IsOpposite(x, y int) bool {
	switch x {
	case LEFT:
		return y == RIGHT
	case RIGHT:
		return y == LEFT
	case UP:
		return y == DOWN
	case DOWN:
		return y == UP
	default:
		return true
	}
}

func (s *Snake) changeDirection(Direction int) {
	if !IsOpposite(s.PrevDirection, Direction) {
		s.Direction = Direction
	}
}

//Head returns head of a snake
func (s *Snake) Head() Point {
	return s.Body[0]
}

//Grow grows a snake
func (s *Snake) Grow(size uint) {
	s.Points++
	s.Eaten += int(size)
}

func (s *Snake) collide(other *Snake) bool {
	for _, b := range other.Body {
		if b == s.Head() {
			return true
		}
	}
	return false
}

func (s *Snake) collideItself() bool {
	for i, b := range s.Body {
		if i == 0 {
			continue
		}
		if b == s.Head() {
			return true
		}
	}
	return false
}

func (s *Snake) lose() {
	s.Lost = true
	s.Body = s.Body[1:]
}

func (s *Snake) includes(point Point) bool {
	for _, b := range s.Body {
		if b == point {
			return true
		}
	}
	return false
}

func (s *Snake) move(width, length int) {
	s.PrevDirection = s.Direction
	currentHead := s.Head()
	switch s.Direction {
	case LEFT:
		currentHead.X = Modulo(currentHead.X-1, width)
	case RIGHT:
		currentHead.X = Modulo(currentHead.X+1, width)
	case UP:
		currentHead.Y = Modulo(currentHead.Y-1, length)
	case DOWN:
		currentHead.Y = Modulo(currentHead.Y+1, length)
	}
	s.Body = append([]Point{currentHead}, s.Body...)
	if s.Eaten == 0 {
		s.Body = s.Body[:len(s.Body)-1]
	} else {
		s.Eaten--
	}
}
