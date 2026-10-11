package books

import (
	"encoding/json"
	"fmt"
	"net/http"

	"book_halal/internal/application/books/queries"
	"book_halal/internal/middleware"
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

	// has_access — можно ли читать книгу и спрашивать по ней RAG (для кнопки «Читать» / «Купить»)
	hasAccess := false
	if user, ok := middleware.UserFromContext(r.Context()); ok {
		if hasAccess, err = c.access.CanRead(r.Context(), user, bookID); err != nil {
			fmt.Println("access check in GetBookByID:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		queries.BookDetailsView
		HasAccess bool `json:"has_access"`
	}{bookDetails, hasAccess})
}