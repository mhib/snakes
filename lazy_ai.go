package main

import "math/rand"

// LazyAI that moves when it has to, or has a fruit as neighbour
type LazyAI struct {
	*BaseAI
}

//Run runs AI
func (ai *LazyAI) Run() {
	for {
		select {
		case <-ai.QuitChannel:
			return
		case board := <-ai.NotifyChannel:
			snake, snakeErr := board.GetSnake(ai.SnakeID)
			if snakeErr != nil {
				break
			}
			direction := getNewDirection(snake, board)
			if direction >= 0 {
				ai.UpdateChannel <- Change{ai.SnakeID, direction}
			}
		}
	}
}

func getNewDirection(snake *Snake, board *Board) int {
	results := make([]int, 0)
	foundCurrent := false
	for direction, point := range board.Neighbours(snake.Head()) {
		if IsOpposite(direction, snake.PrevDirection) || board.PartOfSnake(point) {
			continue
		}
		if board.IsFruit(point) {
			return direction
		}
		if direction == snake.PrevDirection {
			foundCurrent = true
		}
		results = append(results, direction)
	}
	if len(results) == 0 || foundCurrent {
		return -1
	}
	return results[rand.Intn(len(results))]
}

//NewLazyAI creates new LazyAI
func NewLazyAI(updateChannel chan Change, snakeID string) *LazyAI {
	return &LazyAI{NewAI(updateChannel, snakeID)}
}
