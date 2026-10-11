package purchases

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

func (s Status) IsValid() bool {
	return s == StatusPending || s == StatusApproved || s == StatusRejected
}

type ContactType string

const (
	ContactPhone    ContactType = "phone"
	ContactWhatsApp ContactType = "whatsapp"
	ContactTelegram ContactType = "telegram"
)

var (
	ErrInvalidContact       = errors.New("invalid contact")
	ErrPurchaseNotFound     = errors.New("purchase not found")
	ErrActivePurchaseExists = errors.New("purchase is already pending or approved")
	ErrAlreadyReviewed      = errors.New("purchase is already reviewed")
)

// Contact — как связаться с покупателем, если с чеком что-то не так.
type Contact struct {
	Type  ContactType
	Value string
}

var (
	phoneRe      = regexp.MustCompile(`^\+?[0-9][0-9 ()-]{6,19}$`)
	tgUsernameRe = regexp.MustCompile(`^@?[A-Za-z][A-Za-z0-9_]{4,31}$`)
)

func NewContact(contactType, value string) (Contact, error) {
	value = strings.TrimSpace(value)
	switch t := ContactType(contactType); t {
	case ContactPhone, ContactWhatsApp:
		if !phoneRe.MatchString(value) {
			return Contact{}, fmt.Errorf("%w: expected phone number", ErrInvalidContact)
		}
		return Contact{Type: t, Value: value}, nil
	case ContactTelegram:
		if tgUsernameRe.MatchString(value) {
			return Contact{Type: t, Value: "@" + strings.TrimPrefix(value, "@")}, nil
		}
		if phoneRe.MatchString(value) {
			return Contact{Type: t, Value: value}, nil
		}
		return Contact{}, fmt.Errorf("%w: expected telegram @username or phone number", ErrInvalidContact)
	default:
		return Contact{}, fmt.Errorf("%w: type must be phone, whatsapp or telegram", ErrInvalidContact)
	}
}

type Purchase struct {
	ID                string
	UserID            string
	BookID            string
	Amount            int // минорные единицы, как books.price
	Status            Status
	Contact           Contact
	ReceiptPath       string
	ReceiptMime       string
	TelegramMessageID int64
	ReviewedBy        string
	ReviewedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p *Purchase) IsActive() bool {
	return p.Status == StatusPending || p.Status == StatusApproved
}

// PurchaseDetails — покупка вместе с названием книги и данными покупателя (для списков и уведомлений).
type PurchaseDetails struct {
	Purchase
	BookTitle      string
	BuyerFirstName string
	BuyerLastName  string
	BuyerEmail     string
}

// PaymentSettings — реквизиты, которые видит покупатель перед оплатой.
type PaymentSettings struct {
	Details   string
	QRURL     string
	UpdatedAt time.Time
}
