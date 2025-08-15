package handlers

import (
	"encoding/json"
	"net/http"

	"handlers/databases"
)

type Post struct {
	ID int `json:"id"`
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	loggedIn, _ := IsLoggedIn(r)
	if !loggedIn {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	query := "DELETE FROM posts WHERE id = ?"
	_, err = databases.DB.Exec(query, post.ID)
	if err != nil {
		http.Error(w, "error deleting posts from post", http.StatusInternalServerError)
		return
	}

	query1 := "SELECT categoryID from categories_post WHERE postID = ?"
	rows, err1 := databases.DB.Query(query1, post.ID)

	if err1 != nil {
		http.Error(w, "error getting posts from categories_post", http.StatusInternalServerError)
		return
	}
	categoryIDs := []int{}
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			http.Error(w, "error scanning categoryID", http.StatusInternalServerError)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	for _, catecategoryID := range categoryIDs {
		query3 := "DELETE FROM categories_post WHERE postID = ? AND categoryID = ?"
		_, err = databases.DB.Exec(query3, post.ID, catecategoryID)
		if err != nil {
			http.Error(w, "error deleting posts from categories_post", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	loggedIn, _ := IsLoggedIn(r)
	if !loggedIn {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err := databases.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", data.Title, data.Content, data.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}
