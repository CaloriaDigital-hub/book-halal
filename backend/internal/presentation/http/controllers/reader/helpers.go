package reader

import (
	"encoding/json"
	"net/http"
)

func (c *ReaderController) sendError(w http.ResponseWriter, status int, msg string) {
	c.sendJSON(w, status, map[string]string{"error": msg})
}

func (c *ReaderController) sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
