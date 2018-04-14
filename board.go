package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Change represent change of snake direction
type Change struct {
	ID        string
	Direction int
}

// Board represents 2d board
type Board struct {
	Width           int         `json:"width"`
	Length          int         `json:"length"`
	Snakes          []Snake     `json:"snakes"`
	Fruits          []Point     `json:"fruits"`
	State           int         `json:"state"`
	EndOnLastPlayer bool        `json:"-"`
	Tick            int         `json:"tick"`
	Changes         chan Change `json:"-"`
	End             chan bool   `json:"-"`
}

// Represents state of board
const (
	WAITING = iota
	PREPARING
	PLAYING
	ENDED
)

func (b *Board) isFull() bool {
	boardSize := b.Length * b.Width
	sum := len(b.Fruits)
	for _, s := range b.Snakes {
		sum += len(s.Body)
	}
	return sum == boardSize
}

func (b *Board) randomPoint() (Point, error) {
	if b.isFull() {
		return Point{}, errors.New("Board is full")
	}
	for {
		point := Point{rand.Intn(b.Width), rand.Intn(b.Length)}
		if !b.PartOfSnake(point) && !b.IsFruit(point) {
			return point, nil
		}
	}
}

func (b *Board) GetSnake(id string) (*Snake, error) {
	for _, snake := range b.Snakes {
		if snake.ID == id {
			return &snake, nil
		}
	}
	return &Snake{}, errors.New("Snake not found")
}

func (b *Board) AddSnake(id string, name string, color string, size int) {
	point, err := b.randomPoint()
	if err != nil {
		return
	}
	body := []Point{point}
	direction := rand.Intn(DOWN)
	b.Snakes = append(b.Snakes, Snake{
		body, 0, direction, direction, false, size - 1, name, color, id})
}

func (b *Board) tick() {
	needNewFruit := b.moveSnakes()
	if needNewFruit {
		b.generateFruit()
	}
	b.checkCollisions()
}

func (b *Board) generateFruit() {
	point, err := b.randomPoint()
	if err != nil {
		return
	}
	b.Fruits = append(b.Fruits, point)
}

func (b *Board) IsFruit(point Point) bool {
	for _, p := range b.Fruits {
		if p == point {
			return true
		}
	}
	return false
}

func (b *Board) removeFromFruits(point Point) {
	for i := range b.Fruits {
		if b.Fruits[i] == point {
			b.Fruits = append(b.Fruits[:i], b.Fruits[i+1:]...)
			return
		}
	}
}

func (b *Board) moveSnakes() bool {
	fruitEaten := false
	for i := range b.Snakes {
		if !b.Snakes[i].Lost {
			b.Snakes[i].move(b.Width, b.Length)
			if b.IsFruit(b.Snakes[i].Head()) {
				b.Snakes[i].Grow(2)
				b.removeFromFruits(b.Snakes[i].Head())
				fruitEaten = true
			}
		}
	}
	return fruitEaten
}

func (b *Board) checkCollisions() {
	for outer := range b.Snakes {
		for inner := range b.Snakes {
			if b.Snakes[inner].Lost {
				continue
			}
			if outer == inner && b.Snakes[inner].collideItself() {
				b.Snakes[inner].lose()
			} else if outer != inner && b.Snakes[inner].collide(&b.Snakes[outer]) {
				b.Snakes[inner].lose()
			}
		}
	}
}

func (b *Board) PartOfSnake(p Point) bool {
	for _, s := range b.Snakes {
		if s.includes(p) {
			return true
		}
	}
	return false
}

func (b *Board) print() {
	for x := 0; x < b.Width; x++ {
		fmt.Print("-")
	}
	fmt.Println("")
	for y := 0; y < b.Length; y++ {
		for x := 0; x < b.Width; x++ {
			currentPoint := Point{x, y}
			if b.IsFruit(currentPoint) {
				fmt.Print("o")
				continue
			}
			if b.PartOfSnake(currentPoint) {
				fmt.Print("x")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

func (b *Board) changeDirection(change *Change) {
	for snakeID := range b.Snakes {
		if b.Snakes[snakeID].ID == change.ID {
			b.Snakes[snakeID].changeDirection(change.Direction)
			return
		}
	}
}

func createTicker(tick time.Duration) *time.Ticker {
	if tick == 0*time.Second {
		ticker := time.NewTicker(1 * time.Second)
		ticker.Stop()
		return ticker
	}
	return time.NewTicker(tick)
}

func (b *Board) run(moveTick, foodTick time.Duration, callback func(*Board) bool) {
	b.State = PLAYING
	moveTicker := createTicker(moveTick).C
	foodTicker := createTicker(foodTick).C
	b.generateFruit()
	for b.going() {
		select {
		case <-b.End:
			break
		case <-moveTicker:
			b.tick()
			b.Tick++
			if !callback(b) {
				break
			}
		case <-foodTicker:
			b.generateFruit()
			b.Tick++
			if !callback(b) {
				break
			}
		case change := <-b.Changes:
			b.changeDirection(&change)
		}
	}
	b.State = ENDED
	callback(b)
}

func (b *Board) going() bool {
	if b.EndOnLastPlayer {
		return b.hasTwoPlayersPlaying()
	}
	return b.hasOnePlayerPlaying()
}

func (b *Board) hasTwoPlayersPlaying() bool {
	found := false
	for _, s := range b.Snakes {
		if !s.Lost {
			if found {
				return true
			}
			found = true
		}
	}
	return false
}

func (b *Board) hasOnePlayerPlaying() bool {
	for _, s := range b.Snakes {
		if !s.Lost {
			return true
		}
	}
	return false
}

func (b *Board) Neighbours(p Point) []Point {
	left := p
	left.X = Modulo(p.X-1, b.Width)
	right := p
	right.X = Modulo(p.X+1, b.Width)
	up := p
	up.Y = Modulo(p.Y-1, b.Length)
	down := p
	down.Y = Modulo(p.Y+1, b.Length)
	return []Point{left, right, up, down}
}
