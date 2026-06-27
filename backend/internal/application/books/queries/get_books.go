package queries

import (
	"context"
	"fmt"

	repoBooks "book_halal/internal/domain/books"
)

// --- DTO ---

type BookCatalogView struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CoverURL    string `json:"cover_url"`
}

// --- Interface ---

type GetBooksHandler interface {
	Handle(ctx context.Context) ([]BookCatalogView, error)
}

// --- Handler ---

type getBooksHandler struct {
	repo repoBooks.Repository
}

func NewGetBooksHandler(repo repoBooks.Repository) GetBooksHandler {
	return &getBooksHandler{repo: repo}
}

func (h *getBooksHandler) Handle(ctx context.Context) ([]BookCatalogView, error) {
	domainBooks, err := h.repo.GetAllReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books from db: %w", err)
	}

	result := make([]BookCatalogView, 0, len(domainBooks))

	for _, b := range domainBooks {
		result = append(result, BookCatalogView{
			ID:          b.ID.String(),
			Title:       b.Title,
			Author:      b.Author,
			Description: b.Description,
			Price:       b.Price.MajorValue(),
			CoverURL:    b.CoverURL,
		})
	}

	return result, nil
}
