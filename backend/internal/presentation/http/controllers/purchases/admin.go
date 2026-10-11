package purchases

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"book_halal/internal/application/purchases/commands"
	domain "book_halal/internal/domain/purchases"
	"book_halal/internal/middleware"
)

const maxQRSize = 5 << 20 // 5 МБ

// ListPurchases — GET /api/admin/purchases?status=pending&limit=50
func (c *PurchaseController) ListPurchases(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	views, err := c.list.Handle(r.Context(), r.URL.Query().Get("status"), limit)
	if err != nil {
		c.sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.sendJSON(w, http.StatusOK, views)
}

// ApprovePurchase — POST /api/admin/purchases/{id}/approve
func (c *PurchaseController) ApprovePurchase(w http.ResponseWriter, r *http.Request) {
	c.reviewPurchase(w, r, true)
}

// RejectPurchase — POST /api/admin/purchases/{id}/reject
func (c *PurchaseController) RejectPurchase(w http.ResponseWriter, r *http.Request) {
	c.reviewPurchase(w, r, false)
}

func (c *PurchaseController) reviewPurchase(w http.ResponseWriter, r *http.Request, approve bool) {
	admin, _ := middleware.UserFromContext(r.Context())
	view, err := c.review.Handle(r.Context(), commands.ReviewPurchaseCommand{
		PurchaseID: r.PathValue("id"),
		Approve:    approve,
		ReviewedBy: admin.Email.String(),
	})
	switch {
	case errors.Is(err, domain.ErrPurchaseNotFound):
		c.sendError(w, http.StatusNotFound, "purchase not found")
	case errors.Is(err, domain.ErrAlreadyReviewed):
		c.sendError(w, http.StatusConflict, "purchase is already reviewed")
	case err != nil:
		c.logger.Error("review purchase", "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to review purchase")
	default:
		c.sendJSON(w, http.StatusOK, view)
	}
}

// GetReceipt — GET /api/admin/purchases/{id}/receipt: файл чека (только админам).
func (c *PurchaseController) GetReceipt(w http.ResponseWriter, r *http.Request) {
	file, mime, err := c.receipt.Handle(r.Context(), r.PathValue("id"))
	if errors.Is(err, domain.ErrPurchaseNotFound) {
		c.sendError(w, http.StatusNotFound, "purchase not found")
		return
	}
	if err != nil {
		c.logger.Error("open receipt", "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to open receipt")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	io.Copy(w, file)
}

// UpdateRequisites — PUT /api/admin/payments/requisites (multipart: details, qr?, remove_qr?)
func (c *PurchaseController) UpdateRequisites(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxQRSize+1<<20)
	if err := r.ParseMultipartForm(maxQRSize); err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid form data or QR image larger than 5 MB")
		return
	}
	qr, mime, err := readUpload(r, "qr", commands.QRExtensions)
	if errors.Is(err, errUnsupportedType) {
		c.sendError(w, http.StatusBadRequest, "QR must be jpg, png or webp")
		return
	}
	if err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid QR file")
		return
	}
	cmd := commands.UpdatePaymentSettingsCommand{
		Details:  r.FormValue("details"),
		RemoveQR: r.FormValue("remove_qr") == "true",
	}
	if qr != nil {
		defer qr.Close()
		cmd.QR, cmd.QRMime = qr, mime
	}

	if err := c.updateSettings.Handle(r.Context(), cmd); err != nil {
		c.logger.Error("update payment settings", "err", err)
		c.sendError(w, http.StatusInternalServerError, "failed to save requisites")
		return
	}
	c.GetRequisites(w, r)
}
