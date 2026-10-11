package purchases

import (
	"github.com/jackc/pgx/v5"

	domain "book_halal/internal/domain/purchases"
)

func scanPurchase(row pgx.Row, extra ...any) (*domain.Purchase, error) {
	var p domain.Purchase
	var status, contactType string
	dest := append([]any{
		&p.ID, &p.UserID, &p.BookID, &p.Amount, &status, &contactType, &p.Contact.Value,
		&p.ReceiptPath, &p.ReceiptMime, &p.TelegramMessageID, &p.ReviewedBy,
		&p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
	}, extra...)
	if err := row.Scan(dest...); err != nil {
		return nil, err
	}
	p.Status = domain.Status(status)
	p.Contact.Type = domain.ContactType(contactType)
	return &p, nil
}

func scanDetails(row pgx.Row) (*domain.PurchaseDetails, error) {
	var d domain.PurchaseDetails
	p, err := scanPurchase(row, &d.BookTitle, &d.BuyerFirstName, &d.BuyerLastName, &d.BuyerEmail)
	if err != nil {
		return nil, err
	}
	d.Purchase = *p
	return &d, nil
}
