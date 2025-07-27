package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	topics := r.Form["topics"]

	cookie, err := r.Cookie("session")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	query1 := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	var filename string
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		filename = fmt.Sprintf("static/uploads/%d_%s", time.Now().UnixNano(), handler.Filename)
		dst, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Cannot save photo", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to save photo", http.StatusInternalServerError)
			return
		}
	} else {
		filename = ""
	}

	query := `
		INSERT INTO posts (title, content, interest, user_id, photo)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := databases.DB.Exec(query, title, description, strings.Join(topics, ","), userID, filename)
	if err != nil {
		log.Println("Insert post error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		log.Println("Error getting inserted post ID:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Data received successfully",
		"title":    title,
		"content":  description,
		"interest": strings.Join(topics, ","),
		"photo":    filename,
		"post_id":  postID,
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
		var id, userID int
		var content, title, interest string
		var photo sql.NullString
		var createdAt string

		if err := rows.Scan(&id, &userID, &content, &title, &interest, &photo, &createdAt); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		var nickname string
		err = databases.DB.QueryRow(`SELECT nickname FROM users WHERE id = ?`, userID).Scan(&nickname)
		if err != nil {
			log.Println("Nickname not found for user_id:", userID)
			nickname = "Unknown"
		}

		post := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"content":    content,
			"title":      title,
			"interest":   interest,
			"photo":      nil,
			"created_at": createdAt,
			"nickname":   nickname,
		}

		if photo.Valid {
			post["photo"] = photo.String
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
