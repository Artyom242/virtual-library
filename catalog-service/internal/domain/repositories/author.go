package repositories

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type AuthorRepository interface {
	Create(c context.Context, author *entities.Author) (*entities.Author, error)
	GetByID(c context.Context, authorID int) (*entities.Author, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entities.Author, int, error)
	Delete(c context.Context, authorID int) error
	Update(c context.Context, author *entities.Author) (*entities.Author, error)
	GetOrCreate(ctx context.Context, author *entities.Author) (*entities.Author, error)
	HasBooks(c context.Context, authorID int) (bool, error)
}
