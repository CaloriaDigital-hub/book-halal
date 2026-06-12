package queries

// BookDetailsView — DTO для детальной страницы книги
type BookDetailsView struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       int    `json:"price"` 
	CoverURL    string `json:"cover_url"`
	TotalPages  int    `json:"total_pages"` // Добавили, чтобы фронт знал, сколько всего страниц
	Status      string `json:"status"`      // На всякий случай (ready, processing)
}