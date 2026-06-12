package handlers

import (
	"context"
	"fmt"

	"book_halal/internal/application/books/queries"
	repoBooks "book_halal/internal/domain/books"
)

type GetBooksQueryHandler struct {
	repo repoBooks.Repository
}

func NewGetBooksQueryHandler(repo repoBooks.Repository) *GetBooksQueryHandler {
	return &GetBooksQueryHandler{repo: repo}
}

func (h *GetBooksQueryHandler) Handle(ctx context.Context) ([]queries.BookCatalogView, error) {
	domainBooks, err := h.repo.GetAllReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books from db: %w", err)
	}

	result := make([]queries.BookCatalogView, 0, len(domainBooks))
	
	for _, b := range domainBooks {
		result = append(result, queries.BookCatalogView{
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