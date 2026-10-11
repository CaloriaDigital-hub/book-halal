package purchases

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"book_halal/internal/application/purchases/commands"
	"book_halal/internal/application/purchases/queries"
)

type PurchaseController struct {
	create         commands.CreatePurchaseHandler
	review         commands.ReviewPurchaseHandler
	updateSettings commands.UpdatePaymentSettingsHandler
	myPurchase     queries.GetMyPurchaseHandler
	settings       queries.GetPaymentSettingsHandler
	list           queries.ListPurchasesHandler
	receipt        queries.GetReceiptHandler
	logger         *slog.Logger
}

func NewPurchaseController(
	create commands.CreatePurchaseHandler,
	review commands.ReviewPurchaseHandler,
	updateSettings commands.UpdatePaymentSettingsHandler,
	myPurchase queries.GetMyPurchaseHandler,
	settings queries.GetPaymentSettingsHandler,
	list queries.ListPurchasesHandler,
	receipt queries.GetReceiptHandler,
	logger *slog.Logger,
) *PurchaseController {
	return &PurchaseController{
		create:         create,
		review:         review,
		updateSettings: updateSettings,
		myPurchase:     myPurchase,
		settings:       settings,
		list:           list,
		receipt:        receipt,
		logger:         logger,
	}
}

func (c *PurchaseController) sendError(w http.ResponseWriter, status int, msg string) {
	c.sendJSON(w, status, map[string]string{"error": msg})
}

func (c *PurchaseController) sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
