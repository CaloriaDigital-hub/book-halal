package handlers

import (
	"context"
	"fmt"

	uuid "book_halal/internal/pkg"

	"book_halal/internal/application/books/commands"
	repoBooks "book_halal/internal/domain/books"
	book "book_halal/internal/domain/books/entity"
	"book_halal/internal/domain/books/value_objects"
)

type UploadBookCommandHandler struct {
	repo      repoBooks.Repository
	processor repoBooks.BookProcessor
}

func NewUploadBookCommandHandler(repo repoBooks.Repository, processor repoBooks.BookProcessor) *UploadBookCommandHandler {
	return &UploadBookCommandHandler{
		repo:      repo,
		processor: processor,
	}
}

func (h *UploadBookCommandHandler) Handle(ctx context.Context, cmd commands.UploadBookCommand) (*commands.UploadBookResult, error) {
	
	
	priceVO, err := valueobjects.NewPriceFromMajor(cmd.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	bookID, err := valueobjects.NewBookId(uuid.New())
	if err != nil {
		return nil, fmt.Errorf("generate book id: %w", err)
	}

	
	newBook, err := book.New(bookID, cmd.Title, cmd.Author, cmd.Description, priceVO)
	if err != nil {
		return nil, fmt.Errorf("create book entity: %w", err)
	}

	
	if err := h.repo.Create(ctx, newBook); err != nil {
		return nil, fmt.Errorf("create book record: %w", err)
	}

	go h.processor.Process(context.Background(), bookID.String(), cmd.TmpPDFPath)

	return &commands.UploadBookResult{BookID: bookID.String()}, nil
}