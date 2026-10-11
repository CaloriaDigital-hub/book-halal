package purchases

import "context"

type Repository interface {
	// Create возвращает ErrActivePurchaseExists, если у пользователя уже есть ожидающая или подтверждённая покупка книги.
	Create(ctx context.Context, p *Purchase) error
	SetTelegramMessage(ctx context.Context, purchaseID string, messageID int64) error
	// Review атомарно переводит покупку из pending в approved/rejected.
	// ErrAlreadyReviewed — если её уже кто-то обработал, ErrPurchaseNotFound — если её нет.
	Review(ctx context.Context, purchaseID string, status Status, reviewedBy string) (*Purchase, error)
	GetDetails(ctx context.Context, purchaseID string) (*PurchaseDetails, error)
	// FindLatest — последняя покупка пользователя по книге или nil.
	FindLatest(ctx context.Context, userID, bookID string) (*Purchase, error)
	HasApproved(ctx context.Context, userID, bookID string) (bool, error)
	// List — покупки по статусу (пустой статус — все), новые сверху.
	List(ctx context.Context, status Status, limit int) ([]PurchaseDetails, error)

	GetSettings(ctx context.Context) (*PaymentSettings, error)
	SaveSettings(ctx context.Context, settings PaymentSettings) error
}
