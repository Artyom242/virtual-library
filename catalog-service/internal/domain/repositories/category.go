package repositories

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type CategoryRepository interface {
	Create(c context.Context, category *entities.Category) (*entities.Category, error)
	GetByID(c context.Context, categoryID int) (*entities.Category, error)
	GetAll(c context.Context, pagination *pagination.Pagination) ([]*entities.Category, int, error)
	Delete(c context.Context, categoryID int) error
	Update(c context.Context, category *entities.Category) (*entities.Category, error)
}
