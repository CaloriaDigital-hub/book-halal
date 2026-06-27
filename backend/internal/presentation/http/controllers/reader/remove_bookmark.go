package reader

import (
	"log"
	"net/http"
	"strconv"

	"book_halal/internal/application/reader/commands"
	"book_halal/internal/middleware"
)

func (c *ReaderController) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
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

	pageStr := r.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.sendError(w, http.StatusBadRequest, "invalid page number")
		return
	}

	err = c.removeBookmark.Handle(r.Context(), commands.RemoveBookmarkCommand{
		UserID:     user.ID.String(),
		BookID:     bookID,
		PageNumber: page,
	})
	if err != nil {
		log.Printf("[ERROR] RemoveBookmark: %v", err)
		c.sendError(w, http.StatusNotFound, "bookmark not found")
		return
	}

	c.sendJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
