package queries

import (
	"context"
	"fmt"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

const maxListLimit = 200

// --- Interface ---

type ListPurchasesHandler interface {
	Handle(ctx context.Context, status string, limit int) ([]app.AdminPurchaseView, error)
}

// --- Handler ---

type listPurchasesHandler struct {
	repo domain.Repository
}

func NewListPurchasesHandler(repo domain.Repository) ListPurchasesHandler {
	return &listPurchasesHandler{repo: repo}
}

func (h *listPurchasesHandler) Handle(ctx context.Context, status string, limit int) ([]app.AdminPurchaseView, error) {
	s := domain.Status(status)
	if status != "" && !s.IsValid() {
		return nil, fmt.Errorf("invalid status %q", status)
	}
	if limit <= 0 || limit > maxListLimit {
		limit = maxListLimit
	}

	items, err := h.repo.List(ctx, s, limit)
	if err != nil {
		return nil, err
	}
	views := make([]app.AdminPurchaseView, 0, len(items))
	for i := range items {
		views = append(views, app.ToAdminView(&items[i]))
	}
	return views, nil
}
