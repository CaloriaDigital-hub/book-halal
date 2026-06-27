package books

import (
	"encoding/json"
	"net/http"
)

type askRequest struct {
	Question string `json:"question"`
}

func (c *BookController) AskBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	var req askRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
		c.sendError(w, http.StatusBadRequest, "question is required")
		return
	}

	result, err := c.ragClient.AskBook(bookID, req.Question)
	if err != nil {
		c.sendError(w, http.StatusServiceUnavailable, "rag service error: "+err.Error())
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}
