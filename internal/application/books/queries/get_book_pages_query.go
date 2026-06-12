package queries

// PageView — DTO для одной страницы
type PageView struct {
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}

// НОВЫЙ DTO — Обертка для ответа
type BookPagesResponse struct {
	TotalPages int        `json:"total_pages"`
	Pages      []PageView `json:"pages"`
}