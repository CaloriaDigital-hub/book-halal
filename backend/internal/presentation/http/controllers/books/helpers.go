package books

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	repoBooks "book_halal/internal/domain/books"
	"book_halal/internal/middleware"
)

func (c *BookController) sendError(w http.ResponseWriter, status int, msg string) {
	c.sendJSON(w, status, map[string]string{"error": msg})
}

func (c *BookController) sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
// ensureAccess пускает к содержимому книги: админа, всех — к бесплатной, к платной — после подтверждённой оплаты.
// Сам пишет ответ с ошибкой и возвращает false, если доступа нет.
func (c *BookController) ensureAccess(w http.ResponseWriter, r *http.Request, bookID string) bool {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		c.sendError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	allowed, err := c.access.CanRead(r.Context(), user, bookID)
	switch {
	case errors.Is(err, repoBooks.ErrBookNotFound):
		c.sendError(w, http.StatusNotFound, "book not found")
		return false
	case err != nil:
		log.Printf("[ERROR] access check: book=%s err=%v", bookID, err)
		c.sendError(w, http.StatusInternalServerError, "failed to check access")
		return false
	case !allowed:
		c.sendError(w, http.StatusForbidden, "book is not purchased")
		return false
	}
	return true
}
