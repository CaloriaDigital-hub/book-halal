package middleware

import (
	"net/http"
	"strings"

	sessionDomain "book_halal/internal/domain/sessions"
	userDomain "book_halal/internal/domain/users"
	"book_halal/internal/domain/users/value_objects"
)

func Authenticate(sessionRepo sessionDomain.Repository, userRepo userDomain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || token == "" {
				http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
				return
			}

			session, err := sessionRepo.FindByToken(r.Context(), token)
			if err != nil || session == nil {
				http.Error(w, `{"error":"invalid session"}`, http.StatusUnauthorized)
				return
			}

			if session.IsExpired() {
				http.Error(w, `{"error":"session expired"}`, http.StatusUnauthorized)
				return
			}

			userID, err := valueobjects.NewUserId(session.UserID)
			if err != nil {
				http.Error(w, `{"error":"invalid session"}`, http.StatusUnauthorized)
				return
			}

			user, err := userRepo.FindByID(r.Context(), userID)
			if err != nil || user == nil {
				http.Error(w, `{"error":"user not found"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}