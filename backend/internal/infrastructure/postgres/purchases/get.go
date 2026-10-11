package purchases

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	domain "book_halal/internal/domain/purchases"
)

const detailsQuery = `
	SELECT ` + purchaseColumns + `, b.title, u.first_name, u.last_name, u.email
	FROM purchases p
	JOIN books b ON b.id = p.book_id
	JOIN users u ON u.id = p.user_id`

func (r *PurchaseRepository) GetDetails(ctx context.Context, purchaseID string) (*domain.PurchaseDetails, error) {
	d, err := scanDetails(r.pool.QueryRow(ctx, detailsQuery+` WHERE p.id = $1`, purchaseID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPurchaseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get purchase: %w", err)
	}
	return d, nil
}

func (r *PurchaseRepository) FindLatest(ctx context.Context, userID, bookID string) (*domain.Purchase, error) {
	query := `SELECT ` + purchaseColumns + ` FROM purchases p
		WHERE p.user_id = $1 AND p.book_id = $2
		ORDER BY p.created_at DESC LIMIT 1`
	p, err := scanPurchase(r.pool.QueryRow(ctx, query, userID, bookID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find latest purchase: %w", err)
	}
	return p, nil
}

func (r *PurchaseRepository) HasApproved(ctx context.Context, userID, bookID string) (bool, error) {
	var ok bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM purchases WHERE user_id = $1 AND book_id = $2 AND status = 'approved')`,
		userID, bookID,
	).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check approved purchase: %w", err)
	}
	return ok, nil
}

func (r *PurchaseRepository) List(ctx context.Context, status domain.Status, limit int) ([]domain.PurchaseDetails, error) {
	rows, err := r.pool.Query(ctx,
		detailsQuery+` WHERE ($1 = '' OR p.status = $1) ORDER BY p.created_at DESC LIMIT $2`,
		string(status), limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list purchases: %w", err)
	}
	defer rows.Close()

	var result []domain.PurchaseDetails
	for rows.Next() {
		d, err := scanDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("scan purchase: %w", err)
		}
		result = append(result, *d)
	}
	return result, rows.Err()
}
