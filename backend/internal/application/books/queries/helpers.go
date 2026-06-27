package queries

import (
	"book_halal/internal/domain/books/entity"
)

func toPageViews(pages []entity.Page) []PageView {
	result := make([]PageView, 0, len(pages))
	for _, p := range pages {
		result = append(result, PageView{
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		})
	}
	return result
}
