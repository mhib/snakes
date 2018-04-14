package main

type NearestFoodAI struct {
	*BaseAI
}

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
		for direction, point := range neighbours(current.Point, board) {
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
	return getInitialDirection(&lastEntry) // If no path found stay alive as long as possible
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
