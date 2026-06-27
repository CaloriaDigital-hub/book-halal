package books

import (
	"context"
	"fmt"

	"book_halal/internal/domain/books/entity"
	"book_halal/internal/domain/books/value_objects" 
)


func (r *BookRepository) GetAllReady(ctx context.Context) ([]entity.Book, error) {
	query := `
		SELECT id, title, author, description, price, cover_url, total_pages, status, created_at, updated_at
		FROM books
		WHERE status = 'ready'
		ORDER BY created_at DESC
	`
	
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all ready books: %w", err)
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var b entity.Book
		var priceMinor int
		
		
		var description, coverURL *string

		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &description, &priceMinor, &coverURL,
			&b.TotalPages, &b.Status, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan ready book: %w", err)
		}

		
		if description != nil {
			b.Description = *description
		}
		if coverURL != nil {
			b.CoverURL = *coverURL
		}

		// Собираем VO цены ИЗ минорных единиц (копеек/тыйынов)
		priceVO, _ := valueobjects.NewPriceFromMinor(priceMinor)
		b.Price = priceVO

		books = append(books, b)
	}

	return books, nil
}