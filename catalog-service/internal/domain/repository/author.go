package repository

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	"virtual-library/catalog-service/internal/infra/db/repository"
)

type AuthorRepository interface {
	Create(c context.Context, author repository.IAuthor) (*entity.Author, error)
	GetByID(c context.Context, authorID int) (*entity.Author, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entity.Author, int, error)
	Delete(c context.Context, authorID int) error
	Update(c context.Context, author repository.IAuthor) (*entity.Author, error)
	GetOrCreate(ctx context.Context, author repository.IAuthor) (*entity.Author, error)
	HasBooks(c context.Context, authorID int) (bool, error)
}
