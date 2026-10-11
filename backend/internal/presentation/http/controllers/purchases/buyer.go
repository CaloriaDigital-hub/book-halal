package purchases

import (
	"errors"
	"net/http"

	"book_halal/internal/application/purchases/commands"
	repoBooks "book_halal/internal/domain/books"
	domain "book_halal/internal/domain/purchases"
	"book_halal/internal/middleware"
)

const maxReceiptSize = 10 << 20 // 10 МБ

// GetRequisites — GET /api/payments/requisites: реквизиты и QR для оплаты.
func (c *PurchaseController) GetRequisites(w http.ResponseWriter, r *http.Request) {
	view, err := c.settings.Handle(r.Context())
	if err != nil {
		c.logger.Error("get payment settings", "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to load payment requisites")
		return
	}
	c.sendJSON(w, http.StatusOK, view)
}

// CreatePurchase — POST /api/books/{id}/purchase (multipart: receipt, contact_type, contact).
func (c *PurchaseController) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		c.sendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	bookID := r.PathValue("id")

	r.Body = http.MaxBytesReader(w, r.Body, maxReceiptSize+1<<20)
	if err := r.ParseMultipartForm(maxReceiptSize); err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid form data or receipt larger than 10 MB")
		return
	}
	receipt, mime, err := readUpload(r, "receipt", commands.ReceiptExtensions)
	if errors.Is(err, errUnsupportedType) {
		c.sendError(w, http.StatusBadRequest, "receipt must be jpg, png, webp or pdf")
		return
	}
	if err != nil || receipt == nil {
		c.sendError(w, http.StatusBadRequest, "receipt file is required")
		return
	}
	defer receipt.Close()

	view, err := c.create.Handle(r.Context(), commands.CreatePurchaseCommand{
		UserID:       user.ID.String(),
		BookID:       bookID,
		ContactType:  r.FormValue("contact_type"),
		ContactValue: r.FormValue("contact"),
		Receipt:      receipt,
		ReceiptMime:  mime,
	})
	switch {
	case errors.Is(err, domain.ErrInvalidContact):
		c.sendError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repoBooks.ErrBookNotFound):
		c.sendError(w, http.StatusNotFound, "book not found")
	case errors.Is(err, commands.ErrBookIsFree):
		c.sendError(w, http.StatusBadRequest, "book is free, no payment needed")
	case errors.Is(err, domain.ErrActivePurchaseExists):
		c.sendError(w, http.StatusConflict, "payment for this book is already pending or approved")
	case err != nil:
		c.logger.Error("create purchase", "book_id", bookID, "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to create purchase")
	default:
		c.sendJSON(w, http.StatusCreated, view)
	}
}

// GetMyPurchase — GET /api/books/{id}/purchase: статус последней оплаты книги.
func (c *PurchaseController) GetMyPurchase(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		c.sendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	view, err := c.myPurchase.Handle(r.Context(), user.ID.String(), r.PathValue("id"))
	if err != nil {
		c.logger.Error("get my purchase", "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to load purchase")
		return
	}
	if view == nil {
		c.sendError(w, http.StatusNotFound, "no purchase for this book")
		return
	}
	c.sendJSON(w, http.StatusOK, view)
}
