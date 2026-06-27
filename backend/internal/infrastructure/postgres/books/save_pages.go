package books

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	book "book_halal/internal/domain/books/entity"
)
// SavePages переписан на pgx.Batch для максимальной производительности
func (r *BookRepository) SavePages(ctx context.Context, pages []book.Page) error {
	if len(pages) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO book_pages (id, book_id, page_number, image_url, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (book_id, page_number) DO NOTHING
	`

	// Добавляем каждый запрос в пачку
	for _, p := range pages {
		batch.Queue(query, p.ID, p.BookID, p.PageNumber, p.ImageURL)
	}

	// Отправляем всю пачку за один сетевой вызов
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	// Обязательно проверяем выполнение каждого запроса в пачке
	for i := 0; i < len(pages); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch insert failed at page index %d: %w", i, err)
		}
	}

	return nil
}