package queries

import (
	"context"
	"fmt"

	repoBooks "book_halal/internal/domain/books"
)

// --- DTO ---

// PageView — DTO для одной страницы
type PageView struct {
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}

// BookPagesResponse — Обертка для ответа
type BookPagesResponse struct {
	TotalPages int        `json:"total_pages"`
	Pages      []PageView `json:"pages"`
}

// --- Interface ---

type GetBookPagesHandler interface {
	Handle(ctx context.Context, bookID string, offset, limit int) (BookPagesResponse, error)
}

// --- Handler ---

type getBookPagesHandler struct {
	repo repoBooks.Repository
}

func NewGetBookPagesHandler(repo repoBooks.Repository) GetBookPagesHandler {
	return &getBookPagesHandler{repo: repo}
}

func (h *getBookPagesHandler) Handle(ctx context.Context, bookID string, offset, limit int) (BookPagesResponse, error) {
	// limit == 0 → без пагинации (все страницы) — обратная совместимость
	if limit == 0 {
		pages, totalPages, err := h.repo.GetPagesByBookID(ctx, bookID)
		if err != nil {
			return BookPagesResponse{}, fmt.Errorf("failed to get pages from db: %w", err)
		}
		return BookPagesResponse{
			TotalPages: totalPages,
			Pages:      toPageViews(pages),
		}, nil
	}

	pages, totalPages, err := h.repo.GetPagesByBookIDPaginated(ctx, bookID, offset, limit)
	if err != nil {
		return BookPagesResponse{}, fmt.Errorf("failed to get pages from db: %w", err)
	}

	return BookPagesResponse{
		TotalPages: totalPages,
		Pages:      toPageViews(pages),
	}, nil
}
