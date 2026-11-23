package repositories

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type BookRepository interface {
	Create(c context.Context, book *entities.Book) (*entities.Book, error)
	GetByID(c context.Context, bookID int) (*entities.Book, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entities.Book, int, error)
	Delete(c context.Context, bookID int) error
	Update(c context.Context, book *entities.Book) (*entities.Book, error)
}
