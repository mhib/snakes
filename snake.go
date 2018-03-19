package main

type Snake struct {
	Body      []Point `json:"body"`
	Points    int     `json:"points"`
	Direction int     `json:"-"`
	Lost      bool    `json:"-"`
	Eaten     int     `json:"-"`
	ID        string  `json:"id"`
}

const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

func isOpposite(x, y int) bool {
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
	if !isOpposite(s.Direction, Direction) {
		s.Direction = Direction
	}
}

func (s *Snake) head() Point {
	return s.Body[0]
}

func (s *Snake) grow(size uint) {
	s.Points++
	s.Eaten += int(size)
}

func (s *Snake) collide(other *Snake) bool {
	for _, b := range other.Body {
		if b == s.head() {
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
		if b == s.head() {
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

func modulo(n, m int) int {
	val := n % m
	if val >= 0 {
		return val
	}
	return val + m
}

func (s *Snake) move(width, length int) {
	currentHead := s.head()
	switch s.Direction {
	case LEFT:
		currentHead.X = modulo(currentHead.X-1, width)
	case RIGHT:
		currentHead.X = modulo(currentHead.X+1, width)
	case UP:
		currentHead.Y = modulo(currentHead.Y-1, length)
	case DOWN:
		currentHead.Y = modulo(currentHead.Y+1, length)
	}
	s.Body = append([]Point{currentHead}, s.Body...)
	if s.Eaten == 0 {
		s.Body = s.Body[:len(s.Body)-1]
	} else {
		s.Eaten--
	}
}
