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
	IsOpen         bool   `json:"isOpen"`
	Type           string `json:"type"`
	// ClientStatus   bool   `json:"clientStatus"`
}

type Notification struct {
	Type        string `json:"type"` // "notification"
	SenderId    int    `json:"senderId"`
	UnreadCount int    `json:"unreadCount"`
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
	fmt.Println("User connected:", userId)
	mu.Lock()
	broadcastUserStatus(conn, userId)
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
			mu.Lock()
			userOffline(userId, conn)
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
		isOpen := conversationOpened(messageStruct)
		mu.Unlock()

		if len(messageStruct.MessageContent) > 0 {
			Messag := messageHandler(messageStruct)
			fmt.Println("users state is : ", isOpen)
			if ConnectedUsers[messageStruct.ReceiverId] != nil &&
				isOpen {
				fmt.Println("the message content is :", messageStruct.MessageContent)
				updateSeenValue(messageStruct)
				sendNotification(messageStruct)
				err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, []byte(Messag))
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
			} else {
				if ConnectedUsers[messageStruct.ReceiverId] != nil {
					fmt.Println("in the else of the first if in websocket handler inside the loop")
					sendNotification(messageStruct)
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
	// fmt.Println(offset, limit)

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

func broadcastUserStatus(conn *websocket.Conn, userId int) {
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
	fmt.Println("User status broadcasted")
	ConnectedUsers[userId] = conn
}

func userOffline(userId int, conn *websocket.Conn) {
	delete(ConnectedUsers, userId)
	UsersStatus[userId] = "offline"
	newUser := make(map[string]interface{})
	newUser["type"] = "offline" ///////////////////////// hna rh kant online, makaynch dalil niit ms rh jatni khas tjun offline
	newUser["userId"] = userId
	toSend, _ := json.Marshal(newUser)
	for _, value := range ConnectedUsers {
		value.WriteMessage(websocket.TextMessage, []byte(toSend))
	}
	conn.Close()
	fmt.Println("status changed to offline of user :", userId)
}

func conversationOpened(messageStruct Message) bool {
	if OpenedConversations[messageStruct.SenderId] == nil {
		OpenedConversations[messageStruct.SenderId] = make(map[int]bool)
	}
	fmt.Println("messageStruct.IsOpen : ", messageStruct.IsOpen)
	if messageStruct.Type == "closeConversation" {
		OpenedConversations[messageStruct.SenderId][messageStruct.ReceiverId] = false
	} else {
		OpenedConversations[messageStruct.SenderId][messageStruct.ReceiverId] = true
	}

	bothOpen := isMutuallyOpen(messageStruct.SenderId, messageStruct.ReceiverId)

	fmt.Printf("Conversation between %d and %d open status: %v\n",
		messageStruct.SenderId, messageStruct.ReceiverId, bothOpen)

	return bothOpen

}

func messageHandler(messageStruct Message) []byte {
	messageobj := make(map[string]interface{})
	messageobj["type"] = "message"
	messageobj["SenderId"] = messageStruct.SenderId
	messageobj["ReceiverId"] = messageStruct.ReceiverId
	messageobj["content"] = messageStruct.MessageContent
	messageobj["seen"] = false
	Messag, err := json.Marshal(messageobj)
	if err != nil {
		fmt.Println("error in the messageHandler")
	}

	_, err = databases.DB.Exec(`INSERT INTO messages (sender_id,receiver_id,content,seen )
					VALUES (?, ?, ?, ?);`, messageStruct.SenderId, messageStruct.ReceiverId, messageStruct.MessageContent, false)
	if err != nil {
		fmt.Println("Error storing the message in DB : ", err)
	}
	fmt.Println("message Handler worked as we want")
	return Messag
}

func updateSeenValue(messageStruct Message) {
	query := `UPDATE messages
	SET seen = true
	WHERE messages.sender_id = ? AND messages.receiver_id = ?;`
	_, err := databases.DB.Exec(query, messageStruct.SenderId, messageStruct.ReceiverId)
	if err != nil {
		fmt.Println("eror when changing the seen value in database, in updateSeenValue function")
	}

	//update l notifications dyal receiver
	var totalUnread int
	err = databases.DB.QueryRow(`
					SELECT COUNT(*) FROM messages
					WHERE receiver_id = ? AND seen = false;
				`, messageStruct.ReceiverId).Scan(&totalUnread)
	if err != nil {
		fmt.Println("Error fetching total unread count:", err)
	}

	// Send updated total notification count
	notif := Notification{
		Type:        "notifications",
		SenderId:    messageStruct.SenderId,
		UnreadCount: totalUnread,
	}
	notifBytes, _ := json.Marshal(notif)
	err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, notifBytes)
	if err != nil {
		fmt.Println("Error sending total notification count:", err)
	}

	fmt.Println("message sent to the receiver in his conversation")
}

func sendNotification(messageStruct Message) {
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
		UnreadCount: unreadCount,
	}
	fmt.Println("UnreadCount is :", unreadCount)
	notifBytes, _ := json.Marshal(notif)
	err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, notifBytes)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
	fmt.Println("notification count is sent to the receiver ")
}

func isMutuallyOpen(user1, user2 int) bool {
	return OpenedConversations[user1][user2] && OpenedConversations[user2][user1]
}
