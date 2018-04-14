package main

import (
	"errors"
)

type NearestFoodAI struct {
	*BaseAI
}

func (ai *NearestFoodAI) Run() {
	for {
		select {
		case <-ai.QuitChannel:
			break
		case board := <-ai.NotifyChannel:
			snake, snakeErr := board.GetSnake(ai.SnakeID)
			if snakeErr != nil {
				break
			}
			direction, dirErr := findNearestFoodDirection(snake, board)
			if dirErr == nil {
				ai.UpdateChannel <- Change{ai.SnakeID, direction}
			}
		}
	}
}

type bfsEntry struct {
	Point
	direction int
	parent    *bfsEntry
}

func getInitialDirection(entry *bfsEntry) int {
	current := entry
	for current.parent.parent != nil {
		current = current.parent
	}
	return current.direction
}

func findNearestFoodDirection(snake *Snake, board *Board) (int, error) {
	initial := bfsEntry{snake.Head(), snake.PrevDirection, nil}
	queue := []bfsEntry{initial}
	queued := make(map[Point]bool)
	queued[initial.Point] = true
	for len(queue) > 0 {
		var current bfsEntry
		current, queue = queue[0], queue[1:]
		for direction, point := range neighbours(current.Point, board) {
			if IsOpposite(direction, current.direction) || board.PartOfSnake(point) || queued[point] {
				continue
			}
			newEntry := bfsEntry{point, direction, &current}
			if board.IsFruit(point) {
				return getInitialDirection(&newEntry), nil
			}
			queue = append(queue, newEntry)
			queued[point] = true
		}
	}
	return -1, errors.New("Path not found")
}

// NewAI returns new BaseAI
func NewNearestFoodAI(updateChannel chan Change, snakeID string) *NearestFoodAI {
	return &NearestFoodAI{NewAI(updateChannel, snakeID)}
}

func neighbours(p Point, b *Board) []Point {
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
