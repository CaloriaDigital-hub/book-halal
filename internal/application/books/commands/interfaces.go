package commands



import "context"

type UploadBookHandler interface {
	Handle(ctx context.Context, cmd UploadBookCommand) (*UploadBookResult, error)
}