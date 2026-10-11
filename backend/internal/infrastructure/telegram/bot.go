package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"book_halal/internal/application/purchases/commands"
	domain "book_halal/internal/domain/purchases"
)

// Bot слушает нажатия «Подтвердить» / «Отклонить» через long polling —
// публичный адрес и вебхук не нужны, работает и локально.
type Bot struct {
	client *Client
	chatID int64
	admins map[int64]bool // пусто — может нажимать любой участник группы
	review commands.ReviewPurchaseHandler
	logger *slog.Logger
}

func NewBot(client *Client, chatID int64, adminIDs []int64, review commands.ReviewPurchaseHandler, logger *slog.Logger) *Bot {
	admins := make(map[int64]bool, len(adminIDs))
	for _, id := range adminIDs {
		admins[id] = true
	}
	return &Bot{client: client, chatID: chatID, admins: admins, review: review, logger: logger.With("component", "telegram_bot")}
}

func (b *Bot) Run(ctx context.Context) {
	if me, err := b.client.GetMe(ctx); err != nil {
		b.logger.Error("telegram bot unavailable", "err", err)
	} else {
		b.logger.Info("telegram bot started", "username", me.Username, "chat_id", b.chatID)
	}

	var offset int64
	for ctx.Err() == nil {
		updates, err := b.client.GetUpdates(ctx, offset, 30)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			b.logger.Error("get updates", "err", err)
			sleep(ctx, 5*time.Second)
			continue
		}
		for _, u := range updates {
			offset = u.UpdateID + 1
			switch {
			case u.CallbackQuery != nil:
				b.handleCallback(ctx, u.CallbackQuery)
			case u.Message != nil:
				b.handleMessage(ctx, u.Message)
			}
		}
	}
}

func (b *Bot) handleCallback(ctx context.Context, q *CallbackQuery) {
	purchaseID, approve, ok := parseCallback(q.Data)
	if !ok {
		return
	}
	if q.Message == nil || q.Message.Chat.ID != b.chatID {
		b.answer(ctx, q, "Эта кнопка работает только в группе оплат", true)
		return
	}
	if len(b.admins) > 0 && !b.admins[q.From.ID] {
		b.answer(ctx, q, "Нет прав подтверждать оплаты", true)
		return
	}

	reviewer := reviewerName(q.From)
	_, err := b.review.Handle(ctx, commands.ReviewPurchaseCommand{PurchaseID: purchaseID, Approve: approve, ReviewedBy: reviewer})
	switch {
	case errors.Is(err, domain.ErrAlreadyReviewed):
		b.answer(ctx, q, "Эта оплата уже обработана", true)
	case errors.Is(err, domain.ErrPurchaseNotFound):
		b.answer(ctx, q, "Оплата не найдена", true)
	case err != nil:
		b.logger.Error("review purchase", "purchase_id", purchaseID, "err", err)
		b.answer(ctx, q, "Ошибка, попробуйте ещё раз", true)
	case approve:
		b.logger.Info("purchase approved", "purchase_id", purchaseID, "by", reviewer)
		b.answer(ctx, q, "Оплата подтверждена ✅", false)
	default:
		b.logger.Info("purchase rejected", "purchase_id", purchaseID, "by", reviewer)
		b.answer(ctx, q, "Оплата отклонена ❌", false)
	}
}

// /ids — подсказка для настройки: chat_id группы и ID темы, где написана команда.
func (b *Bot) handleMessage(ctx context.Context, m *Message) {
	cmd := strings.SplitN(strings.TrimSpace(m.Text), "@", 2)[0]
	if cmd != "/ids" {
		return
	}
	text := fmt.Sprintf("chat_id: <code>%d</code>\nmessage_thread_id: <code>%d</code>", m.Chat.ID, m.MessageThreadID)
	if topic := m.TopicName(); topic != "" {
		text += "\nтема: " + esc(topic)
	}
	var fromID int64
	if m.From != nil {
		fromID = m.From.ID
		text += fmt.Sprintf("\nваш user id: <code>%d</code>", m.From.ID)
	}
	b.logger.Info("telegram ids requested", "chat_id", m.Chat.ID, "chat_title", m.Chat.Title,
		"message_thread_id", m.MessageThreadID, "topic", m.TopicName(), "from_id", fromID)
	if _, err := b.client.SendMessage(ctx, m.Chat.ID, m.MessageThreadID, text); err != nil {
		b.logger.Error("reply /ids", "err", err)
	}
}

func (b *Bot) answer(ctx context.Context, q *CallbackQuery, text string, alert bool) {
	if err := b.client.AnswerCallback(ctx, q.ID, text, alert); err != nil {
		b.logger.Error("answer callback", "err", err)
	}
}

func sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
