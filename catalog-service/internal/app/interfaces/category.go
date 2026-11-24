package interfaces

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, command *command.CreateCategoryCommand) (*entities.Category, error)
	UpdateCategory(ctx context.Context, command *command.UpdateCategoryCommand) (*entities.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategoryByID(ctx context.Context, id int) (*entities.Category, error)
	GetAllCategories(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Category, int, error)
}
