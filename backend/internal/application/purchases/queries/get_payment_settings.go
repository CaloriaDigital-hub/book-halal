package queries

import (
	"context"
	"time"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

// --- DTO ---

type PaymentSettingsView struct {
	Details   string    `json:"details"`
	QRURL     string    `json:"qr_url"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Interface ---

type GetPaymentSettingsHandler interface {
	Handle(ctx context.Context) (*PaymentSettingsView, error)
}

// --- Handler ---

type getPaymentSettingsHandler struct {
	repo domain.Repository
}

func NewGetPaymentSettingsHandler(repo domain.Repository) GetPaymentSettingsHandler {
	return &getPaymentSettingsHandler{repo: repo}
}

func (h *getPaymentSettingsHandler) Handle(ctx context.Context) (*PaymentSettingsView, error) {
	s, err := h.repo.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return &PaymentSettingsView{Details: s.Details, QRURL: s.QRURL, Currency: app.Currency, UpdatedAt: s.UpdatedAt}, nil
}
