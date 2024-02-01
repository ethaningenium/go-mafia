package chat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Room    string `json:"room"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Command string `json:"command"`
}

var clients = make(map[*websocket.Conn]*Message)

func handleConnections(c *gin.Context) {
	// Upgrade HTTP request to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()


	//registration client 
	var msg Message
	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Println(err)
		return
	}
	defer delete(clients, conn)
	defer notifyJoinLeave(msg.Sender, msg.Room, "left the room")

	
	clients[conn] = &Message{
		Room:    msg.Room,
		Sender:  msg.Sender,
		Message: msg.Message,
		Command: msg.Command,
	}
	notifyJoinLeave(msg.Sender, msg.Room, "joined the room")
	

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		log.Println("message: ", msg)
		if err != nil {
			log.Println(err)
		}

		// Рассылка сообщения всем клиентам в комнате
		broadcastMessage(msg)
	}
}

func notifyJoinLeave(sender, room, action string) {
	log.Println("notify: ", sender, room, action)
	// Создание и отправка уведомления о подключении/отключении
	notification := Message{
		Room:    room,
		Sender:  sender,
		Message: action,
	}

	broadcastMessage(notification)
}

func broadcastMessage(msg Message) {
	for conn, room := range clients {
		if room.Room == msg.Room {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}


