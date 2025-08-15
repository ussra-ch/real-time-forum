package handlers

import (
	"encoding/json"
	"net/http"
)

func errorHandler(errorType int, w http.ResponseWriter) {
	errorr := ErrorStruct{
		Type: "error",
		Text: http.StatusText(errorType),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorType)
	json.NewEncoder(w).Encode(errorr)
}
