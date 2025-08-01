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
