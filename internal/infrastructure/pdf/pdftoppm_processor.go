package pdf


import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	bookRepo "book_halal/internal/domain/books"
	"book_halal/internal/domain/books/entity"
	uuid "book_halal/internal/pkg" // Твой пакет с UUID
)

// Убеждаемся, что структура реализует интерфейс
var _ bookRepo.BookProcessor = (*PDFToPPMProcessor)(nil)

type PDFToPPMProcessor struct {
	repo       bookRepo.Repository
	staticRoot string
	logger     *slog.Logger
}

func NewPDFToPPMProcessor(repo bookRepo.Repository, staticRoot string, logger *slog.Logger) *PDFToPPMProcessor {
	return &PDFToPPMProcessor{
		repo:       repo,
		staticRoot: staticRoot,
		logger:     logger,
	}
}

func (p *PDFToPPMProcessor) Process(ctx context.Context, bookID string, pdfPath string) {
	logger := p.logger.With("book_id", bookID)

	defer func() {
		if err := os.Remove(pdfPath); err != nil {
			logger.Warn("failed to remove tmp pdf", "path", pdfPath, "err", err)
		}
	}()

	outDir := filepath.Join(p.staticRoot, "books", bookID)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		p.markError(ctx, bookID, logger, fmt.Errorf("mkdir: %w", err))
		return
	}

	outputPrefix := filepath.Join(outDir, "page")
	cmd := exec.CommandContext(ctx, "pdftoppm", "-jpeg", "-r", "150", pdfPath, outputPrefix)

	if out, err := cmd.CombinedOutput(); err != nil {
		p.markError(ctx, bookID, logger, fmt.Errorf("pdftoppm: %w, output: %s", err, string(out)))
		return
	}

	pages, err := p.collectPages(bookID, outDir)
	if err != nil {
		p.markError(ctx, bookID, logger, err)
		return
	}

	if len(pages) == 0 {
		p.markError(ctx, bookID, logger, fmt.Errorf("no pages generated"))
		return
	}

	if err := p.repo.SavePages(ctx, pages); err != nil {
		p.markError(ctx, bookID, logger, fmt.Errorf("save pages: %w", err))
		return
	}

	coverURL := pages[0].ImageURL
	if err := p.repo.UpdateCoverURL(ctx, bookID, coverURL); err != nil {
    logger.Error("failed to save cover url", "err", err)
    // не фатально, продолжаем
	}

	if err := p.repo.UpdateStatusWithPageCount(ctx, bookID, entity.StatusReady, len(pages)); err != nil {
		logger.Error("failed to mark book ready", "err", err)
		return
	}

	logger.Info("book processing complete", "total_pages", len(pages))
}

func (p *PDFToPPMProcessor) collectPages(bookID, outDir string) ([]entity.Page, error) {
	entries, err := os.ReadDir(outDir)
	if err != nil {
		return nil, fmt.Errorf("read output dir: %w", err)
	}

	var pages []entity.Page
	pageNum := 1

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") {
			continue
		}

		imageURL := fmt.Sprintf("/static/books/%s/%s", bookID, name)

		pages = append(pages, entity.Page{
			ID:         uuid.New(),
			BookID:     bookID,
			PageNumber: pageNum,
			ImageURL:   imageURL,
		})
		pageNum++
	}

	return pages, nil
}

func (p *PDFToPPMProcessor) markError(ctx context.Context, bookID string, logger *slog.Logger, err error) {
	logger.Error("book processing failed", "err", err)
	if updateErr := p.repo.UpdateStatus(ctx, bookID, entity.StatusError); updateErr != nil {
		logger.Error("failed to mark book as error", "err", updateErr)
	}
}