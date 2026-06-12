package books

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBookByID обрабатывает GET /api/books/{id}
func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Достаем ID из URL (фича Go 1.22)
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	bookDetails, err := c.detailsHandler .Handle(r.Context(), bookID)
	if err != nil {
		// Если это наша ошибка "book not found", отдаем 404
		if err.Error() == "failed to fetch book: book not found" {
			c.sendError(w, http.StatusNotFound, "book not found")
			return
		}

		fmt.Println("ПИЗДЕЦ В ID GET:", err)
		c.sendError(w, http.StatusInternalServerError, "failed to get book details")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookDetails)
}