package reader

import (
	"encoding/json"
	"log"
	"net/http"

	"book_halal/internal/application/reader/commands"
	"book_halal/internal/middleware"
)

func (c *ReaderController) AddBookmark(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		c.sendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	var body struct {
		PageNumber int `json:"page_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.addBookmark.Handle(r.Context(), commands.AddBookmarkCommand{
		UserID:     user.ID.String(),
		BookID:     bookID,
		PageNumber: body.PageNumber,
	})
	if err != nil {
		log.Printf("[ERROR] AddBookmark: %v", err)
		c.sendError(w, http.StatusInternalServerError, "failed to add bookmark")
		return
	}

	c.sendJSON(w, http.StatusCreated, result)
}
