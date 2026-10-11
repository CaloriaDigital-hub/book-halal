package commands

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	app "book_halal/internal/application/purchases"
	domain "book_halal/internal/domain/purchases"
)

var QRExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// --- DTO ---

type UpdatePaymentSettingsCommand struct {
	Details string
	// QR — новая картинка QR-кода; nil — оставить текущую.
	QR     io.Reader
	QRMime string
	// RemoveQR — убрать QR-код совсем.
	RemoveQR bool
}

// --- Interface ---

type UpdatePaymentSettingsHandler interface {
	Handle(ctx context.Context, cmd UpdatePaymentSettingsCommand) error
}

// --- Handler ---

type updatePaymentSettingsHandler struct {
	repo   domain.Repository
	static app.FileStorage // корень static: QR раздаётся публично по /static/...
}

func NewUpdatePaymentSettingsHandler(repo domain.Repository, static app.FileStorage) UpdatePaymentSettingsHandler {
	return &updatePaymentSettingsHandler{repo: repo, static: static}
}

func (h *updatePaymentSettingsHandler) Handle(ctx context.Context, cmd UpdatePaymentSettingsCommand) error {
	current, err := h.repo.GetSettings(ctx)
	if err != nil {
		return err
	}

	settings := domain.PaymentSettings{Details: strings.TrimSpace(cmd.Details), QRURL: current.QRURL}
	if cmd.RemoveQR {
		settings.QRURL = ""
	}
	if cmd.QR != nil {
		ext, ok := QRExtensions[cmd.QRMime]
		if !ok {
			return fmt.Errorf("unsupported QR image type %q", cmd.QRMime)
		}
		// Новое имя на каждую загрузку — чтобы клиенты не видели старый QR из кэша.
		relPath := fmt.Sprintf("payments/qr-%d%s", time.Now().Unix(), ext)
		if err := h.static.Save(relPath, cmd.QR); err != nil {
			return fmt.Errorf("save QR: %w", err)
		}
		settings.QRURL = "/static/" + relPath
	}

	if err := h.repo.SaveSettings(ctx, settings); err != nil {
		return err
	}
	if current.QRURL != "" && current.QRURL != settings.QRURL {
		h.static.Remove(strings.TrimPrefix(current.QRURL, "/static/"))
	}
	return nil
}
