package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Game struct {
	ID              string
	Board           Board
	UsersCount      int
	Users           []*websocket.Conn
	FoodTick        time.Duration
	MoveTick        time.Duration
	EndPingChannels map[string]chan bool
}

type UserMoveMessage struct {
	Direction string `json:"direction"`
}

type UserConnectionMessage struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (g *Game) sendChange(msg UserMoveMessage, userID string) {
	var change Change
	switch msg.Direction {
	case "LEFT":
		change = Change{userID, LEFT}
	case "RIGHT":
		change = Change{userID, RIGHT}
	case "UP":
		change = Change{userID, UP}
	case "DOWN":
		change = Change{userID, DOWN}
	}
	g.Board.Changes <- change
}

func (g *Game) addUser(connection *websocket.Conn, userID string, name string, color string) {
	g.Board.Lock()
	g.Users = append(g.Users, connection)
	g.EndPingChannels[userID] = make(chan bool, 1)
	g.Board.Unlock()
	g.Board.addSnake(userID, name, color, 3)
	lobbyUpdateChannel <- true
	if g.UsersCount == len(g.Users) {
		g.Board.State = PREPARING
		g.StopPing()
		g.NotifyUsers()
		go func() {
			time.Sleep(5 * time.Second)
			handleChanges(g)
		}()
	}
}

func (g *Game) RemoveConnection(connection *websocket.Conn, userID string) {
	g.Board.Lock()
	defer g.Board.Unlock()
	delete(g.EndPingChannels, userID)
	newUsers := g.Users[:0]
	for _, user := range g.Users {
		if user != connection {
			newUsers = append(newUsers, user)
		}
	}
	g.Users = newUsers
	lobbyUpdateChannel <- true
}

func (g *Game) StopPing() {
	for _, v := range g.EndPingChannels {
		v <- true
	}
}

type gamesType struct {
	sync.RWMutex
	m map[string]*Game
}

func milliseconds(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1000000
}

func (gameMap gamesType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")

	for _, value := range gameMap.m {
		if value.Board.State != WAITING || value.UsersCount == 1 || len(value.Users) == 0 {
			continue
		}
		buffer.WriteString("{")
		buffer.WriteString(fmt.Sprintf("\"id\":\"%s\",", value.ID))
		buffer.WriteString(fmt.Sprintf("\"players\":%d,", value.UsersCount))
		buffer.WriteString(fmt.Sprintf("\"connected\":%d,", len(value.Users)))
		buffer.WriteString(fmt.Sprintf("\"width\":%d,", value.Board.Width))
		buffer.WriteString(fmt.Sprintf("\"length\":%d,", value.Board.Length))
		buffer.WriteString(fmt.Sprintf("\"foodTick\":%d,", milliseconds(value.FoodTick)))
		buffer.WriteString(fmt.Sprintf("\"moveTick\":%d", milliseconds(value.MoveTick)))
		buffer.WriteString("}")
		buffer.WriteString(",")
	}
	buffer = bytes.NewBuffer(bytes.TrimSuffix(buffer.Bytes(), []byte(",")))
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

var games = gamesType{m: make(map[string]*Game)}

var lobbyConnections = struct {
	sync.RWMutex
	m map[*websocket.Conn]bool
}{m: make(map[*websocket.Conn]bool)}

var lobbyUpdateChannel = make(chan bool, 100)

func normalizeToRange(val, min, max int) int {
	if val > max {
		val = max
	}
	if val < min {
		val = min
	}
	return val
}

func getNumericFromForm(r *http.Request, field string, def int) int {
	value, hasValue := r.Form[field]
	if !hasValue {
		return def
	}
	ret, err := strconv.Atoi(value[0])
	if err != nil {
		return def
	}
	return ret
}

func getUserData(ws *websocket.Conn) (UserConnectionMessage, error) {
	var msg UserConnectionMessage
	err := ws.ReadJSON(&msg)
	return msg, err
}

func addGame(r *http.Request) string {
	games.Lock()
	defer games.Unlock()
	var gameID string
	for {
		uid := uuid.Must(uuid.NewV4())
		gameID = uid.String()
		_, ok := games.m[gameID]
		if !ok {
			break
		}
	}
	board := Board{
		Width: normalizeToRange(
			getNumericFromForm(r, "width", 20), 1, 100),
		Length: normalizeToRange(
			getNumericFromForm(r, "length", 20), 1, 100),
		Fruits:  make([]Point, 0),
		State:   WAITING,
		Changes: make(chan Change, 100),
		End:     make(chan bool),
	}
	game := &Game{
		ID:    gameID,
		Board: board,
		Users: make([]*websocket.Conn, 0),
		UsersCount: normalizeToRange(
			getNumericFromForm(r, "players", 1), 1, 30),
		FoodTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "food_tick", 2000), 0, 120000)) * time.Millisecond,
		MoveTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "move_tick", 10), 1, 20000)) * time.Millisecond,
		EndPingChannels: make(map[string]chan bool, 30),
	}
	games.m[gameID] = game
	lobbyUpdateChannel <- true
	return gameID
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := addGame(r)

	http.Redirect(w, r, fmt.Sprintf("/game/%s", id), http.StatusFound)
}

func (g *Game) NotifyUsers() {
	newUsers := g.Users[:0]
	for _, user := range g.Users {
		err := user.WriteJSON(g.Board)
		if err != nil {
			user.Close()
		} else {
			newUsers = append(newUsers, user)
		}
	}
	g.Users = newUsers
}

func handleChanges(g *Game) {
	g.Board.run(g.MoveTick, g.FoodTick, func(b *Board) {
		g.NotifyUsers()
		if len(g.Users) == 0 {
			b.End <- true
			games.Lock()
			delete(games.m, g.ID)
			games.Unlock()
		}
	})
}

func notifyLobbyConnections() {
	for {
		<-lobbyUpdateChannel
		games.RLock()
		toSend, _ := json.Marshal(games)
		games.RUnlock()
		for key := range lobbyConnections.m {
			err := key.WriteMessage(websocket.TextMessage, toSend)
			if err != nil {
				key.Close()
				lobbyConnections.Lock()
				delete(lobbyConnections.m, key)
				lobbyConnections.Unlock()
			}
		}
	}
}

func lobbyHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	lobbyConnections.Lock()
	lobbyConnections.m[ws] = true
	lobbyConnections.Unlock()
	games.RLock()
	toSend, _ := json.Marshal(games)
	games.RUnlock()
	ws.WriteMessage(websocket.TextMessage, toSend)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/lobby.html")
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	games.RLock()
	defer games.RUnlock()
	gameID := strings.TrimPrefix(r.URL.Path, "/game/")
	game, ok := games.m[gameID]
	if !ok || game.Board.State != WAITING {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "frontend/game.html")
}

func gameConnectionHandler(w http.ResponseWriter, r *http.Request) {
	games.RLock()
	gameID := strings.TrimPrefix(r.URL.Path, "/gamews/")
	game, ok := games.m[gameID]
	games.RUnlock()
	if !ok || game.Board.State != WAITING {
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	userID := uuid.Must(uuid.NewV4()).String()
	userData, serr := getUserData(ws)
	if serr != nil {
		return
	}
	game.addUser(ws, userID, userData.Name, userData.Color)
	go pingGameConnection(ws, userID, game)
	for {
		var msg UserMoveMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}
		game.sendChange(msg, userID)
	}
}

func pingGameConnection(ws *websocket.Conn, userID string, game *Game) {
	pingTicker := time.NewTicker(time.Millisecond * 500)
	endPingChan := game.EndPingChannels[userID]
	defer pingTicker.Stop()
	for {
		select {
		case <-endPingChan:
			return
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(time.Millisecond * 250))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				game.RemoveConnection(ws, userID)
				ws.Close()
				return
			}
			ws.SetWriteDeadline(time.Time{})
		}
	}
}

func main() {
	publicServer := http.FileServer(http.Dir("frontend/public"))

	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/lobby", lobbyHandler)
	http.HandleFunc("/game/", gameHandler)
	http.HandleFunc("/gamews/", gameConnectionHandler)
	http.HandleFunc("/new_game/", newGameHandler)
	http.Handle("/public/", http.StripPrefix("/public/", publicServer))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go notifyLobbyConnections()
	fmt.Println("Waiting on " + port)
	http.ListenAndServe(":"+port, nil)
}
