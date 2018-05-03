package ai

import "github.com/mhib/snakes/board"

// AI defines basic AI methods
type AI interface {
	Notify(*board.Board)
	Quit()
	Run()
}

// BaseAI defines basic AI fields
type BaseAI struct {
	UpdateChannel chan board.Change
	NotifyChannel chan *board.Board
	SnakeID       string
	QuitChannel   chan bool
}

// Run Runs BaseAI (should be overwritten)
func (ai *BaseAI) Run() {
	<-ai.QuitChannel
}

// Quit stops ai from working
func (ai *BaseAI) Quit() {
	ai.QuitChannel <- true
}

// Notify notifies ai about board.Board update
func (ai *BaseAI) Notify(b *board.Board) {
	ai.NotifyChannel <- b
}

// NewAI returns new BaseAI
func NewAI(updateChannel chan board.Change, snakeID string) *BaseAI {
	return &BaseAI{
		UpdateChannel: updateChannel,
		NotifyChannel: make(chan *board.Board, 100),
		SnakeID:       snakeID,
		QuitChannel:   make(chan bool, 1),
	}
}
