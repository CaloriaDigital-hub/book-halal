package queries

import (
	"context"
	"time"

	"book_halal/internal/domain/reader"
)

// --- DTO ---

type ProgressResponse struct {
	PageNumber int       `json:"page_number"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// --- Interface ---

type GetProgressHandler interface {
	Handle(ctx context.Context, userID, bookID string) (*ProgressResponse, error)
}

// --- Handler ---

type getProgressHandler struct {
	repo reader.Repository
}

func NewGetProgressHandler(repo reader.Repository) GetProgressHandler {
	return &getProgressHandler{repo: repo}
}

func (h *getProgressHandler) Handle(ctx context.Context, userID, bookID string) (*ProgressResponse, error) {
	progress, err := h.repo.GetProgress(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	if progress == nil {
		// Ещё не начинал читать — страница 1
		return &ProgressResponse{PageNumber: 1}, nil
	}

	return &ProgressResponse{
		PageNumber: progress.PageNumber,
		UpdatedAt:  progress.UpdatedAt,
	}, nil
}
