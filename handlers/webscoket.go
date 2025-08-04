package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	SenderId       int    `json:"senderId"`
	ReceiverId     int    `json:"receiverId"`
	MessageContent string `json:"messageContent"`
	// ClientStatus   bool   `json:"clientStatus"`
}

var ConnectedUsers = make(map[int]*websocket.Conn)

// openedConversations = make(map[int]int)

// Send
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error when upgrading the http: ", err)
		return
	}
	_, userId := IsLoggedIn(r)

	if _, exists := ConnectedUsers[userId]; !exists {
		newUser := make(map[string]interface{})
		newUser["type"] = "online"
		newUser["userId"] = userId
		toSend, err := json.Marshal(newUser)
		if err != nil {
			fmt.Println("error when sending the user's status : ", err)
		}
		// fmt.Println("11")
		for _, value := range ConnectedUsers {

			fmt.Println("dkhal l loop bach ysift status dluser")
			value.WriteMessage(websocket.TextMessage, []byte(toSend))
		}
	}
	ConnectedUsers[userId] = conn

	defer func(){
			delete(ConnectedUsers, userId)
		 conn.Close()
	}()

	for {
		_, message, err := conn.NextReader()
		if err != nil {
			fmt.Println("error when reading the upcoming message : ", err)
			return
		}

		var messageStruct Message
		decoder := json.NewDecoder(message)
		err = decoder.Decode(&messageStruct)

		messageobj := make(map[string]interface{})
		messageobj["type"] = "message"
		messageobj["SenderId"] = messageStruct.SenderId
		messageobj["ReceiverId"] = messageStruct.ReceiverId
		messageobj["content"] = messageStruct.MessageContent
		Messag, err := json.Marshal(messageobj)
		if err != nil {
			fmt.Println("erooooooor f decoder")
		}

		_, err = databases.DB.Exec(`INSERT INTO messages (sender_id,receiver_id,content)
					VALUES (?, ?, ?);`, messageStruct.SenderId, messageStruct.ReceiverId, messageStruct.MessageContent)
		if err != nil {
			fmt.Println("Error storing the message in DB : ", err)
		}
		// fmt.Println("2222")
		if ConnectedUsers[messageStruct.ReceiverId] != nil {

			err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, []byte(Messag))
			if err != nil {
				fmt.Println("Error storing the message in DB : ", err)
			}
		} else {
			// fmt.Println("4444")
		}

	}

}

func FetchMessages(w http.ResponseWriter, r *http.Request) {
	// fetch data
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	_, userId := IsLoggedIn(r)
	query := `SELECT * from messages WHERE sender_id = ? OR receiver_id = ?
	ORDER BY sent_at DESC
	LIMIT 10;`
	rows, err := databases.DB.Query(query, userId, userId)
	if err != nil {
		fmt.Println("error geting messages from db : ", err)
	}
	var messages []map[string]interface{}
	for rows.Next() {
		var id, userId, sender_id int
		var content string
		var time time.Time

		if err := rows.Scan(&id, &sender_id, &userId, &content, &time); err != nil {
			fmt.Println("error in a message")
		}
		message := map[string]interface{}{
			"id":        id,
			"sender_id": sender_id,
			"userId":    userId,
			"content":   content,
			"time":      time,
		}
		messages = append(messages, message)
		// fmt.Println(id, sender_id, userId, content, time)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
