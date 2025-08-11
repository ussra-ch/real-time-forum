package handlers

import (
	"encoding/json"
	"handlers/databases"
	"net/http"
)

type Post struct {
	ID int `json:"id"`
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var post Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	query := "DELETE FROM posts WHERE id = ?"
	_, err = databases.DB.Exec(query, post.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

func EditPost(w http.ResponseWriter, r *http.Request) {
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
