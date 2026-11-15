package catalog

import "errors"

type BookRepository interface {
	Create(book Book) (Book, error)
	GetByID(bookID int) (*Book, error)
	GetAll() ([]Book, error)
	Delete(bookID int) error
}

var ErrBookNotFound = errors.New("book not found")
