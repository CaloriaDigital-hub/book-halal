package purchases

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	pool *pgxpool.Pool
}

func NewPurchaseRepository(pool *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{pool: pool}
}

const purchaseColumns = `
	p.id, p.user_id, p.book_id, p.amount, p.status, p.contact_type, p.contact_value,
	p.receipt_path, p.receipt_mime, COALESCE(p.telegram_message_id, 0), COALESCE(p.reviewed_by, ''),
	p.reviewed_at, p.created_at, p.updated_at`
