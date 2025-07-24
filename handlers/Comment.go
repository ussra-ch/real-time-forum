package handlers

import (
	"encoding/json"
	"fmt"
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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println("Error decoding comment data:", err)
		return
	}
	fmt.Println("cd:", cd)
	// Get user ID from session
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	var userID int
	err = databases.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Error retrieving session cookie:", err)
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	// Insert comment
	_, err = databases.DB.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)
	`, cd.PostID, userID, cd.Content)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	fmt.Println("Comment created successfully for post ID:", cd.PostID, userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment created successfully",
	})
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func FetchCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}

	rows, err := databases.DB.Query(`
		SELECT id, post_id, user_id, content, created_at 
		FROM comments 
		WHERE post_id = ? 
		ORDER BY created_at ASC`, postID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
		if err != nil {
			http.Error(w, "Error scanning comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
