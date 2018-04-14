package main

import (
	"math/rand"
)

// RandomMoveAI goes to fruit if fruit is neighbour; goes to random free point otherwise
type RandomMoveAI struct {
	*BaseAI
}

func (ai *RandomMoveAI) Run() {
	for {
		select {
		case <-ai.QuitChannel:
			return
		case board := <-ai.NotifyChannel:
			snake, snakeErr := board.GetSnake(ai.SnakeID)
			if snakeErr != nil {
				break
			}
			direction := getRandomDirection(snake, board)
			if direction >= 0 {
				ai.UpdateChannel <- Change{ai.SnakeID, direction}
			}
		}
	}
}

func getRandomDirection(snake *Snake, board *Board) int {
	results := make([]int, 0)
	for direction, point := range board.Neighbours(snake.Head()) {
		if IsOpposite(direction, snake.PrevDirection) || board.PartOfSnake(point) {
			continue
		}
		if board.IsFruit(point) {
			return direction
		}
		results = append(results, direction)
	}
	if len(results) == 0 {
		return -1
	}
	return results[rand.Intn(len(results))]
}

func NewRandomMoveAI(updateChannel chan Change, snakeID string) *RandomMoveAI {
	return &RandomMoveAI{NewAI(updateChannel, snakeID)}
}
