package communication

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"github.com/mhib/snakes/ai"
	"github.com/mhib/snakes/board"
)

// Client represents user connection
type Client struct {
	conn           *websocket.Conn
	name           string
	color          string
	id             string
	game           *Game
	send           chan []byte
	endPingChannel chan bool
}

// NewClient Initializes new client and returns pointer to it
func NewClient(conn *websocket.Conn, name, color, id string, game *Game) *Client {
	return &Client{
		conn,
		name,
		color,
		id,
		game,
		make(chan []byte, 10),
		make(chan bool, 1),
	}
}

func (client *Client) pingGameConnection() {
	pingTicker := time.NewTicker(time.Millisecond * 500)
	defer pingTicker.Stop()
	for {
		select {
		case <-client.endPingChannel:
			return
		case <-pingTicker.C:
			client.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 250))
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				client.game.Unregister <- client
				return
			}
			client.conn.SetWriteDeadline(time.Time{})
		}
	}
}

func (client *Client) read() {
	for {
		var msg userMoveMessage
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			client.game.Unregister <- client
			break
		}
		client.game.sendChange(msg, client.id)
	}

}

func (client *Client) write() {
	for {
		message := <-client.send
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.game.Unregister <- client
			return
		}
	}
}

func (client *Client) run() {
	go client.read()
	go client.write()
}

// Game represent Game state
type Game struct {
	ID                 string
	Board              board.Board
	UsersCount         int
	Users              map[string]*Client
	FoodTick           time.Duration
	MoveTick           time.Duration
	Broadcast          chan []byte
	Register           chan *Client
	Unregister         chan *Client
	ChangeStateChannel chan bool
	DisposeChannel     chan *Game
	Bots               []ai.AI
}

type userMoveMessage struct {
	Direction string `json:"direction"`
}

func (g *Game) sendChange(msg userMoveMessage, userID string) {
	var change board.Change
	switch msg.Direction {
	case "LEFT":
		change = board.Change{userID, board.LEFT}
	case "RIGHT":
		change = board.Change{userID, board.RIGHT}
	case "UP":
		change = board.Change{userID, board.UP}
	case "DOWN":
		change = board.Change{userID, board.DOWN}
	}
	g.Board.Changes <- change
}

// Run runs game and broadcasts states update
func (g *Game) Run() {
	for {
		select {
		case client := <-g.Register:
			g.addUser(client)
		case client := <-g.Unregister:
			if _, ok := g.Users[client.id]; ok {
				g.ChangeStateChannel <- true
				delete(g.Users, client.id)
				close(client.send)
			}
		case message := <-g.Broadcast:
			for _, user := range g.Users {
				select {
				case user.send <- message:
				default:
					close(user.send)
					delete(g.Users, user.id)
				}
			}
		}
	}
}

func (g *Game) addUser(user *Client) {
	g.Users[user.id] = user
	g.Board.AddSnake(user.id, user.name, user.color, 3)
	go user.pingGameConnection()
	g.ChangeStateChannel <- true
	if g.UsersCount == len(g.Users) {
		go g.start()
	}
}

func (g *Game) runBots() {
	for _, bot := range g.Bots {
		go func(b ai.AI) {
			b.Run()
		}(bot)
	}
}

func (g *Game) start() {
	g.Board.State = board.PREPARING
	g.ChangeStateChannel <- true
	g.prepareUsers()
	g.broadcastBoard()
	g.runBots()
	time.Sleep(5 * time.Second)
	g.runBoard()
}

func (g *Game) prepareUsers() {
	for _, user := range g.Users {
		user.endPingChannel <- true
		go func(c *Client) {
			c.run()
		}(user)
	}
}

func (g *Game) broadcastBoard() {
	val, err := json.Marshal(g.Board)
	if err != nil {
		return
	}
	g.Broadcast <- val
	for _, bot := range g.Bots {
		go func(b ai.AI) {
			b.Notify(&g.Board)
		}(bot)
	}
}

func (g *Game) quitBots() {
	for _, bot := range g.Bots {
		bot.Quit()
	}
}

func (g *Game) handleBoardChange(b *board.Board) bool {
	g.broadcastBoard()
	if len(g.Users) == 0 {
		go g.quitBots()
		g.DisposeChannel <- g
		return false
	}
	return true
}

func (g *Game) runBoard() {
	g.Board.Run(g.MoveTick, g.FoodTick, g.handleBoardChange)
}
