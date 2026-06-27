package ports

import entity "book_halal/internal/domain/books/entity"

type RAGIndexer interface {
	IndexBook(bookID string, baseURL string, pages []entity.Page) error
}