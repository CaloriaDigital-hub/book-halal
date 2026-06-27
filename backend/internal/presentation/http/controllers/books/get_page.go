package books

import (
	"log"
	"net/http"
	"strconv"
)

func (c *BookController) GetPage(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	pageStr := r.PathValue("page")
	pageNumber, err := strconv.Atoi(pageStr)
	if err != nil || pageNumber < 1 {
		c.sendError(w, http.StatusBadRequest, "invalid page number")
		return
	}

	result, err := c.pageHandler.Handle(r.Context(), bookID, pageNumber)
	if err != nil {
		log.Printf("[ERROR] GetPage: book=%s page=%d err=%v", bookID, pageNumber, err)
		c.sendError(w, http.StatusNotFound, "page not found")
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}
