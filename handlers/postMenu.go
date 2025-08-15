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
	mu.Lock()
	loggedIn, _ := IsLoggedIn(r)
	if !loggedIn {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	mu.Unlock()

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		errorHandler(http.StatusBadRequest, w)
		return
	}
	query := "DELETE FROM posts WHERE id = ?"
	_, err = databases.DB.Exec(query, post.ID)
	if err != nil {
		errorHandler(http.StatusInternalServerError, w)
		return
	}

	query1 := "SELECT categoryID from categories_post WHERE postID = ?"
	rows, err1 := databases.DB.Query(query1, post.ID)
	if err1 != nil {
		errorHandler(http.StatusInternalServerError, w)
		return
	}
	categoryIDs := []int{}
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			errorHandler(http.StatusInternalServerError, w)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	for _, catecategoryID := range categoryIDs {
		query3 := "DELETE FROM categories_post WHERE postID = ? AND categoryID = ?"
		_, err = databases.DB.Exec(query3, post.ID, catecategoryID)
		if err != nil {
			errorHandler(http.StatusInternalServerError, w)
			return
		}
	}

	// errorHandler(http.StatusOK, w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))

}

func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(http.StatusMethodNotAllowed, w)
		return
	}

	var data struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		errorHandler(http.StatusBadRequest, w)
		return
	}

	_, err := databases.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", data.Title, data.Content, data.ID)
	if err != nil {
		errorHandler(http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}
