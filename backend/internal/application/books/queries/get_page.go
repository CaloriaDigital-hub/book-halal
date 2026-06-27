package queries

import (
	"context"
	"fmt"

	repoBooks "book_halal/internal/domain/books"
)

// --- DTO ---

type SinglePageResponse struct {
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}

// --- Interface ---

type GetPageHandler interface {
	Handle(ctx context.Context, bookID string, pageNumber int) (*SinglePageResponse, error)
}

// --- Handler ---

type getPageHandler struct {
	repo repoBooks.Repository
}

func NewGetPageHandler(repo repoBooks.Repository) GetPageHandler {
	return &getPageHandler{repo: repo}
}

func (h *getPageHandler) Handle(ctx context.Context, bookID string, pageNumber int) (*SinglePageResponse, error) {
	if pageNumber < 1 {
		return nil, fmt.Errorf("page number must be at least 1")
	}

	page, err := h.repo.GetPageByNumber(ctx, bookID, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	if page == nil {
		return nil, fmt.Errorf("page %d not found", pageNumber)
	}

	return &SinglePageResponse{
		PageNumber: page.PageNumber,
		ImageURL:   page.ImageURL,
	}, nil
}
