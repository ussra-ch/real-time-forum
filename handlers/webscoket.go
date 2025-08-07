package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	Seen           bool   `json:"seen"`
	IsOpen bool `json:"isOpen"`
	// ClientStatus   bool   `json:"clientStatus"`
}

type Notification struct {
    Type       string `json:"type"` // "notification"
    SenderId   int    `json:"senderId"`
    UnreadCount int   `json:"unreadCount"`
}

var ConnectedUsers = make(map[int]*websocket.Conn)
var OpenedConversations = make(map[int]map[int]bool)

// openedConversations = make(map[int]int)

// Send
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error when upgrading the http: ", err)
		return
	}
	_, userId := IsLoggedIn(r)
	mu.Lock()
	if _, exists := ConnectedUsers[userId]; !exists {
		newUser := make(map[string]interface{})
		newUser["type"] = "online"
		newUser["userId"] = userId
		toSend, err := json.Marshal(newUser)
		if err != nil {
			fmt.Println("error when sending the user's status : ", err)
		}
		for _, value := range ConnectedUsers {
			value.WriteMessage(websocket.TextMessage, []byte(toSend))
		}
	}

	ConnectedUsers[userId] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(ConnectedUsers, userId)
		conn.Close()
		mu.Unlock()
	}()

	for {
		_, message, err := conn.NextReader()
		if message == nil {
			fmt.Println("11")
			mu.Lock()
			delete(ConnectedUsers, userId)
			UsersStatus[userId] = "offline"
			newUser := make(map[string]interface{})
			newUser["type"] = "online"
			newUser["userId"] = userId
			toSend, _ := json.Marshal(newUser)
			for _, value := range ConnectedUsers {
				value.WriteMessage(websocket.TextMessage, []byte(toSend))
			}
			conn.Close()

			mu.Unlock()

		}
		if err != nil {
			fmt.Println("error when reading the upcoming message : ", err)
			return
		}

		var messageStruct Message
		decoder := json.NewDecoder(message)
		_ = decoder.Decode(&messageStruct)
		mu.Lock()
		if OpenedConversations[messageStruct.SenderId] == nil {
			OpenedConversations[messageStruct.SenderId] = make(map[int]bool)
		}
		OpenedConversations[messageStruct.SenderId][messageStruct.ReceiverId] = messageStruct.IsOpen
		// fmt.Println("conver :",OpenedConversations[messageStruct.SenderId][messageStruct.ReceiverId])
		mu.Unlock()

		if len(messageStruct.MessageContent) > 0{
		messageobj := make(map[string]interface{})
		messageobj["type"] = "message"
		messageobj["SenderId"] = messageStruct.SenderId
		messageobj["ReceiverId"] = messageStruct.ReceiverId
		messageobj["content"] = messageStruct.MessageContent
		messageobj["seen"] = messageStruct.Seen
		Messag, err := json.Marshal(messageobj)
		if err != nil {
			fmt.Println("erooooooor f decoder")
		}

		_, err = databases.DB.Exec(`INSERT INTO messages (sender_id,receiver_id,content,seen )
					VALUES (?, ?, ?, ?);`, messageStruct.SenderId, messageStruct.ReceiverId, messageStruct.MessageContent, false)
		if err != nil {
			fmt.Println("Error storing the message in DB : ", err)
		}
		// fmt.Println("2222")
		if ConnectedUsers[messageStruct.ReceiverId] != nil  && OpenedConversations[messageStruct.ReceiverId][messageStruct.SenderId]{
			err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, []byte(Messag))
			if err != nil {
				fmt.Println("Error storing the message in DB : ", err)
			}
			query := `UPDATE messages
						SET seen = true
						WHERE messages.sender_id = ? AND messages.receiver_id = ?;`
			_, err = databases.DB.Exec(query, messageStruct.SenderId, messageStruct.ReceiverId)
			if err != nil {
				fmt.Println("eroooooor fach kanbdlo seen=true ")
			}
		} else {
			fmt.Println("dkhaaal l else")

			if ConnectedUsers[messageStruct.ReceiverId] != nil {
				var unreadCount int
				err := databases.DB.QueryRow(`
					SELECT COUNT(*) FROM messages
					WHERE receiver_id = ? AND seen = false;
				`, messageStruct.ReceiverId).Scan(&unreadCount)
				if err != nil {
					fmt.Println("Error fetching unread count:", err)
				}
		
				notif := Notification{
					Type:        "notification",
					SenderId:    messageStruct.SenderId,
					UnreadCount: unreadCount,
				}
		
				notifBytes, _ := json.Marshal(notif)
				err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, notifBytes)
				if err != nil {
					fmt.Println("Error sending notification:", err)
				}
			}
		}

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
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	senderID := r.URL.Query().Get("sender")

	offset, err1 := strconv.Atoi(offsetStr)
	limit, err2 := strconv.Atoi(limitStr)
	fmt.Println(offset, limit)

	if err1 != nil || err2 != nil || limit <= 0 {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
    SELECT * FROM messages 
    WHERE (sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)
    ORDER BY sent_at DESC
    LIMIT %d OFFSET %d;`, limit, offset)

	rows, err := databases.DB.Query(query, userId, senderID, userId, senderID)
	if err != nil {
		fmt.Println("error geting messages from db : ", err)
	}
	var messages []map[string]interface{}
	for rows.Next() {
		var id, userId, sender_id int
		var content string
		var time time.Time

		if err := rows.Scan(&id, &sender_id, &userId, &content, &time); err != nil {
			// fmt.Println("error in a message")
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
