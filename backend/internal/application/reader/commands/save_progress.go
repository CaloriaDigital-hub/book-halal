package commands

import (
	"context"
	"fmt"

	"book_halal/internal/domain/reader"
)

// --- DTO ---

type SaveProgressCommand struct {
	UserID     string
	BookID     string
	PageNumber int
}

// --- Interface ---

type SaveProgressHandler interface {
	Handle(ctx context.Context, cmd SaveProgressCommand) error
}

// --- Handler ---

type saveProgressHandler struct {
	repo reader.Repository
}

func NewSaveProgressHandler(repo reader.Repository) SaveProgressHandler {
	return &saveProgressHandler{repo: repo}
}

func (h *saveProgressHandler) Handle(ctx context.Context, cmd SaveProgressCommand) error {
	if cmd.PageNumber < 1 {
		return fmt.Errorf("page number must be at least 1")
	}

	return h.repo.UpsertProgress(ctx, cmd.UserID, cmd.BookID, cmd.PageNumber)
}
