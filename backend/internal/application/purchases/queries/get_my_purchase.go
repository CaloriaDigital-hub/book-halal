package queries

import (
	"context"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

// --- Interface ---

type GetMyPurchaseHandler interface {
	// Handle возвращает последнюю покупку книги пользователем или nil.
	Handle(ctx context.Context, userID, bookID string) (*app.PurchaseView, error)
}

// --- Handler ---

type getMyPurchaseHandler struct {
	repo domain.Repository
}

func NewGetMyPurchaseHandler(repo domain.Repository) GetMyPurchaseHandler {
	return &getMyPurchaseHandler{repo: repo}
}

func (h *getMyPurchaseHandler) Handle(ctx context.Context, userID, bookID string) (*app.PurchaseView, error) {
	p, err := h.repo.FindLatest(ctx, userID, bookID)
	if err != nil || p == nil {
		return nil, err
	}
	view := app.ToView(p)
	return &view, nil
}
