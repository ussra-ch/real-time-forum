package handlers

import (
	"encoding/json"
	"html"
	"log"
	"net/http"

	"handlers/databases"
)

type CommentData struct {
	PostID  string `json:"post_id"`
	Content string `json:"comment"`
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	var cd CommentData
	if err := json.NewDecoder(r.Body).Decode(&cd); err != nil {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}
	if cd.Content == "" {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}
	var exists bool
	err := databases.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", cd.PostID).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}
	// Get user ID from session
	cookie, err := r.Cookie("session")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	var userID int
	err = databases.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, cookie.Value).Scan(&userID)
	if err != nil {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorr)
		return
	}
	// Insert comment
	_, err = databases.DB.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)
	`, cd.PostID, userID, html.EscapeString(cd.Content))
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment created successfully",
	})
}

type Comment struct {
	ID        int
	Content   string
	CreatedAt string
	UserID    string
	PostID    string
	Name      string
}

func FetchCommentsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := databases.DB.Query(`
		SELECT 
			comments.id,
			comments.content,
			comments.created_at,
			comments.user_id,
			comments.post_id,
			users.nickname
		FROM comments
		JOIN users ON comments.user_id = users.id
		ORDER BY comments.created_at DESC;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.PostID, &c.Name)
		if err != nil {
			log.Println("Error scanning comment:", err)
			continue
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
