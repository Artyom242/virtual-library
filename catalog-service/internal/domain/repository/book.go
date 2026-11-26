package repository

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	IRepository "virtual-library/catalog-service/internal/infra/db/repository"
)

type BookRepository interface {
	Create(c context.Context, book IRepository.IBook) (*entity.Book, error)
	GetByID(c context.Context, bookID int) (*entity.Book, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entity.Book, int, error)
	Delete(c context.Context, bookID int) error
	Update(c context.Context, book IRepository.IBook) (*entity.Book, error)
}
