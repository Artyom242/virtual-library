package interfaces

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type BookService interface {
	CreateBook(ctx context.Context, command *command.CreateBookCommand) (*entities.Book, error)
	UpdateBook(ctx context.Context, command *command.UpdateBookCommand) (*entities.Book, error)
	DeleteBook(ctx context.Context, id int) error
	GetBookByID(ctx context.Context, id int) (*entities.Book, error)
	GetAllBooks(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Book, int, error)
	CreateBookWithNewAuthor(ctx context.Context, commandBook *command.CreateBookCommand, commandAuthor *command.CreateAuthorCommand) (*entities.Book, error)
}
