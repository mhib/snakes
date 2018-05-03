package ai

import "math/rand"

import "github.com/mhib/snakes/board"

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
		case b := <-ai.NotifyChannel:
			snake, snakeErr := b.GetSnake(ai.SnakeID)
			if snakeErr != nil {
				break
			}
			direction := getNewDirection(snake, b)
			if direction >= 0 {
				ai.UpdateChannel <- board.Change{ai.SnakeID, direction}
			}
		}
	}
}

func getNewDirection(snake *board.Snake, b *board.Board) int {
	results := make([]int, 0)
	foundCurrent := false
	for direction, point := range b.Neighbours(snake.Head()) {
		if board.IsOpposite(direction, snake.PrevDirection) || b.PartOfSnake(point) {
			continue
		}
		if b.IsFruit(point) {
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
func NewLazyAI(updateChannel chan board.Change, snakeID string) *LazyAI {
	return &LazyAI{NewAI(updateChannel, snakeID)}
}
