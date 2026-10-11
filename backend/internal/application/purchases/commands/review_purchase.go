package commands

import (
	"context"
	"log/slog"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

// --- DTO ---

type ReviewPurchaseCommand struct {
	PurchaseID string
	Approve    bool
	// ReviewedBy — кто принял решение: "@username" из Telegram или email админа.
	ReviewedBy string
}

// --- Interface ---

type ReviewPurchaseHandler interface {
	Handle(ctx context.Context, cmd ReviewPurchaseCommand) (*app.AdminPurchaseView, error)
}

// --- Handler ---

type reviewPurchaseHandler struct {
	repo     domain.Repository
	notifier app.Notifier
	logger   *slog.Logger
}

func NewReviewPurchaseHandler(repo domain.Repository, notifier app.Notifier, logger *slog.Logger) ReviewPurchaseHandler {
	return &reviewPurchaseHandler{repo: repo, notifier: notifier, logger: logger}
}

func (h *reviewPurchaseHandler) Handle(ctx context.Context, cmd ReviewPurchaseCommand) (*app.AdminPurchaseView, error) {
	status := domain.StatusRejected
	if cmd.Approve {
		status = domain.StatusApproved
	}
	if _, err := h.repo.Review(ctx, cmd.PurchaseID, status, cmd.ReviewedBy); err != nil {
		return nil, err
	}

	details, err := h.repo.GetDetails(ctx, cmd.PurchaseID)
	if err != nil {
		return nil, err
	}
	// Решение уже сохранено — сбой уведомления не должен его откатывать.
	if err := h.notifier.PurchaseReviewed(ctx, app.NoticeFromDetails(details)); err != nil {
		h.logger.Error("notify about reviewed purchase", "purchase_id", cmd.PurchaseID, "err", err)
	}

	view := app.ToAdminView(details)
	return &view, nil
}
