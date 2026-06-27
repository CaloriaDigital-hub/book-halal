package queries

import (
	"context"
	"time"

	"book_halal/internal/domain/reader"
)

// --- DTO ---

type BookmarkView struct {
	ID         string    `json:"id"`
	PageNumber int       `json:"page_number"`
	CreatedAt  time.Time `json:"created_at"`
}

// --- Interface ---

type GetBookmarksHandler interface {
	Handle(ctx context.Context, userID, bookID string) ([]BookmarkView, error)
}

// --- Handler ---

type getBookmarksHandler struct {
	repo reader.Repository
}

func NewGetBookmarksHandler(repo reader.Repository) GetBookmarksHandler {
	return &getBookmarksHandler{repo: repo}
}

func (h *getBookmarksHandler) Handle(ctx context.Context, userID, bookID string) ([]BookmarkView, error) {
	bookmarks, err := h.repo.GetBookmarks(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	result := make([]BookmarkView, 0, len(bookmarks))
	for _, b := range bookmarks {
		result = append(result, BookmarkView{
			ID:         b.ID,
			PageNumber: b.PageNumber,
			CreatedAt:  b.CreatedAt,
		})
	}

	return result, nil
}
