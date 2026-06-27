package controllers

import (
	"encoding/json"
	"net/http"

	"book_halal/internal/application/users/commands"
)

type initiateRequest struct {
	Email string `json:"email"`
}

// InitiateRegistration handles POST /api/auth/register/initiate
func (c *AuthController) InitiateRegistration(w http.ResponseWriter, r *http.Request) {
	var req initiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cmd := commands.InitiateRegistrationCommand{
		Email: req.Email,
	}

	if err := c.initiateUseCase.Handle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "verification code sent to email",
	})
}