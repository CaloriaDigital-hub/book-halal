package books

import (
	"book_halal/internal/application/books/commands"
	"book_halal/internal/application/books/queries"
	purchaseQueries "book_halal/internal/application/purchases/queries"
	"book_halal/internal/infrastructure/ragclient"
)

type BookController struct {
	uploadHandler  commands.UploadBookHandler
	catalogHandler queries.GetBooksHandler
	detailsHandler queries.GetBookByIDHandler
	pagesHandler   queries.GetBookPagesHandler
	pageHandler    queries.GetPageHandler
	ragClient      *ragclient.RAGClient
	access         purchaseQueries.AccessChecker
	baseURL        string // e.g. http://localhost:8090
}

func NewBookController(
	uploadHandler commands.UploadBookHandler,
	catalogHandler queries.GetBooksHandler,
	detailsHandler queries.GetBookByIDHandler,
	pagesHandler queries.GetBookPagesHandler,
	pageHandler queries.GetPageHandler,
	ragClient *ragclient.RAGClient,
	access purchaseQueries.AccessChecker,
	baseURL string,
) *BookController {
	return &BookController{
		uploadHandler:  uploadHandler,
		catalogHandler: catalogHandler,
		detailsHandler: detailsHandler,
		pagesHandler:   pagesHandler,
		pageHandler:    pageHandler,
		ragClient:      ragClient,
		access:         access,
		baseURL:        baseURL,
	}
}