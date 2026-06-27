package commands

import (
	"context"
	"fmt"

	"book_halal/internal/domain/reader"
)

// --- DTO ---

type RemoveBookmarkCommand struct {
	UserID     string
	BookID     string
	PageNumber int
}

// --- Interface ---

type RemoveBookmarkHandler interface {
	Handle(ctx context.Context, cmd RemoveBookmarkCommand) error
}

// --- Handler ---

type removeBookmarkHandler struct {
	repo reader.Repository
}

func NewRemoveBookmarkHandler(repo reader.Repository) RemoveBookmarkHandler {
	return &removeBookmarkHandler{repo: repo}
}

func (h *removeBookmarkHandler) Handle(ctx context.Context, cmd RemoveBookmarkCommand) error {
	if cmd.PageNumber < 1 {
		return fmt.Errorf("page number must be at least 1")
	}

	return h.repo.DeleteBookmark(ctx, cmd.UserID, cmd.BookID, cmd.PageNumber)
}
