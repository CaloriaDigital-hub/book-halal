package purchases

import (
	"time"

	domain "book_halal/internal/domain/purchases"
)

// PurchaseView — покупка глазами покупателя.
type PurchaseView struct {
	ID         string     `json:"id"`
	BookID     string     `json:"book_id"`
	Amount     int        `json:"amount"`
	Currency   string     `json:"currency"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}

func ToView(p *domain.Purchase) PurchaseView {
	return PurchaseView{
		ID:         p.ID,
		BookID:     p.BookID,
		Amount:     p.Amount / 100,
		Currency:   Currency,
		Status:     string(p.Status),
		CreatedAt:  p.CreatedAt,
		ReviewedAt: p.ReviewedAt,
	}
}

// AdminPurchaseView — покупка в списке админки.
type AdminPurchaseView struct {
	PurchaseView
	BookTitle    string `json:"book_title"`
	BuyerName    string `json:"buyer_name"`
	BuyerEmail   string `json:"buyer_email"`
	ContactType  string `json:"contact_type"`
	ContactValue string `json:"contact_value"`
	ReceiptMime  string `json:"receipt_mime"`
	ReviewedBy   string `json:"reviewed_by,omitempty"`
}

func ToAdminView(d *domain.PurchaseDetails) AdminPurchaseView {
	return AdminPurchaseView{
		PurchaseView: ToView(&d.Purchase),
		BookTitle:    d.BookTitle,
		BuyerName:    d.BuyerFirstName + " " + d.BuyerLastName,
		BuyerEmail:   d.BuyerEmail,
		ContactType:  string(d.Contact.Type),
		ContactValue: d.Contact.Value,
		ReceiptMime:  d.ReceiptMime,
		ReviewedBy:   d.ReviewedBy,
	}
}
