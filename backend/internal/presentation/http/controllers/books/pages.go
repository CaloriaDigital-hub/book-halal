package books

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	// TODO: ВРЕМЕННОЕ РЕШЕНИЕ - потом убрать и использовать app error handler
	repoBooks "book_halal/internal/domain/books"
)

func (c *BookController) GetPages(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}
	
	fromStr := r.URL.Query().Get("from")
	limitStr := r.URL.Query().Get("limit")

	offset := 0
	limit := 0

	if fromStr != "" {
		from, err := strconv.Atoi(fromStr)
		if err != nil || from < 1 {
			c.sendError(w, http.StatusBadRequest, "invalid 'from' parameter")
			return
		}
		offset = from - 1
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 0 {
			c.sendError(w, http.StatusBadRequest, "invalid 'limit' parameter")
			return
		}
		limit = l
	}

	result, err := c.pagesHandler.Handle(r.Context(), bookID, offset, limit)
	if err != nil {
		if errors.Is(err, repoBooks.ErrBookNotFound) {
			c.sendError(w, http.StatusNotFound, "book not found")
			return
		}

		log.Printf("[ERROR] pagesHandler.Handle failed for book %s: %v", bookID, err)
		c.sendError(w, http.StatusInternalServerError, "failed to get book pages")
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}