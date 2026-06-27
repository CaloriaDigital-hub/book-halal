package commands

import (
	"context"
	"fmt"
	"time"

	"book_halal/internal/domain/reader"
)

// --- DTO ---

type AddBookmarkCommand struct {
	UserID     string
	BookID     string
	PageNumber int
}

type BookmarkView struct {
	ID         string    `json:"id"`
	PageNumber int       `json:"page_number"`
	CreatedAt  time.Time `json:"created_at"`
}

// --- Interface ---

type AddBookmarkHandler interface {
	Handle(ctx context.Context, cmd AddBookmarkCommand) (*BookmarkView, error)
}

// --- Handler ---

type addBookmarkHandler struct {
	repo reader.Repository
}

func NewAddBookmarkHandler(repo reader.Repository) AddBookmarkHandler {
	return &addBookmarkHandler{repo: repo}
}

func (h *addBookmarkHandler) Handle(ctx context.Context, cmd AddBookmarkCommand) (*BookmarkView, error) {
	if cmd.PageNumber < 1 {
		return nil, fmt.Errorf("page number must be at least 1")
	}

	bookmark, err := h.repo.AddBookmark(ctx, cmd.UserID, cmd.BookID, cmd.PageNumber)
	if err != nil {
		return nil, err
	}

	return &BookmarkView{
		ID:         bookmark.ID,
		PageNumber: bookmark.PageNumber,
		CreatedAt:  bookmark.CreatedAt,
	}, nil
}
