package telegram

import (
	"context"

	app "book_halal/internal/application/purchases"
)

// Notifier кладёт чеки в тему «Платежи», а итоги проверки — в «Историю».
type Notifier struct {
	client         *Client
	chatID         int64
	paymentsThread int64
	historyThread  int64
}

func NewNotifier(client *Client, chatID, paymentsThread, historyThread int64) *Notifier {
	return &Notifier{client: client, chatID: chatID, paymentsThread: paymentsThread, historyThread: historyThread}
}

func (n *Notifier) PurchaseCreated(ctx context.Context, notice app.Notice) (int64, error) {
	fileName := "receipt-" + shortID(notice.PurchaseID) + extFor(notice.ReceiptMime)
	msg, err := n.client.SendFile(ctx, n.chatID, n.paymentsThread, notice.Receipt, fileName,
		notice.ReceiptMime, purchaseCaption(notice), reviewKeyboard(notice.PurchaseID))
	if err != nil {
		return 0, err
	}
	return msg.MessageID, nil
}

func (n *Notifier) PurchaseReviewed(ctx context.Context, notice app.Notice) error {
	// Сообщение с чеком: обновляем статус и убираем кнопки
	if notice.TelegramMessageID != 0 {
		if err := n.client.EditCaption(ctx, n.chatID, notice.TelegramMessageID, purchaseCaption(notice), nil); err != nil {
			return err
		}
	}
	_, err := n.client.SendMessage(ctx, n.chatID, n.historyThread, historyLine(notice))
	return err
}

func extFor(mime string) string {
	switch mime {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "application/pdf":
		return ".pdf"
	default:
		return ".jpg"
	}
}
