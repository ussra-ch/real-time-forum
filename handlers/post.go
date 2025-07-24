package handlers

import (
	"encoding/json"
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
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	query1 := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
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
