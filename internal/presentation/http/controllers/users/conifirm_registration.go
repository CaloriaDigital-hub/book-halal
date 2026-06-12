package controllers

import (
	"encoding/json"
	"net/http"

	"book_halal/internal/application/users/commands"
)

type confirmRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

// ConfirmRegistration handles POST /api/auth/register/confirm
func (c *AuthController) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	var req confirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cmd := commands.ConfirmRegistrationCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Code:      req.Code,
	}

	if err := c.confirmUseCase.Handle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user successfully registered",
	})
}