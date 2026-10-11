package telegram

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	app "book_halal/internal/application/purchases"
	"book_halal/internal/application/purchases/commands"
	domain "book_halal/internal/domain/purchases"
)

// fakeAPI записывает вызовы Bot API и отвечает успехом.
type fakeAPI struct {
	mu    sync.Mutex
	calls []fakeCall
}

type fakeCall struct {
	method string
	fields map[string]string
	file   string
}

func (f *fakeAPI) server(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := fakeCall{method: r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:], fields: map[string]string{}}
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Errorf("parse multipart: %v", err)
			}
			for k, v := range r.MultipartForm.Value {
				call.fields[k] = v[0]
			}
			for k, files := range r.MultipartForm.File {
				f, _ := files[0].Open()
				b, _ := io.ReadAll(f)
				call.file = k + ":" + files[0].Filename + ":" + string(b)
			}
		} else {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			for k, v := range body {
				b, _ := json.Marshal(v)
				call.fields[k] = strings.Trim(string(b), `"`)
			}
		}
		f.mu.Lock()
		f.calls = append(f.calls, call)
		f.mu.Unlock()
		w.Write([]byte(`{"ok":true,"result":{"message_id":42,"chat":{"id":-100}}}`))
	}))
}

func testClient(srv *httptest.Server) *Client {
	c := NewClient("TOKEN")
	c.baseURL = srv.URL + "/botTOKEN/"
	return c
}

func notice() app.Notice {
	return app.Notice{
		PurchaseID: "15e01377-caf6-444c-be2b-953da693eb75",
		BookTitle:  "Намаз",
		BuyerName:  "Иван Иванов",
		BuyerEmail: "ivan@example.com",
		Amount:     1500,
		Contact:    domain.Contact{Type: domain.ContactTelegram, Value: "@ivan_buyer"},
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}
}

func TestNotifierSendsReceiptToPaymentsTopic(t *testing.T) {
	api := &fakeAPI{}
	srv := api.server(t)
	defer srv.Close()

	n := NewNotifier(testClient(srv), -100, 7, 9)
	in := notice()
	in.Receipt, in.ReceiptMime = strings.NewReader("JPEGDATA"), "image/jpeg"

	id, err := n.PurchaseCreated(context.Background(), in)
	if err != nil || id != 42 {
		t.Fatalf("PurchaseCreated = %d, %v", id, err)
	}
	c := api.calls[0]
	if c.method != "sendPhoto" || c.fields["chat_id"] != "-100" || c.fields["message_thread_id"] != "7" {
		t.Fatalf("unexpected call: %+v", c)
	}
	if c.file != "photo:receipt-15e01377.jpg:JPEGDATA" {
		t.Errorf("file = %q", c.file)
	}
	if !strings.Contains(c.fields["reply_markup"], "pay:approve:15e01377-caf6-444c-be2b-953da693eb75") {
		t.Errorf("no approve button: %s", c.fields["reply_markup"])
	}
	if !strings.Contains(c.fields["caption"], "https://t.me/ivan_buyer") {
		t.Errorf("caption: %s", c.fields["caption"])
	}
}

func TestNotifierSendsPDFAsDocument(t *testing.T) {
	api := &fakeAPI{}
	srv := api.server(t)
	defer srv.Close()

	in := notice()
	in.Receipt, in.ReceiptMime = strings.NewReader("%PDF"), "application/pdf"
	if _, err := NewNotifier(testClient(srv), -100, 7, 9).PurchaseCreated(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if api.calls[0].method != "sendDocument" || !strings.HasPrefix(api.calls[0].file, "document:receipt-15e01377.pdf") {
		t.Fatalf("unexpected call: %+v", api.calls[0])
	}
}

func TestNotifierReviewedEditsMessageAndWritesHistory(t *testing.T) {
	api := &fakeAPI{}
	srv := api.server(t)
	defer srv.Close()

	in := notice()
	in.Status, in.ReviewedBy, in.ReviewedAt, in.TelegramMessageID = domain.StatusApproved, "@admin", time.Now(), 42
	if err := NewNotifier(testClient(srv), -100, 7, 9).PurchaseReviewed(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if len(api.calls) != 2 {
		t.Fatalf("calls: %+v", api.calls)
	}
	edit, history := api.calls[0], api.calls[1]
	if edit.method != "editMessageCaption" || edit.fields["message_id"] != "42" || !strings.Contains(edit.fields["caption"], "Подтверждено") {
		t.Errorf("edit: %+v", edit)
	}
	if edit.fields["reply_markup"] != `{"inline_keyboard":[]}` {
		t.Errorf("buttons must be removed, got %s", edit.fields["reply_markup"])
	}
	if history.method != "sendMessage" || history.fields["message_thread_id"] != "9" || !strings.Contains(history.fields["text"], "подтверждена — @admin") {
		t.Errorf("history: %+v", history)
	}
}

// fakeReview — подмена команды подтверждения.
type fakeReview struct {
	got []commands.ReviewPurchaseCommand
	err error
}

func (f *fakeReview) Handle(_ context.Context, cmd commands.ReviewPurchaseCommand) (*app.AdminPurchaseView, error) {
	f.got = append(f.got, cmd)
	return &app.AdminPurchaseView{}, f.err
}

func callback(chatID, fromID int64, data string) *CallbackQuery {
	return &CallbackQuery{
		ID:      "cb1",
		From:    User{ID: fromID, Username: "admin"},
		Message: &Message{MessageID: 42, Chat: Chat{ID: chatID}},
		Data:    data,
	}
}

func TestBotCallbacks(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()

	t.Run("approve from group", func(t *testing.T) {
		api, review := &fakeAPI{}, &fakeReview{}
		srv := api.server(t)
		defer srv.Close()
		NewBot(testClient(srv), -100, nil, review, logger).handleCallback(ctx, callback(-100, 1, "pay:approve:p1"))
		if len(review.got) != 1 || !review.got[0].Approve || review.got[0].PurchaseID != "p1" || review.got[0].ReviewedBy != "@admin" {
			t.Fatalf("review: %+v", review.got)
		}
		if !strings.Contains(api.calls[0].fields["text"], "подтверждена") {
			t.Errorf("answer: %+v", api.calls[0])
		}
	})

	t.Run("button from another chat is ignored", func(t *testing.T) {
		api, review := &fakeAPI{}, &fakeReview{}
		srv := api.server(t)
		defer srv.Close()
		NewBot(testClient(srv), -100, nil, review, logger).handleCallback(ctx, callback(-555, 1, "pay:approve:p1"))
		if len(review.got) != 0 {
			t.Fatal("review must not run for foreign chat")
		}
	})

	t.Run("not in admin list", func(t *testing.T) {
		api, review := &fakeAPI{}, &fakeReview{}
		srv := api.server(t)
		defer srv.Close()
		NewBot(testClient(srv), -100, []int64{777}, review, logger).handleCallback(ctx, callback(-100, 1, "pay:reject:p1"))
		if len(review.got) != 0 || !strings.Contains(api.calls[0].fields["text"], "Нет прав") {
			t.Fatalf("review=%+v calls=%+v", review.got, api.calls)
		}
	})

	t.Run("already reviewed", func(t *testing.T) {
		api, review := &fakeAPI{}, &fakeReview{err: domain.ErrAlreadyReviewed}
		srv := api.server(t)
		defer srv.Close()
		NewBot(testClient(srv), -100, nil, review, logger).handleCallback(ctx, callback(-100, 1, "pay:reject:p1"))
		if !strings.Contains(api.calls[0].fields["text"], "уже обработана") {
			t.Fatalf("answer: %+v", api.calls[0])
		}
	})
}
