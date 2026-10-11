package purchases

import (
	"context"
	"fmt"

	domain "book_halal/internal/domain/purchases"
)

func (r *PurchaseRepository) GetSettings(ctx context.Context) (*domain.PaymentSettings, error) {
	var s domain.PaymentSettings
	err := r.pool.QueryRow(ctx, `SELECT details, qr_url, updated_at FROM payment_settings WHERE id = 1`).
		Scan(&s.Details, &s.QRURL, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get payment settings: %w", err)
	}
	return &s, nil
}

func (r *PurchaseRepository) SaveSettings(ctx context.Context, s domain.PaymentSettings) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO payment_settings (id, details, qr_url, updated_at) VALUES (1, $1, $2, NOW())
		ON CONFLICT (id) DO UPDATE SET details = EXCLUDED.details, qr_url = EXCLUDED.qr_url, updated_at = NOW()`,
		s.Details, s.QRURL,
	)
	if err != nil {
		return fmt.Errorf("save payment settings: %w", err)
	}
	return nil
}
