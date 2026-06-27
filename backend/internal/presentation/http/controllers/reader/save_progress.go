package reader

import (
	"encoding/json"
	"log"
	"net/http"

	"book_halal/internal/application/reader/commands"
	"book_halal/internal/middleware"
)

func (c *ReaderController) SaveProgress(w http.ResponseWriter, r *http.Request) {
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

	err := c.saveProgress.Handle(r.Context(), commands.SaveProgressCommand{
		UserID:     user.ID.String(),
		BookID:     bookID,
		PageNumber: body.PageNumber,
	})
	if err != nil {
		log.Printf("[ERROR] SaveProgress: %v", err)
		c.sendError(w, http.StatusInternalServerError, "failed to save progress")
		return
	}

	c.sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
