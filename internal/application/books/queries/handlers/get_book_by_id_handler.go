package handlers

import (
	"context"
	"fmt"

	"book_halal/internal/application/books/queries"
	repoBooks "book_halal/internal/domain/books"
)

type GetBookByIDQueryHandler struct {
	repo repoBooks.Repository
}

func NewGetBookByIDQueryHandler(repo repoBooks.Repository) *GetBookByIDQueryHandler {
	return &GetBookByIDQueryHandler{repo: repo}
}

func (h *GetBookByIDQueryHandler) Handle(ctx context.Context, id string) (queries.BookDetailsView, error) {
	b, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return queries.BookDetailsView{}, fmt.Errorf("failed to fetch book: %w", err)
	}


	return queries.BookDetailsView{
		ID:          b.ID.String(), 
		Title:       b.Title,
		Author:      b.Author,
		Description: b.Description,
		Price:       b.Price.MajorValue(),
		CoverURL:    b.CoverURL,
		TotalPages:  b.TotalPages,
		Status:      string(b.Status),
	}, nil
}