package reader

import (
	"context"
	"time"
)


type ReadingProgress struct {
	UserID     string
	BookID     string
	PageNumber int
	UpdatedAt  time.Time
}


type Bookmark struct {
	ID         string
	UserID     string
	BookID     string
	PageNumber int
	CreatedAt  time.Time
}


type Repository interface {
	// Progress
	UpsertProgress(ctx context.Context, userID, bookID string, page int) error
	GetProgress(ctx context.Context, userID, bookID string) (*ReadingProgress, error)

	// Bookmarks
	AddBookmark(ctx context.Context, userID, bookID string, page int) (*Bookmark, error)
	GetBookmarks(ctx context.Context, userID, bookID string) ([]Bookmark, error)
	DeleteBookmark(ctx context.Context, userID, bookID string, page int) error
}
