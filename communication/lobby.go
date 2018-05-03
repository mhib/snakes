package communication

import (
	"time"

	"github.com/gorilla/websocket"
)

// LobbyConnection - websocket lobby connection
type LobbyConnection struct {
	conn  *websocket.Conn
	lobby *Lobby
	send  chan []byte
}

// NewLobbyConnection initializes lobby connection and returns pointer to it
func NewLobbyConnection(conn *websocket.Conn, lobby *Lobby) *LobbyConnection {
	return &LobbyConnection{
		conn:  conn,
		lobby: lobby,
		send:  make(chan []byte, 256),
	}
}

// Run send pings and write messages
func (lc *LobbyConnection) Run() {
	pingTicker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-pingTicker.C:
			lc.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 250))
			if err := lc.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				lc.lobby.Unregister <- lc
				return
			}
			lc.conn.SetWriteDeadline(time.Time{})
		case message := <-lc.send:
			err := lc.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				lc.lobby.Unregister <- lc
				return
			}
		}
	}
}

// Lobby represents lobby state
type Lobby struct {
	Broadcast      chan []byte
	Register       chan *LobbyConnection
	Unregister     chan *LobbyConnection
	Connections    map[*LobbyConnection]bool
	InitialMessage []byte
}

// NewLobby initializes new lobby and returns pointer to it
func NewLobby(initialMessage []byte) *Lobby {
	return &Lobby{
		Broadcast:      make(chan []byte),
		Register:       make(chan *LobbyConnection),
		Unregister:     make(chan *LobbyConnection),
		Connections:    make(map[*LobbyConnection]bool),
		InitialMessage: initialMessage,
	}
}

// Run Registers, unregisters Connections and Broadcast messages to them
func (l *Lobby) Run() {
	lastMessage := l.InitialMessage
	for {
		select {
		case connection := <-l.Register:
			l.Connections[connection] = true
			connection.send <- lastMessage
			go connection.Run()
		case connection := <-l.Unregister:
			if _, ok := l.Connections[connection]; ok {
				delete(l.Connections, connection)
				close(connection.send)
			}
		case lastMessage = <-l.Broadcast:
			for connection := range l.Connections {
				select {
				case connection.send <- lastMessage:
				default:
					close(connection.send)
					delete(l.Connections, connection)
				}
			}
		}
	}
}
