package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"handlers/databases"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id       int
	Username string
	Conn     *websocket.Conn
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	SenderId       string `json:"senderId"`
	ReceiverId     string `json:"receiverId"`
	MessageContent string `json:"messageContent"`
}

// Send
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error when upgrading the http: ", err)
		return
	}

	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		fmt.Println(messageType)
		fmt.Println(len(string(message)))
		fmt.Println(string(message))
		if err != nil {
			fmt.Println("error when reading the upcoming message : ", err)
			return
		}
		var messageStruct Message
		err = json.Unmarshal(message, &messageStruct)
		// fmt.Println()
		_, err = databases.DB.Exec(`INSERT INTO messages (sender_id,receiver_id,content,sent_at)
					VALUES (?, ?, ?, DATETIME('now'));`, messageStruct.SenderId, messageStruct.ReceiverId, messageStruct.MessageContent)
		if err != nil {
			fmt.Println("Error storing the message in DB : ", err)
		}
		fmt.Println("The message is :", messageStruct.MessageContent)
		err = conn.WriteMessage(messageType, []byte(messageStruct.MessageContent))
		if err != nil {
			fmt.Println("Error storing the message in DB : ", err)
		}

	}
}
