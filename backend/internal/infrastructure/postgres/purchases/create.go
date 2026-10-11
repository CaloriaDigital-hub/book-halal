package purchases

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	domain "book_halal/internal/domain/purchases"
)

func (r *PurchaseRepository) Create(ctx context.Context, p *domain.Purchase) error {
	query := `
		INSERT INTO purchases (id, user_id, book_id, amount, status, contact_type, contact_value,
		                       receipt_path, receipt_mime, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.ID, p.UserID, p.BookID, p.Amount, p.Status, p.Contact.Type, p.Contact.Value,
		p.ReceiptPath, p.ReceiptMime,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "unique_active_purchase" {
		return domain.ErrActivePurchaseExists
	}
	if err != nil {
		return fmt.Errorf("insert purchase: %w", err)
	}
	return nil
}

func (r *PurchaseRepository) SetTelegramMessage(ctx context.Context, purchaseID string, messageID int64) error {
	_, err := r.pool.Exec(ctx, `UPDATE purchases SET telegram_message_id = $2 WHERE id = $1`, purchaseID, messageID)
	if err != nil {
		return fmt.Errorf("set telegram message: %w", err)
	}
	return nil
}
