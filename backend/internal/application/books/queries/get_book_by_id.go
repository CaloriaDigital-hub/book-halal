package queries

import (
	"context"
	"fmt"

	repoBooks "book_halal/internal/domain/books"
)

// --- DTO ---

// BookDetailsView — DTO для детальной страницы книги
type BookDetailsView struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CoverURL    string `json:"cover_url"`
	TotalPages  int    `json:"total_pages"` // Добавили, чтобы фронт знал, сколько всего страниц
	Status      string `json:"status"`      // На всякий случай (ready, processing)
}

// --- Interface ---

type GetBookByIDHandler interface {
	Handle(ctx context.Context, bookID string) (BookDetailsView, error)
}

// --- Handler ---

type getBookByIDHandler struct {
	repo repoBooks.Repository
}

func NewGetBookByIDHandler(repo repoBooks.Repository) GetBookByIDHandler {
	return &getBookByIDHandler{repo: repo}
}

func (h *getBookByIDHandler) Handle(ctx context.Context, id string) (BookDetailsView, error) {
	b, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return BookDetailsView{}, fmt.Errorf("failed to fetch book: %w", err)
	}

	return BookDetailsView{
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
