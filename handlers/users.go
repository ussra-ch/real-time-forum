package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"handlers/databases"
)

var mu sync.Mutex

type User struct {
	Nickname string         `json:"nickname"`
	UserId   int            `json:"userId"`
	Photo    sql.NullString `json:"photo"`
	Status   string         `json:"status"`
}

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	loggedIn, userID := IsLoggedIn(r)
	if !loggedIn {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn":    false,
			"nickname":    nil,
			"onlineUsers": []string{},
		})
		return
	}

	var myNickname string
	err := databases.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&myNickname)
	if err != nil {
		log.Println("Error getting nickname:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rows, err := databases.DB.Query(`
		SELECT u.nickname, u.id, u.photo
		FROM users u
		
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type User struct {
		Nickname string         `json:"nickname"`
		UserId   int            `json:"userId"`
		Photo    sql.NullString `json:"photo"`
		Status   string         `json:"status"`
		Time     time.Time      `json:"time"`
	}
	var onlineUsers []User
	for rows.Next() {
		var photo sql.NullString
		var nickname string
		var userId int
		q := `SELECT sent_at  FROM messages 
    	WHERE (sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)
    	ORDER BY sent_at DESC
    	LIMIT 1`
		row, _ := databases.DB.Query(q, userId, userID, userId, userID)
		var T time.Time
		for row.Next() {
			var time time.Time

			if err := row.Scan(&time); err != nil {
				// fmt.Println("error in a message")
			}
			T = time
		}
		if err := rows.Scan(&nickname, &userId, &photo); err != nil {
			log.Fatal("error", err)
		}
		// fmt.Println(ConnectedUsers[userId])
		mu.Lock()
		_, exists := ConnectedUsers[float64(userId)]
		mu.Unlock()
		if exists {
			mu.Lock()
			UsersStatus[userId] = "online"
			mu.Unlock()
		}
		// else {
		// 	mu.Lock()
		// 	// fmt.Println(userId)
		// 	UsersStatus[userId] = "offline"
		// 	mu.Unlock()

		// }
		mu.Lock()
		onlineUsers = append(onlineUsers, User{Nickname: nickname, UserId: userId, Photo: photo, Status: UsersStatus[userId], Time: T})
		mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	// fmt.Println(offlineUsers)
	mu.Lock()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn":    true,
		"nickname":    myNickname,
		"onlineUsers": onlineUsers,
		"UserId":      userID,
		"status":      UsersStatus[userID],
	})
	mu.Unlock()
}
