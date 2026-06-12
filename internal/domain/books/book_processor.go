package book

import "context"


type BookProcessor interface {
	Process(ctx context.Context, bookID string, pdfPath string)
}