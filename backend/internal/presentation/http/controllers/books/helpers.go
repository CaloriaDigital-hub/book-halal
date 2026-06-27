package books

import (
	"encoding/json"
	"net/http"
)

func (c *BookController) sendError(w http.ResponseWriter, status int, msg string) {
	c.sendJSON(w, status, map[string]string{"error": msg})
}

func (c *BookController) sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}