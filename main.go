package main

import (
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
	Board      Board
	UsersCount int
	Users      []*websocket.Conn
	FoodTick   time.Duration
	MoveTick   time.Duration
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
	g.Board.addSnake(userID, name, color, 3)
	g.Board.Unlock()
	if g.UsersCount == len(g.Users) {
		g.Board.State = PREPARING
		g.NotifyUsers()
		go func() {
			time.Sleep(5 * time.Second)
			handleChanges(g)
		}()
	}
}

var Games = struct {
	sync.RWMutex
	m map[string]*Game
}{m: make(map[string]*Game)}

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

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	Games.Lock()
	defer Games.Unlock()
	r.ParseForm()

	board := Board{
		Width: normalizeToRange(
			getNumericFromForm(r, "width", 20), 1, 100),
		Length: normalizeToRange(
			getNumericFromForm(r, "length", 20), 1, 100),
		Fruits:  make([]Point, 0),
		State:   WAITING,
		Changes: make(chan Change),
		End:     make(chan bool),
	}

	var boardID string
	for {
		uid := uuid.Must(uuid.NewV4())
		boardID = uid.String()
		_, ok := Games.m[boardID]
		if !ok {
			break
		}
	}
	game := &Game{
		Board: board,
		Users: make([]*websocket.Conn, 0),
		UsersCount: normalizeToRange(
			getNumericFromForm(r, "players", 1), 1, 30),
		FoodTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "food_tick", 2000), 0, 20000)) * time.Millisecond,
		MoveTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "move_tick", 10), 1, 20000)) * time.Millisecond,
	}
	Games.m[boardID] = game
	http.Redirect(w, r, fmt.Sprintf("/game/%s", boardID), http.StatusFound)
}

func (g *Game) NotifyUsers() {
	for _, user := range g.Users {
		err := user.WriteJSON(g.Board)
		if err != nil {
			user.Close()
			g.Board.End <- true
		}
	}
}

func handleChanges(g *Game) {
	g.Board.run(g.MoveTick, g.FoodTick, func(b *Board) {
		g.NotifyUsers()
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	Games.RLock()
	defer Games.RUnlock()
	gameID := strings.TrimPrefix(r.URL.Path, "/game/")
	game, ok := Games.m[gameID]
	if !ok || game.Board.State != WAITING {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "frontend/game.html")
}

func gameConnectionHandler(w http.ResponseWriter, r *http.Request) {
	Games.RLock()
	gameID := strings.TrimPrefix(r.URL.Path, "/gamews/")
	game, ok := Games.m[gameID]
	Games.RUnlock()
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
	for {
		var msg UserMoveMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}
		game.sendChange(msg, userID)
	}
}

func main() {
	publicServer := http.FileServer(http.Dir("frontend/public"))

	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/game/", gameHandler)
	http.HandleFunc("/gamews/", gameConnectionHandler)
	http.HandleFunc("/new_game/", newGameHandler)
	http.Handle("/public/", http.StripPrefix("/public/", publicServer))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Waiting on " + port)
	http.ListenAndServe(":"+port, nil)
}
