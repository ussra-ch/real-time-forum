package handlers

import (
	"encoding/json"
	"net/http"
)

type PostData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var postData PostData
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post created successfully"))
}
