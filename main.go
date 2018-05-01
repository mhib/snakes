package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type gamesType struct {
	sync.RWMutex
	m map[string]*Game
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func milliseconds(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1000000
}

func (gameMap *gamesType) MarshalJSON() ([]byte, error) {
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
		buffer.WriteString(fmt.Sprintf("\"moveTick\":%d,", milliseconds(value.MoveTick)))
		buffer.WriteString(fmt.Sprintf("\"bots\":%d", len(value.Bots)))
		buffer.WriteString("}")
		buffer.WriteString(",")
	}
	buffer = bytes.NewBuffer(bytes.TrimSuffix(buffer.Bytes(), []byte(",")))
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

func initialGameInfo(g *Game) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteString("{")
	buffer.WriteString(fmt.Sprintf("\"width\":%d,", g.Board.Width))
	buffer.WriteString(fmt.Sprintf("\"length\":%d,", g.Board.Length))
	if len(g.Users) == g.UsersCount-1 {
		buffer.WriteString("\"isLast\":true")
	} else {
		buffer.WriteString("\"isLast\":false")
	}
	buffer.WriteString("}")
	return buffer.Bytes()
}

var games = gamesType{m: make(map[string]*Game)}

var lobby = NewLobby([]byte("[]"))

var lobbyUpdateChannel = make(chan bool, 1024)

func notifyLobby() {
	for {
		<-lobbyUpdateChannel
		games.RLock()
		toSend, _ := json.Marshal(&games)
		games.RUnlock()
		lobby.broadcast <- toSend
	}
}

var removeGameChannel = make(chan *Game, 1024)

func removeGameWorker() {
	for {
		game := <-removeGameChannel
		games.Lock()
		delete(games.m, game.ID)
		games.Unlock()
	}
}

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

func getBoolFromForm(r *http.Request, field string, def bool) bool {
	value, hasValue := r.Form[field]
	if !hasValue {
		return def
	}
	return value[0] == "on"
}

type userConnectionMessage struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func getUserData(ws *websocket.Conn) (userConnectionMessage, error) {
	var msg userConnectionMessage
	err := ws.ReadJSON(&msg)
	return msg, err
}

func addBots(b *Board, r *http.Request) []AI {
	ret := make([]AI, 0)
	nearestFruitCount := normalizeToRange(
		getNumericFromForm(r, "nearestFruitBots", 0), 0, 4)
	randomCount := normalizeToRange(
		getNumericFromForm(r, "randomMoveBots", 0), 0, 4)
	lazyCount := normalizeToRange(
		getNumericFromForm(r, "lazyBots", 0), 0, 4)

	for i := 1; i <= nearestFruitCount; i++ {
		name := fmt.Sprintf("Nearest-Fruit-Bot#%d", i)
		ai := NewNearestFoodAI(b.Changes, name)
		b.AddSnake(name, name, "#999999", 3)
		ret = append(ret, ai)
	}
	for i := 1; i <= randomCount; i++ {
		name := fmt.Sprintf("Random-Bot#%d", i)
		ai := NewRandomMoveAI(b.Changes, name)
		b.AddSnake(name, name, "#999999", 3)
		ret = append(ret, ai)
	}
	for i := 1; i <= lazyCount; i++ {
		name := fmt.Sprintf("Lazy-Bot#%d", i)
		ai := NewLazyAI(b.Changes, name)
		b.AddSnake(name, name, "#999999", 3)
		ret = append(ret, ai)
	}
	return ret
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
	usersCount := normalizeToRange(getNumericFromForm(r, "players", 1), 1, 30)
	board := Board{
		Width: normalizeToRange(
			getNumericFromForm(r, "width", 20), 1, 100),
		Length: normalizeToRange(
			getNumericFromForm(r, "length", 20), 1, 100),
		Fruits:          make([]Point, 0),
		State:           WAITING,
		EndOnLastPlayer: getBoolFromForm(r, "endOnLastPlayer", false) && usersCount > 1,
		Tick:            0,
		Changes:         make(chan Change, 100),
		End:             make(chan bool, 1),
	}
	bots := addBots(&board, r)
	game := &Game{
		ID:         gameID,
		Board:      board,
		Users:      make(map[string]*Client),
		UsersCount: usersCount,
		FoodTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "foodTick", 2000), 0, 120000)) * time.Millisecond,
		MoveTick: time.Duration(normalizeToRange(
			getNumericFromForm(r, "moveTick", 10), 1, 20000)) * time.Millisecond,
		Broadcast:          make(chan []byte),
		Register:           make(chan *Client),
		Unregister:         make(chan *Client),
		ChangeStateChannel: lobbyUpdateChannel,
		DisposeChannel:     removeGameChannel,
		Bots:               bots,
	}
	games.m[gameID] = game
	go game.Run()
	return gameID
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := addGame(r)

	http.Redirect(w, r, fmt.Sprintf("/game/%s", id), http.StatusFound)
}

func lobbyHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	lobby.register <- NewLobbyConnection(ws, lobby)
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
	t, _ := template.ParseFiles("frontend/game.html")
	t.Execute(w, template.HTML(initialGameInfo(game)))
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
	userID := uuid.Must(uuid.NewV4()).String()
	userData, serr := getUserData(ws)
	if serr != nil {
		return
	}
	game.Register <- NewClient(ws, userData.Name, userData.Color, userID, game)
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
	go lobby.Run()
	go notifyLobby()
	go removeGameWorker()
	fmt.Println("Waiting on " + port)
	http.ListenAndServe(":"+port, nil)
}
