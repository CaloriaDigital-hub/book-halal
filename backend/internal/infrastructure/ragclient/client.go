package ragclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	entity "book_halal/internal/domain/books/entity"
)

type RAGClient struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *RAGClient {
	return &RAGClient{baseURL: baseURL, client: &http.Client{}}
}

// --- Index by page images ---

type indexPage struct {
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}

type indexByPagesRequest struct {
	Pages []indexPage `json:"pages"`
}

// IndexBook sends page image URLs to the RAG service.
// Python will download each image and run OCR locally.
func (r *RAGClient) IndexBook(bookID string, baseURL string, pages []entity.Page) error {
	req := indexByPagesRequest{}
	for _, p := range pages {
		req.Pages = append(req.Pages, indexPage{
			PageNumber: p.PageNumber,
			// image_url must be absolute so Python can download it
			ImageURL: baseURL + p.ImageURL,
		})
	}

	body, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/admin/books/%s/index", r.baseURL, bookID)
	resp, err := r.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rag index request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Ask ---

type askRequest struct {
	Question string `json:"question"`
}

type AskSource struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}

type AskResponse struct {
	Answer  string      `json:"answer"`
	Sources []AskSource `json:"sources"`
}

func (r *RAGClient) AskBook(bookID string, question string) (*AskResponse, error) {
	body, _ := json.Marshal(askRequest{Question: question})
	url := fmt.Sprintf("%s/books/%s/ask", r.baseURL, bookID)

	resp, err := r.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("rag service unavailable: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rag response: %w", err)
	}

	var result AskResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse rag response: %w", err)
	}

	return &result, nil
}