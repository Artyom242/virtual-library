package domain

import (
	"context"
)

type BookRepository interface {
	Create(c context.Context, book Book) (Book, error)
	GetByID(c context.Context, bookID int) (*Book, error)
	GetAll(c context.Context, page, limit int) ([]Book, int, error)
	Delete(c context.Context, bookID int) error
	Update(c context.Context, id int, book BookPatch) (Book, error)
}
