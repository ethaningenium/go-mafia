package timer

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Phases = map[string]string{"waiting": "chat", "chat": "action", "action": "choice", "choice" : "chat"}

type Response struct {
	Phase string `json:"phase"`
	Time  int    `json:"time"`
}

type Message struct {
	Room    string `json:"room"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type Game struct {
	Owner   *websocket.Conn
	Players map[*websocket.Conn]string
	Timer   *time.Ticker // Глобальный таймер
	Current int
	Phase   string
}

 func createGame(owner *websocket.Conn) *Game {
	var players = make(map[*websocket.Conn]string)

	players[owner] = "owner"
	// Инициализация таймера с периодом 1 секунда
	ticker := time.NewTicker(1 * time.Second)

	// Возвращаем созданный объект Game
	return &Game{
		Owner:   owner,
		Players: players,
		Timer:   ticker,
		Current: 0,
		Phase:   "waiting",
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var games = make(map[string]*Game)

func handleConnections(c *gin.Context) {
	// Upgrade HTTP request to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	// Регистрация клиента
	var msg Message
	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Println(err)
		return
	}

	var game *Game

	// Проверка наличия игры в карте
	value, exists := games[msg.Room]
	if exists {
		value.Players[conn] = msg.Sender
	} else {
		// Если игры нет, создаем новую
		game = createGame(conn)
		games[msg.Room] = game

		// Запускаем горутину для управления глобальным таймером
		go manageGlobalTimer(game)
		
	}
	for  {
		time.Sleep(10 * time.Second)
		var msg Message
		err := conn.ReadJSON(&msg)
		log.Println("message: ", msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
	
}

func manageGlobalTimer(game *Game) {
	for range game.Timer.C {
		// Обновляем текущее время
		if game.Current >= 30 {
			game.Phase = Phases[game.Phase]
			game.Current = 0
		}else{
			game.Current++
		}
		
		if len(game.Players) == 0 {
			game.Phase = "waiting"
			game.Current = 0
			game.Timer.Stop()
		}
		// Отправляем всем клиентам обновленные данные
		broadcastTimeUpdate(game)
	}
}

func broadcastTimeUpdate(game *Game) {
	fmt.Println("Players", game.Players)
	// Создаем сообщение с обновленным временем
	response := Response{
		Phase: game.Phase,
		Time:  game.Current,
	}

	// Создаем временный срез для соединений
	var connectionsToRemove []*websocket.Conn

	// Рассылаем сообщение всем подключенным клиентам
	for conn := range game.Players {
		if err := conn.WriteJSON(response); err != nil {
			// Если произошла ошибка, удаляем соединение
			log.Println(err)
			// Добавляем соединение в срез для последующего удаления
			connectionsToRemove = append(connectionsToRemove, conn)
		}
	}

	// Удаляем соединения после завершения итерации
	for _, conn := range connectionsToRemove {
		delete(game.Players, conn)
		conn.Close()
	}
}
