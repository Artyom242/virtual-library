package _interface

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type BookService interface {
	CreateBook(ctx context.Context, command *command.CreateBookCommand) (*entity.Book, error)
	UpdateBook(ctx context.Context, command *command.UpdateBookCommand) (*entity.Book, error)
	DeleteBook(ctx context.Context, id int) error
	GetBookByID(ctx context.Context, id int) (*entity.Book, error)
	GetAllBooks(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Book, int, error)
	CreateBookWithNewAuthor(ctx context.Context, commandBook *command.CreateBookCommand, commandAuthor *command.CreateAuthorCommand) (*entity.Book, error)
}
