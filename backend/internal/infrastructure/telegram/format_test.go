package telegram

import (
	"strings"
	"testing"
	"time"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

func TestParseCallback(t *testing.T) {
	id, approve, ok := parseCallback("pay:approve:abc")
	if !ok || !approve || id != "abc" {
		t.Fatalf("approve: got %q %v %v", id, approve, ok)
	}
	id, approve, ok = parseCallback("pay:reject:abc")
	if !ok || approve || id != "abc" {
		t.Fatalf("reject: got %q %v %v", id, approve, ok)
	}
	if _, _, ok := parseCallback("something:else"); ok {
		t.Fatal("foreign callback must be ignored")
	}
}

func TestCallbackDataFitsTelegramLimit(t *testing.T) {
	// callback_data в Telegram — не больше 64 байт
	for _, row := range reviewKeyboard("15e01377-caf6-444c-be2b-953da693eb75").InlineKeyboard {
		for _, btn := range row {
			if len(btn.CallbackData) > 64 {
				t.Errorf("%q: %d bytes", btn.CallbackData, len(btn.CallbackData))
			}
		}
	}
}

func TestPurchaseCaption(t *testing.T) {
	n := app.Notice{
		PurchaseID: "15e01377-caf6-444c-be2b-953da693eb75",
		BookTitle:  "Намаз <опора>",
		BuyerName:  "Иван Иванов",
		BuyerEmail: "ivan@example.com",
		Amount:     1500,
		Contact:    domain.Contact{Type: domain.ContactWhatsApp, Value: "+996 555 123 456"},
		Status:     domain.StatusPending,
		CreatedAt:  time.Date(2026, 10, 11, 14, 0, 0, 0, time.UTC),
	}
	got := purchaseCaption(n)
	for _, want := range []string{"#15e01377", "Намаз &lt;опора&gt;", "1 500 сом", "https://wa.me/996555123456", "11.10.2026 20:00", "Ожидает проверки"} {
		if !strings.Contains(got, want) {
			t.Errorf("caption has no %q:\n%s", want, got)
		}
	}

	n.Status, n.ReviewedBy, n.ReviewedAt = domain.StatusApproved, "@admin", n.CreatedAt
	if got := purchaseCaption(n); !strings.Contains(got, "Подтверждено</b> — @admin") {
		t.Errorf("approved caption:\n%s", got)
	}
}

func TestFormatAmount(t *testing.T) {
	for in, want := range map[int]string{0: "0", 999: "999", 1500: "1 500", 1234567: "1 234 567"} {
		if got := formatAmount(in); got != want {
			t.Errorf("formatAmount(%d) = %q, want %q", in, got, want)
		}
	}
}
