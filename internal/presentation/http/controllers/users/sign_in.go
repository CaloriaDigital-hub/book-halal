package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"book_halal/internal/application/users/commands"
	userHandlers "book_halal/internal/application/users/commands/handlers"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cmd := commands.SignInCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := c.signInUseCase.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, userHandlers.ErrInvalidCredentials) {
			http.Error(w, `{"error":"invalid email or password"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error":"failed to sign in"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token":      result.Token,
		"expires_at": result.ExpiresAt,
	})
}