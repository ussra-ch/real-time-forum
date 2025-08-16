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
type Category struct {
	Id   int
	Name string
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(http.StatusMethodNotAllowed, w)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorHandler(http.StatusBadRequest, w)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	topics := r.Form["topics"]
	if title == "" || description == "" || len(topics) == 0 {
		errorHandler(http.StatusBadRequest, w)
		return
	}

	categoriesRows, err := databases.DB.Query("SELECT * FROM categories")
	defer categoriesRows.Close()

	if err != nil {
		errorHandler(http.StatusInternalServerError, w)
		return
	}

	var allCategories []Category
	for categoriesRows.Next() {
		var name string
		var categoryId int
		if err := categoriesRows.Scan(&categoryId, &name); err != nil {
			log.Fatal(err)
		}
		allCategories = append(allCategories, Category{Id: categoryId, Name: name})
	}

	var updatedTopics []Category
	found := true
	for _, topic := range topics {
		ok, id := contains(allCategories, topic)
		if !ok {
			found = false
			break
		} else {
			updatedTopics = append(updatedTopics, Category{Id: id, Name: topic})
		}
	}

	if !found {
		errorHandler(http.StatusBadRequest, w)
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
			errorHandler(http.StatusInternalServerError, w)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			errorHandler(http.StatusInternalServerError, w)
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
		errorHandler(http.StatusInternalServerError, w)
		return
	}
	mu.Lock()
	postID, err := res.LastInsertId()
	if err != nil {
		errorHandler(http.StatusInternalServerError, w)
		return
	}
	for _, x := range updatedTopics {
		query2 := `INSERT INTO categories_post (categoryID, postID) VALUES (?, ?)`
		_, err := databases.DB.Exec(query2, x.Id, postID)
		if err != nil {
			errorHandler(http.StatusInternalServerError, w)
			return
		}
	}
	mu.Unlock()
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
		errorHandler(http.StatusMethodNotAllowed, w)
		return
	}
	mu.Lock()
	_, UserID := IsLoggedIn(r)
	mu.Unlock()

	query := `SELECT id, user_id, content, title, interest, photo, created_at FROM posts`
	rows, err := databases.DB.Query(query)
	if err != nil {
		errorHandler(http.StatusInternalServerError, w)
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

func contains(slice []Category, item string) (bool, int) {
	for _, v := range slice {
		if v.Name == item {
			return true, v.Id
		}
	}
	return false, -1
}
