package purchases

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	domain "book_halal/internal/domain/purchases"
)

func (r *PurchaseRepository) Review(ctx context.Context, purchaseID string, status domain.Status, reviewedBy string) (*domain.Purchase, error) {
	// Условие status = 'pending' делает переход атомарным: двойное нажатие
	// или два админа одновременно — обработает только первый.
	query := `
		UPDATE purchases p
		SET status = $2, reviewed_by = $3, reviewed_at = NOW(), updated_at = NOW()
		WHERE p.id = $1 AND p.status = 'pending'
		RETURNING ` + purchaseColumns

	p, err := scanPurchase(r.pool.QueryRow(ctx, query, purchaseID, status, reviewedBy))
	if err == nil {
		return p, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("review purchase: %w", err)
	}

	var exists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM purchases WHERE id = $1)`, purchaseID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check purchase: %w", err)
	}
	if !exists {
		return nil, domain.ErrPurchaseNotFound
	}
	return nil, domain.ErrAlreadyReviewed
}
