package queries

import (
	"context"

	repoBooks "book_halal/internal/domain/books"
	domain "book_halal/internal/domain/purchases"
	"book_halal/internal/domain/users/entity"
)

// AccessChecker решает, может ли пользователь читать книгу и задавать по ней вопросы.
type AccessChecker interface {
	// CanRead: админ — всегда; бесплатная книга — всем; платная — после подтверждённой оплаты.
	// Возвращает ошибку книги (например, ErrBookNotFound), если её нет.
	CanRead(ctx context.Context, user *entity.User, bookID string) (bool, error)
}

type accessChecker struct {
	repo  domain.Repository
	books repoBooks.Repository
}

func NewAccessChecker(repo domain.Repository, books repoBooks.Repository) AccessChecker {
	return &accessChecker{repo: repo, books: books}
}

func (a *accessChecker) CanRead(ctx context.Context, user *entity.User, bookID string) (bool, error) {
	book, err := a.books.GetByID(ctx, bookID)
	if err != nil {
		return false, err
	}
	if user.Role.IsAdmin() || book.Price.IsFree() {
		return true, nil
	}
	return a.repo.HasApproved(ctx, user.ID.String(), bookID)
}
