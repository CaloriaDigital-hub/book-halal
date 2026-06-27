package books

import (
	"context"
	
)

func (r *BookRepository) UpdateCoverURL(ctx context.Context, bookID string, coverURL string) error {
    query := `UPDATE books SET cover_url = $1, updated_at = NOW() WHERE id = $2`
    _, err := r.pool.Exec(ctx, query, coverURL, bookID)
    return err
}