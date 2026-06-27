package commands

import (
	"context"
	"fmt"
	"log/slog"

	"book_halal/internal/application/books/ports"
	repoBooks "book_halal/internal/domain/books"
	book "book_halal/internal/domain/books/entity"
	valueobjects "book_halal/internal/domain/books/value_objects"
	uuid "book_halal/internal/pkg"
)

// --- DTO ---

type UploadBookCommand struct {
	Title       string
	Author      string
	Description string
	Price       int
	TmpPDFPath  string
	// BaseURL is the public URL of the Go server (e.g. http://localhost:8090).
	// Used to build absolute image URLs sent to the RAG service.
	BaseURL string
}

type UploadBookResult struct {
	BookID string
}

// --- Interface ---

type UploadBookHandler interface {
	Handle(ctx context.Context, cmd UploadBookCommand) (*UploadBookResult, error)
}

// --- Handler ---

type uploadBookHandler struct {
	repo      repoBooks.Repository
	processor repoBooks.BookProcessor
	indexer   ports.RAGIndexer
	logger    *slog.Logger
}

func NewUploadBookHandler(
	repo repoBooks.Repository,
	processor repoBooks.BookProcessor,
	indexer ports.RAGIndexer,
	logger *slog.Logger,
) UploadBookHandler {
	return &uploadBookHandler{
		repo:      repo,
		processor: processor,
		indexer:   indexer,
		logger:    logger,
	}
}

func (h *uploadBookHandler) Handle(ctx context.Context, cmd UploadBookCommand) (*UploadBookResult, error) {
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

	go func() {
		pages, err := h.processor.Process(context.Background(), bookID.String(), cmd.TmpPDFPath)
		if err != nil {
			h.logger.Error("book processing failed", "book_id", bookID.String(), "err", err)
			return
		}

		if len(pages) == 0 {
			h.logger.Warn("no pages generated, skipping RAG indexing", "book_id", bookID.String())
			return
		}

		if err := h.indexer.IndexBook(bookID.String(), cmd.BaseURL, pages); err != nil {
			h.logger.Error("rag indexing failed", "book_id", bookID.String(), "err", err)
		}
	}()

	return &UploadBookResult{BookID: bookID.String()}, nil
}