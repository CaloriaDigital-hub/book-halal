package books

import (
	"book_halal/internal/application/books/commands"
	"book_halal/internal/application/books/queries"
	"book_halal/internal/infrastructure/ragclient"
)

type BookController struct {
	uploadHandler  commands.UploadBookHandler
	catalogHandler queries.GetBooksHandler
	detailsHandler queries.GetBookByIDHandler
	pagesHandler   queries.GetBookPagesHandler
	pageHandler    queries.GetPageHandler
	ragClient      *ragclient.RAGClient
	baseURL        string // e.g. http://localhost:8090
}

func NewBookController(
	uploadHandler commands.UploadBookHandler,
	catalogHandler queries.GetBooksHandler,
	detailsHandler queries.GetBookByIDHandler,
	pagesHandler queries.GetBookPagesHandler,
	pageHandler queries.GetPageHandler,
	ragClient *ragclient.RAGClient,
	baseURL string,
) *BookController {
	return &BookController{
		uploadHandler:  uploadHandler,
		catalogHandler: catalogHandler,
		detailsHandler: detailsHandler,
		pagesHandler:   pagesHandler,
		pageHandler:    pageHandler,
		ragClient:      ragClient,
		baseURL:        baseURL,
	}
}