package repository

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	IRepository "virtual-library/catalog-service/internal/infra/db/repository"
)

type CategoryRepository interface {
	Create(c context.Context, category IRepository.ICategory) (*entity.Category, error)
	GetByID(c context.Context, categoryID int) (*entity.Category, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entity.Category, int, error)
	Delete(c context.Context, categoryID int) error
	Update(c context.Context, category IRepository.ICategory) (*entity.Category, error)
}
