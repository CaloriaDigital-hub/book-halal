package book
import (
	"context"
	"book_halal/internal/domain/books/entity"
)



type Repository interface {
	Create(ctx context.Context, book *entity.Book) error
	UpdateStatus(ctx context.Context, bookID string, status entity.Status) error
	UpdateStatusWithPageCount(ctx context.Context, bookID string, status entity.Status, totalPages int) error
	SavePages(ctx context.Context, pages []entity.Page) error
	GetPagesByBookID(ctx context.Context, bookID string) ([]entity.Page, int, error)
	UpdateCoverURL(ctx context.Context, bookID string, coverURL string) error
	GetAllReady(ctx context.Context) ([]entity.Book, error)
	GetByID(ctx context.Context, bookID string) (*entity.Book, error)
	
}