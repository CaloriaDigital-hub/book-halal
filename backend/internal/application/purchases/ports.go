package purchases

import (
	"context"
	"io"
	"time"

	domain "book_halal/internal/domain/purchases"
)

// Currency — валюта цен книг (books.price хранится в минорных единицах — тыйынах).
const Currency = "KGS"

// FileStorage — хранилище файлов (чеки — приватное, QR — публичное static).
type FileStorage interface {
	Save(relPath string, r io.Reader) error
	Open(relPath string) (io.ReadCloser, error)
	Remove(relPath string) error
}

// Notice — всё, что нужно уведомлению о покупке (Telegram и т.п.).
type Notice struct {
	PurchaseID        string
	BookTitle         string
	BuyerName         string
	BuyerEmail        string
	Amount            int // в сомах
	Contact           domain.Contact
	Status            domain.Status
	ReviewedBy        string
	ReviewedAt        time.Time
	CreatedAt         time.Time
	TelegramMessageID int64

	// Только для PurchaseCreated
	Receipt     io.Reader
	ReceiptMime string
}

// Notifier сообщает админам о новых чеках и итогах проверки.
type Notifier interface {
	// PurchaseCreated отправляет чек на проверку и возвращает ID сообщения (0, если отправки не было).
	PurchaseCreated(ctx context.Context, n Notice) (messageID int64, err error)
	PurchaseReviewed(ctx context.Context, n Notice) error
}

// NopNotifier — когда Telegram не настроен: покупки работают, подтверждение — из админки.
type NopNotifier struct{}

func (NopNotifier) PurchaseCreated(context.Context, Notice) (int64, error) { return 0, nil }
func (NopNotifier) PurchaseReviewed(context.Context, Notice) error         { return nil }

// NoticeFromDetails собирает уведомление из покупки с данными книги и покупателя.
func NoticeFromDetails(d *domain.PurchaseDetails) Notice {
	n := Notice{
		PurchaseID:        d.ID,
		BookTitle:         d.BookTitle,
		BuyerName:         d.BuyerFirstName + " " + d.BuyerLastName,
		BuyerEmail:        d.BuyerEmail,
		Amount:            d.Amount / 100,
		Contact:           d.Contact,
		Status:            d.Status,
		ReviewedBy:        d.ReviewedBy,
		CreatedAt:         d.CreatedAt,
		TelegramMessageID: d.TelegramMessageID,
	}
	if d.ReviewedAt != nil {
		n.ReviewedAt = *d.ReviewedAt
	}
	return n
}
