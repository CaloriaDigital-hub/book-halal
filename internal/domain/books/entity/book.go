package entity

import (
	"book_halal/internal/domain/books/value_objects"
	"errors"
	"time"
)

type Status string

const (
	StatusProcessing Status = "processing"
	StatusReady      Status = "ready"
	StatusError      Status = "error"
)

type Book struct {
	ID          valueobjects.BookId
	Title       string
	Author      string
	Description string
	Price       valueobjects.Price
	CoverURL	string   
	TotalPages  int
	Status      Status
	Pages       []Page
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


func New(id valueobjects.BookId, title, author, description string, price valueobjects.Price) (*Book, error){
	if title == "" {
		return nil, errors.New("book title cannot be empty")
	}
	if author == "" {
		return nil, errors.New("book author cannot be empty")
	}
	if description == "" {
		return nil, errors.New("book description cannot be empty")
	}

	return &Book{
		ID: id,
		Title: title,
		Author: author,
		Description: description,
		Price: price,
		Status: StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	
	}, nil

}

func (b *Book) AddPages(pages []Page) {
	b.Pages = pages
	b.TotalPages = len(pages)
	b.UpdatedAt = time.Now()


}

func (b *Book) SetCover(url string) {
	b.CoverURL = url
	b.UpdatedAt = time.Now()
}