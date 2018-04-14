package main

type AI interface {
	Notify(*Board)
	Quit()
	Run()
}

// BaseAI defines basic AI fields
type BaseAI struct {
	UpdateChannel chan Change
	NotifyChannel chan *Board
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

// Notify notifies ai about Board update
func (ai *BaseAI) Notify(b *Board) {
	ai.NotifyChannel <- b
}

// NewAI returns new BaseAI
func NewAI(updateChannel chan Change, snakeID string) *BaseAI {
	return &BaseAI{
		UpdateChannel: updateChannel,
		NotifyChannel: make(chan *Board, 100),
		SnakeID:       snakeID,
		QuitChannel:   make(chan bool, 1),
	}
}
