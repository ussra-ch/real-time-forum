package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"handlers/databases"
)

type PostData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var pd PostData
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("Received post data:", pd)
	cookie, err := r.Cookie("session")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	query1 := `SELECT user_id FROM sessions WHERE session = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
	fmt.Println("User ID from session:", err)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	query := `
    INSERT INTO posts (title, content, interest, user_id)
    VALUES (?, ?, ?, ?)
`
	_, err = databases.DB.Exec(query, pd.Title, pd.Description, strings.Join(pd.Topics, ","), userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data received successfully",
		"title":    pd.Title,
		"content":  pd.Description,
		"interest": strings.Join(pd.Topics, ","),
	})
}
func FetchPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT id, user_id, content, title, interest, photo, created_at FROM posts`
	rows, err := databases.DB.Query(query)
	if err != nil {
		log.Println("Error fetching posts:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var userID int
		var content, title, interest, photo string
		var createdAt string
		if err := rows.Scan(&id, &userID, &content, &title, &interest, &photo, &createdAt); err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		post := map[string]interface{}{
			"id":        id,
			"user_id":   userID,
			"content":   content,
			"title":     title,
			"interest":  interest,
			"photo":     photo,
			"created_at": createdAt,
		}
		posts = append(posts, post)
	}
	fmt.Println("Fetched posts:", posts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
