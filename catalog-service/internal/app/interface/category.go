package _interface

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, command *command.CreateCategoryCommand) (*entity.Category, error)
	UpdateCategory(ctx context.Context, command *command.UpdateCategoryCommand) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategoryByID(ctx context.Context, id int) (*entity.Category, error)
	GetAllCategories(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Category, int, error)
}
