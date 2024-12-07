package writer

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Trace   string `json:"trace,omitempty"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
