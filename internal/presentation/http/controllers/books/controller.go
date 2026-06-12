package books

import (
	"book_halal/internal/application/books/commands"
	"book_halal/internal/application/books/queries"
)

type BookController struct {
	uploadHandler  commands.UploadBookHandler
	catalogHandler queries.GetBooksHandler
	detailsHandler queries.GetBookByIDHandler
	pagesHandler   queries.GetBookPagesHandler
}

func NewBookController(
	uploadHandler commands.UploadBookHandler,
	catalogHandler queries.GetBooksHandler,
	detailsHandler queries.GetBookByIDHandler,
	pagesHandler queries.GetBookPagesHandler,
) *BookController {
	return &BookController{
		uploadHandler:  uploadHandler,
		catalogHandler: catalogHandler,
		detailsHandler: detailsHandler,
		pagesHandler:   pagesHandler,
	}
}