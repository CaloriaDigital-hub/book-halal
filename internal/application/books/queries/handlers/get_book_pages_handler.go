package handlers

import (
	"context"
	"fmt"

	"book_halal/internal/application/books/queries"
	repoBooks "book_halal/internal/domain/books"
)

type GetBookPagesQueryHandler struct {
	repo repoBooks.Repository
}

func NewGetBookPagesQueryHandler(repo repoBooks.Repository) *GetBookPagesQueryHandler {
	return &GetBookPagesQueryHandler{repo: repo}
}

func (h *GetBookPagesQueryHandler) Handle(ctx context.Context, bookID string) (queries.BookPagesResponse, error) {

	pages, totalPages, err := h.repo.GetPagesByBookID(ctx, bookID)
	if err != nil {
		return queries.BookPagesResponse{}, fmt.Errorf("failed to get pages from db: %w", err)
	}

	result := make([]queries.PageView, 0, len(pages))
	for _, p := range pages {
		result = append(result, queries.PageView{
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		})
	}

	return queries.BookPagesResponse{
		TotalPages: totalPages,
		Pages:      result,
	}, nil
}