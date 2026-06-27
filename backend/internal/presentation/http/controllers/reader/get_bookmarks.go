package reader

import (
	"log"
	"net/http"

	"book_halal/internal/middleware"
)

func (c *ReaderController) GetBookmarks(w http.ResponseWriter, r *http.Request) {
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

	result, err := c.getBookmarks.Handle(r.Context(), user.ID.String(), bookID)
	if err != nil {
		log.Printf("[ERROR] GetBookmarks: %v", err)
		c.sendError(w, http.StatusInternalServerError, "failed to get bookmarks")
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}
