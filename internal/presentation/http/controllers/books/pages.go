package books

import "net/http"

func (c *BookController) GetPages(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	result, err := c.pagesHandler.Handle(r.Context(), bookID)
	if err != nil {
		c.sendError(w, http.StatusInternalServerError, "failed to get book pages")
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}