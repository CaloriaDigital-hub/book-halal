package book
import (
	"errors"
	"context"
	"book_halal/internal/domain/books/entity"
)

var ErrBookNotFound = errors.New("book not found")

type Repository interface {
	Create(ctx context.Context, book *entity.Book) error
	UpdateStatus(ctx context.Context, bookID string, status entity.Status) error
	UpdateStatusWithPageCount(ctx context.Context, bookID string, status entity.Status, totalPages int) error
	SavePages(ctx context.Context, pages []entity.Page) error
	GetPagesByBookID(ctx context.Context, bookID string) ([]entity.Page, int, error)
	GetPagesByBookIDPaginated(ctx context.Context, bookID string, offset, limit int) ([]entity.Page, int, error)
	GetPageByNumber(ctx context.Context, bookID string, pageNumber int) (*entity.Page, error)
	UpdateCoverURL(ctx context.Context, bookID string, coverURL string) error
	GetAllReady(ctx context.Context) ([]entity.Book, error)
	GetByID(ctx context.Context, bookID string) (*entity.Book, error)
}