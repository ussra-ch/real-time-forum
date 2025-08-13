package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
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
		errorr := ErrorStruct{
			Type: "error",
			Text: "Method not allowed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorr)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	topics := r.Form["topics"]
	if title == "" || description == "" || len(topics) == 0 {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}

	categoriesRows, err := databases.DB.Query("SELECT name FROM categories")
	defer categoriesRows.Close()

	var allCategories []string
	for categoriesRows.Next() {
		var name string
		if err := categoriesRows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		allCategories = append(allCategories, name)
	}

	found := true
	for _, topic := range topics {
		if !contains(allCategories, topic) {
			found = false
			break
		}
	}

	if !found {
		errorr := ErrorStruct{
			Type: "error",
			Text: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
		return
	}

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
			// fmt.Println(err)
			errorr := ErrorStruct{
				Type: "error",
				Text: "Internal server error",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorr)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Internal server error",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorr)
			return
		}
	} else {
		filename = ""
	}

	query := `
		INSERT INTO posts (title, content, interest, user_id, photo)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := databases.DB.Exec(query, html.EscapeString(title), html.EscapeString(description), strings.Join(topics, ","), userID, filename)
	if err != nil {
		// log.Println("Insert post error:", err)
		errorr := ErrorStruct{
			Type: "error",
			Text: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorr)
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		// log.Println("Error getting inserted post ID:", err)
		errorr := ErrorStruct{
			Type: "error",
			Text: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorr)
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
		errorr := ErrorStruct{
			Type: "error",
			Text: "Method Not Allowed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorr)
		return
	}
	_, UserID := IsLoggedIn(r)

	query := `SELECT id, user_id, content, title, interest, photo, created_at FROM posts`
	rows, err := databases.DB.Query(query)
	if err != nil {
		// log.Println("Error fetching posts:", err)
		errorr := ErrorStruct{
			Type: "error",
			Text: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorr)
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
			"myId":       UserID,
		}

		if photo.Valid {
			post["photo"] = photo.String
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
