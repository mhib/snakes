package ai

import (
	"math/rand"

	"github.com/mhib/snakes/board"
)

// RandomMoveAI goes to fruit if fruit is neighbour; goes to random free point otherwise
type RandomMoveAI struct {
	*BaseAI
}

// Run runs AI
func (ai *RandomMoveAI) Run() {
	for {
		select {
		case <-ai.QuitChannel:
			return
		case b := <-ai.NotifyChannel:
			snake, snakeErr := b.GetSnake(ai.SnakeID)
			if snakeErr != nil || snake.Lost {
				break
			}
			direction := getRandomDirection(snake, b)
			if direction >= 0 {
				ai.UpdateChannel <- board.Change{ai.SnakeID, direction}
			}
		}
	}
}

func getRandomDirection(snake *board.Snake, b *board.Board) int {
	results := make([]int, 0)
	for direction, point := range b.Neighbours(snake.Head()) {
		if board.IsOpposite(direction, snake.PrevDirection) || b.PartOfSnake(point) {
			continue
		}
		if b.IsFruit(point) {
			return direction
		}
		results = append(results, direction)
	}
	if len(results) == 0 {
		return -1
	}
	return results[rand.Intn(len(results))]
}

//NewRandomMoveAI creates new RandomMoveAI
func NewRandomMoveAI(updateChannel chan board.Change, snakeID string) *RandomMoveAI {
	return &RandomMoveAI{NewAI(updateChannel, snakeID)}
}
