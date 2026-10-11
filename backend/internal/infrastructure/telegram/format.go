package telegram

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

const (
	callbackApprove = "pay:approve:"
	callbackReject  = "pay:reject:"
	captionLimit    = 1024 // лимит Telegram на подпись к фото/документу
)

var nonDigits = regexp.MustCompile(`\D`)

// Время в подписях — по Бишкеку.
var bishkek = time.FixedZone("Asia/Bishkek", 6*60*60)

func reviewKeyboard(purchaseID string) *InlineKeyboard {
	return &InlineKeyboard{InlineKeyboard: [][]InlineButton{{
		{Text: "✅ Подтвердить", CallbackData: callbackApprove + purchaseID},
		{Text: "❌ Отклонить", CallbackData: callbackReject + purchaseID},
	}}}
}

// parseCallback разбирает callback_data кнопок. ok=false — чужая кнопка.
func parseCallback(data string) (purchaseID string, approve bool, ok bool) {
	switch {
	case strings.HasPrefix(data, callbackApprove):
		return strings.TrimPrefix(data, callbackApprove), true, true
	case strings.HasPrefix(data, callbackReject):
		return strings.TrimPrefix(data, callbackReject), false, true
	}
	return "", false, false
}

func purchaseCaption(n app.Notice) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🧾 <b>Оплата книги</b> · #%s\n\n", shortID(n.PurchaseID))
	fmt.Fprintf(&b, "📖 %s\n", esc(n.BookTitle))
	fmt.Fprintf(&b, "💰 <b>%s сом</b>\n", formatAmount(n.Amount))
	fmt.Fprintf(&b, "👤 %s · %s\n", esc(strings.TrimSpace(n.BuyerName)), esc(n.BuyerEmail))
	fmt.Fprintf(&b, "📞 %s\n", contactLine(n.Contact))
	fmt.Fprintf(&b, "🕓 %s\n\n", n.CreatedAt.In(bishkek).Format("02.01.2006 15:04"))
	b.WriteString(statusLine(n))
	return truncate(b.String(), captionLimit)
}

func historyLine(n app.Notice) string {
	icon, verb := "✅", "подтверждена"
	if n.Status == domain.StatusRejected {
		icon, verb = "❌", "отклонена"
	}
	return fmt.Sprintf("%s Оплата #%s %s — %s\n📖 %s · 💰 %s сом\n👤 %s · %s",
		icon, shortID(n.PurchaseID), verb, esc(n.ReviewedBy),
		esc(n.BookTitle), formatAmount(n.Amount),
		esc(strings.TrimSpace(n.BuyerName)), contactLine(n.Contact))
}

func statusLine(n app.Notice) string {
	when := n.ReviewedAt.In(bishkek).Format("02.01 15:04")
	switch n.Status {
	case domain.StatusApproved:
		return fmt.Sprintf("✅ <b>Подтверждено</b> — %s, %s", esc(n.ReviewedBy), when)
	case domain.StatusRejected:
		return fmt.Sprintf("❌ <b>Отклонено</b> — %s, %s", esc(n.ReviewedBy), when)
	default:
		return "⏳ Ожидает проверки"
	}
}

func contactLine(c domain.Contact) string {
	v := esc(c.Value)
	switch c.Type {
	case domain.ContactTelegram:
		if strings.HasPrefix(c.Value, "@") {
			return fmt.Sprintf(`Telegram: <a href="https://t.me/%s">%s</a>`, esc(strings.TrimPrefix(c.Value, "@")), v)
		}
		return "Telegram: " + v
	case domain.ContactWhatsApp:
		return fmt.Sprintf(`WhatsApp: <a href="https://wa.me/%s">%s</a>`, nonDigits.ReplaceAllString(c.Value, ""), v)
	default:
		return "Телефон: " + v
	}
}

// formatAmount: 1500 → "1 500".
func formatAmount(n int) string {
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ' ')
		}
		out = append(out, s[i])
	}
	return string(out)
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

func esc(s string) string { return html.EscapeString(s) }

func truncate(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit-1]) + "…"
}

// reviewerName — как подписать того, кто нажал кнопку.
func reviewerName(u User) string {
	if u.Username != "" {
		return "@" + u.Username
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	return fmt.Sprintf("%s (id %d)", name, u.ID)
}
