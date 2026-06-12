package entity

import "time"

type Session struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewSession(id, userID, token string, expiresAt time.Time) *Session {
	return &Session{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}