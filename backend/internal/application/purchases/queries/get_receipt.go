package queries

import (
	"context"
	"io"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

// --- Interface ---

type GetReceiptHandler interface {
	// Handle открывает файл чека; вызывающий закрывает reader.
	Handle(ctx context.Context, purchaseID string) (io.ReadCloser, string, error)
}

// --- Handler ---

type getReceiptHandler struct {
	repo     domain.Repository
	receipts app.FileStorage
}

func NewGetReceiptHandler(repo domain.Repository, receipts app.FileStorage) GetReceiptHandler {
	return &getReceiptHandler{repo: repo, receipts: receipts}
}

func (h *getReceiptHandler) Handle(ctx context.Context, purchaseID string) (io.ReadCloser, string, error) {
	d, err := h.repo.GetDetails(ctx, purchaseID)
	if err != nil {
		return nil, "", err
	}
	f, err := h.receipts.Open(d.ReceiptPath)
	if err != nil {
		return nil, "", err
	}
	return f, d.ReceiptMime, nil
}
