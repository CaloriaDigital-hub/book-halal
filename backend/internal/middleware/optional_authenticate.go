package middleware

import (
	"net/http"
	"strings"

	sessionDomain "book_halal/internal/domain/sessions"
	userDomain "book_halal/internal/domain/users"
	"book_halal/internal/domain/users/value_objects"
)

// OptionalAuthenticate — для публичных эндпоинтов (каталог, карточка книги):
// с валидным токеном кладёт пользователя в контекст, без токена или с протухшим — пропускает анонимно.
func OptionalAuthenticate(sessionRepo sessionDomain.Repository, userRepo userDomain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				next.ServeHTTP(w, r)
				return
			}
			session, err := sessionRepo.FindByToken(r.Context(), token)
			if err != nil || session == nil || session.IsExpired() {
				next.ServeHTTP(w, r)
				return
			}
			userID, err := valueobjects.NewUserId(session.UserID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user, err := userRepo.FindByID(r.Context(), userID)
			if err != nil || user == nil {
				next.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}
}
