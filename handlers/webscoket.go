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
	SenderId       float64 `json:"senderId"`
	ReceiverId     float64 `json:"receiverId"`
	MessageContent string  `json:"messageContent"`
	Seen           bool    `json:"seen"`
	Type           string  `json:"type"`
	// Notifications  int     `json:notifications`
}

type Notification struct {
	Type string `json:"type"` // "notification"
	// SenderId    int    `json:"senderId"`
	UnreadCount int `json:"unreadCount"`
}

var ConnectedUsers = make(map[float64]*websocket.Conn)
var OpenedConversations = make(map[float64]map[float64]bool)

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
	sendUnreadNotifications(userId, conn)
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(ConnectedUsers, float64(userId))
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

		var isConversationOpened bool
		var toolMap map[string]interface{}
		decoder := json.NewDecoder(message)
		_ = decoder.Decode(&toolMap)

		if typeValue, ok := toolMap["type"].(string); ok {
			fmt.Println("Type value is :", typeValue)
			if typeValue == "OpenConversation" || typeValue == "CloseConversation" {
				fmt.Println("\nTool map is :\n\n", toolMap)
				mu.Lock()
				conversationOpened(toolMap["senderId"].(float64), toolMap["receiverId"].(float64), toolMap["type"].(string))
				if (typeValue == "OpenConversation" ){
					fmt.Println("inside openConversation::::::::::::::::")
					updateSeenValue(int(toolMap["receiverId"].(float64)),int(toolMap["senderId"].(float64)))
				}
				sendUnreadNotifications(userId, conn)
				mu.Unlock()
			}
			if typeValue == "message" {
				messageStruct.SenderId = toolMap["senderId"].(float64)
				messageStruct.ReceiverId = toolMap["receiverId"].(float64)
				messageStruct.Type = toolMap["type"].(string)
				messageStruct.MessageContent = toolMap["messageContent"].(string)
				messageHandler(messageStruct)
				isConversationOpened = OpenedConversations[toolMap["receiverId"].(float64)][toolMap["senderId"].(float64)]
			}
		}

		if len(messageStruct.MessageContent) > 0 {
			fmt.Println("isConversationOpened :", isConversationOpened)
			if isConversationOpened {

				//update seen = true
				updateSeenValue(int(messageStruct.ReceiverId), int(messageStruct.SenderId))

				//update notification's value
				// messageStruct.Notifications = unreadMessages(int(messageStruct.ReceiverId))

				// fmt.Println("Notifications in backend is :", messageStruct.Notifications)

				//send message to the receiver
				Message, err := json.Marshal(messageStruct)
				if err != nil {
					fmt.Println("error in the messageHandler")
				}
				err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, []byte(Message))
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

				sendUnreadNotifications(int(messageStruct.ReceiverId), ConnectedUsers[messageStruct.ReceiverId])
			} else {
				//update notification's value
				// messageStruct.Notifications = unreadMessages(int(messageStruct.ReceiverId))
				// fmt.Println("Notifications in backend is :", messageStruct.Notifications)
				sendUnreadNotifications(int(messageStruct.ReceiverId), ConnectedUsers[messageStruct.ReceiverId])
				// notifications, err := json.Marshal(notifs)
				// if err != nil {
				// 	fmt.Println("error in the messageHandler")
				// }
				// err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, []byte(notifications))
				// if err != nil {
				// 	fmt.Println("Error sending message:", err)
				// }
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
		var seen bool

		if err := rows.Scan(&id, &sender_id, &userId, &content, &time, &seen); err != nil {
			fmt.Println("error in a message", err)
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
	if _, exists := ConnectedUsers[float64(userId)]; !exists {
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
	ConnectedUsers[float64(userId)] = conn
}

func userOffline(userId int, conn *websocket.Conn) {
	delete(ConnectedUsers, float64(userId))
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

func conversationOpened(senderId, receiverId float64, typeValue string) {
	if OpenedConversations[senderId] == nil {
		OpenedConversations[senderId] = make(map[float64]bool)
	}
	fmt.Printf("Type for : %f, %f, %s\n", senderId, receiverId, typeValue)
	if typeValue == "closeConversation" {
		OpenedConversations[senderId][receiverId] = false
	} else {
		OpenedConversations[senderId][receiverId] = true
		// updateSeenValue(messageStruct)
		// sendNotification(&messageStruct)
	}

	// bothOpen := isMutuallyOpen(messageStruct.SenderId, messageStruct.ReceiverId)

	// fmt.Printf("Conversation between %d and %d open status: %v\n",
	// 	messageStruct.SenderId, messageStruct.ReceiverId, bothOpen)

	// return bothOpen

}

func messageHandler(messageStruct Message) {
	_, err := databases.DB.Exec(`INSERT INTO messages (sender_id,receiver_id,content,seen )
					VALUES (?, ?, ?, ?);`, messageStruct.SenderId, messageStruct.ReceiverId, messageStruct.MessageContent, false)
	if err != nil {
		fmt.Println("Error storing the message in DB : ", err)
	}
	fmt.Println("the message is stored in the database")
	// return Message
}

func updateSeenValue(receiverId, senderId int) {
	fmt.Println("updaaaate seeeeeen")
	query := `UPDATE messages
	SET seen = 1
	WHERE messages.sender_id = ? AND messages.receiver_id = ?;`
	_, err := databases.DB.Exec(query, senderId, receiverId)
	if err != nil {
		fmt.Println("eror when changing the seen value in database, in updateSeenValue function")
	}

	//update l notifications dyal receiver
	// var totalUnread int
	// err = databases.DB.QueryRow(`
	// 				SELECT COUNT(*) FROM messages
	// 				WHERE receiver_id = ? AND seen = false;
	// 			`, messageStruct.ReceiverId).Scan(&totalUnread)
	// if err != nil {
	// 	fmt.Println("Error fetching total unread count:", err)
	// }

	// Send updated total notification count
	// notif := Notification{
	// 	Type:        "notifications",
	// 	SenderId:    messageStruct.SenderId,
	// 	UnreadCount: totalUnread,
	// }
	// notifBytes, _ := json.Marshal(notif)
	// err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, notifBytes)
	// if err != nil {
	// 	fmt.Println("Error sending total notification count:", err)
	// }

	fmt.Println("Seen value has been updated")
}

func unreadMessages(receiverId int) int {
	var unreadCount int
	err := databases.DB.QueryRow(`
					SELECT COUNT(*) FROM messages
					WHERE receiver_id = ? AND seen = false;
				`, receiverId).Scan(&unreadCount)
	if err != nil {
		fmt.Println("Error fetching unread count:", err)
	}

	return unreadCount

	// notif := Notification{
	// 	Type:        "notification",
	// 	UnreadCount: unreadCount,
	// }
	// fmt.Println("UnreadCount is :", unreadCount)
	// notifBytes, _ := json.Marshal(notif)
	// err = ConnectedUsers[messageStruct.ReceiverId].WriteMessage(websocket.TextMessage, notifBytes)
	// if err != nil {
	// 	fmt.Println("Error sending notification:", err)
	// }
	// fmt.Println("notification count is sent to the receiver ")
}

// func isMutuallyOpen(user1, user2 int) bool {
// 	return OpenedConversations[user1][user2] && OpenedConversations[user2][user1]
// }


func sendUnreadNotifications(userId int, conn *websocket.Conn){
	fmt.Println("USERRRRR ID :", userId)
	notifs := Notification{
		Type:        "unreadMessage",
		UnreadCount: unreadMessages(userId),
	}
	Notifs, err := json.Marshal(notifs)
	if err != nil {
		fmt.Println("error in the messageHandler")
	}
	// fmt.Println("user id is ;", userId)
	// fmt.Println("unread messages are : ", notifs.UnreadCount)
	conn.WriteMessage(websocket.TextMessage, Notifs)
}