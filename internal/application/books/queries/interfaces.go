package queries

import "context"

type GetBooksHandler interface {
	Handle(ctx context.Context) ([]BookCatalogView, error)
}

type GetBookByIDHandler interface {
	Handle(ctx context.Context, bookID string) (BookDetailsView, error)
}

type GetBookPagesHandler interface {
	Handle(ctx context.Context, bookID string) (BookPagesResponse, error)
}