package main

import (
	"github.com/gorilla/websocket"
)

type lobbyConnection struct {
	conn  *websocket.Conn
	lobby *Lobby
	send  chan []byte
}

func NewLobbyConnection(conn *websocket.Conn, lobby *Lobby) *lobbyConnection {
	return &lobbyConnection{
		conn:  conn,
		lobby: lobby,
		send:  make(chan []byte, 256),
	}
}

func (lc *lobbyConnection) Run() {
	for {
		message := <-lc.send
		err := lc.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			lc.lobby.unregister <- lc
			return
		}
	}
}

type Lobby struct {
	broadcast      chan []byte
	register       chan *lobbyConnection
	unregister     chan *lobbyConnection
	connections    map[*lobbyConnection]bool
	initialMessage []byte
}

func NewLobby(initialMessage []byte) *Lobby {
	return &Lobby{
		broadcast:      make(chan []byte),
		register:       make(chan *lobbyConnection),
		unregister:     make(chan *lobbyConnection),
		connections:    make(map[*lobbyConnection]bool),
		initialMessage: initialMessage,
	}
}

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
