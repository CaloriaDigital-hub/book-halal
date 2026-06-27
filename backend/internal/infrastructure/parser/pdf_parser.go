package parser


import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type PDFParser struct {
	outputDir string // Папка, куда будем сохранять картинки (например, "./static/books")
}

func NewPDFParser(outputDir string) *PDFParser {
	// Создаем папку, если ее нет
	os.MkdirAll(outputDir, 0755)
	return &PDFParser{outputDir: outputDir}
}

// Parse принимает путь к временно сохраненному PDF файлу и возвращает пути к картинкам страниц
func (p *PDFParser) Parse(ctx context.Context, pdfFilePath string, bookID string) ([]string, error) {
	// Папка для конкретной книги: ./static/books/{bookID}
	bookDir := filepath.Join(p.outputDir, bookID)
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create book directory: %w", err)
	}

	// Команда: pdftoppm -jpeg {input.pdf} {output_prefix}
	// Она создаст файлы вида page-1.jpg, page-2.jpg и т.д.
	outputPrefix := filepath.Join(bookDir, "page")
	cmd := exec.CommandContext(ctx, "pdftoppm", "-jpeg", pdfFilePath, outputPrefix)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to extract images from pdf: %w", err)
	}

	// Читаем, какие файлы создались
	files, err := os.ReadDir(bookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read book directory: %w", err)
	}

	var imageUrls []string
	for _, f := range files {
		if !f.IsDir() {
			// Формируем относительный URL для базы данных, например: /static/books/123/page-01.jpg
			imageUrls = append(imageUrls, fmt.Sprintf("/static/books/%s/%s", bookID, f.Name()))
		}
	}

	// pdftoppm добавляет нули (page-01, page-02), поэтому обычная сортировка строк сработает идеально
	sort.Strings(imageUrls)

	return imageUrls, nil
}