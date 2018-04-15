package main

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
		case board := <-ai.NotifyChannel:
			snake, snakeErr := board.GetSnake(ai.SnakeID)
			if snakeErr != nil {
				break
			}
			direction := findNearestFoodDirection(snake, board)
			ai.UpdateChannel <- Change{ai.SnakeID, direction}
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
	if current.parent == nil {
		return current.direction
	}
	for current.parent.parent != nil {
		current = current.parent
	}
	return current.direction
}

func findNearestFoodDirection(snake *Snake, board *Board) int {
	lastEntry := bfsEntry{snake.Head(), snake.PrevDirection, nil}
	queue := []bfsEntry{lastEntry}
	queued := make(map[Point]bool)
	queued[lastEntry.Point] = true
	for len(queue) > 0 {
		var current bfsEntry
		current, queue = queue[0], queue[1:]
		for direction, point := range board.Neighbours(current.Point) {
			if IsOpposite(direction, current.direction) || board.PartOfSnake(point) || queued[point] {
				continue
			}
			newEntry := bfsEntry{point, direction, &current}
			if board.IsFruit(point) {
				return getInitialDirection(&newEntry)
			}
			queue = append(queue, newEntry)
			queued[point] = true
			lastEntry = newEntry
		}
	}
	return getInitialDirection(&lastEntry) // If no path found, stay alive as long as possible
}

//NewNearestFoodAI creates new NearestFoodAI
func NewNearestFoodAI(updateChannel chan Change, snakeID string) *NearestFoodAI {
	return &NearestFoodAI{NewAI(updateChannel, snakeID)}
}
