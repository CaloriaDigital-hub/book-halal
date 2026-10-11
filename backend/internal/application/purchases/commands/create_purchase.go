package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	app "book_halal/internal/application/purchases"
	repoBooks "book_halal/internal/domain/books"
	domain "book_halal/internal/domain/purchases"
	uuid "book_halal/internal/pkg"
)

var ErrBookIsFree = errors.New("book is free")

// ReceiptExtensions — допустимые форматы чека.
var ReceiptExtensions = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

// --- DTO ---

type CreatePurchaseCommand struct {
	UserID       string
	BookID       string
	ContactType  string
	ContactValue string
	Receipt      io.Reader
	ReceiptMime  string
}

// --- Interface ---

type CreatePurchaseHandler interface {
	Handle(ctx context.Context, cmd CreatePurchaseCommand) (*app.PurchaseView, error)
}

// --- Handler ---

type createPurchaseHandler struct {
	repo     domain.Repository
	books    repoBooks.Repository
	receipts app.FileStorage
	notifier app.Notifier
	logger   *slog.Logger
}

func NewCreatePurchaseHandler(repo domain.Repository, books repoBooks.Repository, receipts app.FileStorage, notifier app.Notifier, logger *slog.Logger) CreatePurchaseHandler {
	return &createPurchaseHandler{repo: repo, books: books, receipts: receipts, notifier: notifier, logger: logger}
}

func (h *createPurchaseHandler) Handle(ctx context.Context, cmd CreatePurchaseCommand) (*app.PurchaseView, error) {
	contact, err := domain.NewContact(cmd.ContactType, cmd.ContactValue)
	if err != nil {
		return nil, err
	}
	ext, ok := ReceiptExtensions[cmd.ReceiptMime]
	if !ok {
		return nil, fmt.Errorf("unsupported receipt type %q", cmd.ReceiptMime)
	}

	book, err := h.books.GetByID(ctx, cmd.BookID)
	if err != nil {
		return nil, err
	}
	if book.Price.IsFree() {
		return nil, ErrBookIsFree
	}

	// Быстрая проверка до сохранения файла; гонку закрывает уникальный индекс в БД.
	if latest, err := h.repo.FindLatest(ctx, cmd.UserID, cmd.BookID); err != nil {
		return nil, err
	} else if latest != nil && latest.IsActive() {
		return nil, domain.ErrActivePurchaseExists
	}

	id := uuid.New()
	receiptPath := "receipts/" + id + ext
	if err := h.receipts.Save(receiptPath, cmd.Receipt); err != nil {
		return nil, fmt.Errorf("save receipt: %w", err)
	}

	p := &domain.Purchase{
		ID:          id,
		UserID:      cmd.UserID,
		BookID:      cmd.BookID,
		Amount:      book.Price.MinorValue(),
		Status:      domain.StatusPending,
		Contact:     contact,
		ReceiptPath: receiptPath,
		ReceiptMime: cmd.ReceiptMime,
	}
	if err := h.repo.Create(ctx, p); err != nil {
		h.receipts.Remove(receiptPath)
		return nil, err
	}

	h.notify(ctx, p.ID)
	view := app.ToView(p)
	return &view, nil
}

// notify — best effort: если Telegram недоступен, покупка всё равно сохранена
// и её можно подтвердить из админки.
func (h *createPurchaseHandler) notify(ctx context.Context, purchaseID string) {
	logger := h.logger.With("purchase_id", purchaseID)

	details, err := h.repo.GetDetails(ctx, purchaseID)
	if err != nil {
		logger.Error("load purchase for notification", "err", err)
		return
	}
	receipt, err := h.receipts.Open(details.ReceiptPath)
	if err != nil {
		logger.Error("open receipt for notification", "err", err)
		return
	}
	defer receipt.Close()

	notice := app.NoticeFromDetails(details)
	notice.Receipt = receipt
	notice.ReceiptMime = details.ReceiptMime

	messageID, err := h.notifier.PurchaseCreated(ctx, notice)
	if err != nil {
		logger.Error("notify about new purchase", "err", err)
		return
	}
	if messageID != 0 {
		if err := h.repo.SetTelegramMessage(ctx, purchaseID, messageID); err != nil {
			logger.Error("save telegram message id", "err", err)
		}
	}
}
