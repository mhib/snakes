package ai

import "github.com/mhib/snakes/board"

//NearestFoodAI finds shortest path to fruit
type NearestFoodAI struct {
	*BaseAI
}

//Run runs ai
func (ai *NearestFoodAI) Run() {
	for {
		select {
		case <-ai.QuitChannel:
			return
		case b := <-ai.NotifyChannel:
			snake, snakeErr := b.GetSnake(ai.SnakeID)
			if snakeErr != nil || snake.Lost {
				break
			}
			direction := findNearestFoodDirection(snake, b)
			ai.UpdateChannel <- board.Change{ai.SnakeID, direction}
		}
	}
}

type bfsEntry struct {
	board.Point
	direction int
	parent    *bfsEntry
}

func getInitialDirection(entry *bfsEntry) int {
	current := entry
	if current.parent == nil {
		return current.direction
	}
	for current.parent.parent != nil {
		current = current.parent
	}
	return current.direction
}

func findNearestFoodDirection(snake *board.Snake, b *board.Board) int {
	lastEntry := bfsEntry{snake.Head(), snake.PrevDirection, nil}
	queue := []bfsEntry{lastEntry}
	queued := make(map[board.Point]bool)
	queued[lastEntry.Point] = true
	for len(queue) > 0 {
		var current bfsEntry
		current, queue = queue[0], queue[1:]
		for direction, point := range b.Neighbours(current.Point) {
			if board.IsOpposite(direction, current.direction) || b.PartOfSnake(point) || queued[point] {
				continue
			}
			if current.parent == nil && mayCollideWithOtherSnake(point, snake, b) {
				continue
			}
			newEntry := bfsEntry{point, direction, &current}
			if b.IsFruit(point) {
				return getInitialDirection(&newEntry)
			}
			queue = append(queue, newEntry)
			queued[point] = true
			lastEntry = newEntry
		}
	}
	return getInitialDirection(&lastEntry) // If no path found, stay alive as long as possible
}

func mayCollideWithOtherSnake(p board.Point, currentSnake *board.Snake, b *board.Board) bool {
	for _, neighbour := range b.Neighbours(p) {
		for _, snake := range b.Snakes {
			if snake.ID == currentSnake.ID {
				continue
			}
			if snake.Head() == neighbour {
				return true
			}
		}
	}
	return false
}

//NewNearestFoodAI creates new NearestFoodAI
func NewNearestFoodAI(updateChannel chan board.Change, snakeID string) *NearestFoodAI {
	return &NearestFoodAI{NewAI(updateChannel, snakeID)}
}
