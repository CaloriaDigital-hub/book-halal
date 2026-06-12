package books

import (
	"errors"
	"log"
	"net/http"

	repoBooks "book_halal/internal/domain/books" // Импортируем твой домен, где должна лежать ошибка
)

func (c *BookController) GetPages(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		c.sendError(w, http.StatusBadRequest, "book id is required")
		return
	}

	result, err := c.pagesHandler.Handle(r.Context(), bookID)
	if err != nil {
		// Проверяем, является ли ошибка именно "книга не найдена"
		if errors.Is(err, repoBooks.ErrBookNotFound) {
			c.sendError(w, http.StatusNotFound, "book not found")
			return
		}

		// Технические ошибки (упала БД и т.д.) логируем, чтобы видеть в консоли/кибане
		log.Printf("[ERROR] pagesHandler.Handle failed for book %s: %v", bookID, err)
		
		// А пользователю отдаем общую 500
		c.sendError(w, http.StatusInternalServerError, "failed to get book pages")
		return
	}

	c.sendJSON(w, http.StatusOK, result)
}