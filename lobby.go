package main

import (
	"github.com/gorilla/websocket"
	"time"
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
				lc.lobby.unregister <- lc
				return
			}
			lc.conn.SetWriteDeadline(time.Time{})
		case message := <-lc.send:
			err := lc.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				lc.lobby.unregister <- lc
				return
			}
		}
	}
}

// Lobby represents lobby state
type Lobby struct {
	broadcast      chan []byte
	register       chan *LobbyConnection
	unregister     chan *LobbyConnection
	connections    map[*LobbyConnection]bool
	initialMessage []byte
}

// NewLobby initializes new lobby and returns pointer to it
func NewLobby(initialMessage []byte) *Lobby {
	return &Lobby{
		broadcast:      make(chan []byte),
		register:       make(chan *LobbyConnection),
		unregister:     make(chan *LobbyConnection),
		connections:    make(map[*LobbyConnection]bool),
		initialMessage: initialMessage,
	}
}

// Run Registers, unregisters connections and broadcast messages to them
func (l *Lobby) Run() {
	lastMessage := l.initialMessage
	for {
		select {
		case connection := <-l.register:
			l.connections[connection] = true
			connection.send <- lastMessage
			go connection.Run()
		case connection := <-l.unregister:
			if _, ok := l.connections[connection]; ok {
				delete(l.connections, connection)
				close(connection.send)
			}
		case lastMessage = <-l.broadcast:
			for connection := range l.connections {
				select {
				case connection.send <- lastMessage:
				default:
					close(connection.send)
					delete(l.connections, connection)
				}
			}
		}
	}
}
