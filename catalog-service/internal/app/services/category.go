package services

import (
	"context"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/domain/repositories"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type CategoryService struct {
	tx   app.Transactor
	repo repositories.CategoryRepository
}

func NewCategoryService(tx app.Transactor, repo repositories.CategoryRepository) interfaces.CategoryService {
	return &CategoryService{tx: tx, repo: repo}
}

func (svc *CategoryService) CreateCategory(ctx context.Context, command *command.CreateCategoryCommand) (*entities.Category, error) {
	category, err := entities.NewCategory(command.Name)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return svc.repo.Create(ctx, category)
}

func (svc *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return svc.repo.Delete(ctx, id)
}

func (svc *CategoryService) GetAllCategories(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Category, int, error) {
	return svc.repo.GetAll(ctx, pagination)
}

func (svc *CategoryService) GetCategoryByID(ctx context.Context, id int) (*entities.Category, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *CategoryService) UpdateCategory(ctx context.Context, command *command.UpdateCategoryCommand) (*entities.Category, error) {
	category, err := svc.repo.GetByID(ctx, command.ID)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	categoryChanges := entities.CategoryChanges{
		Name: command.Name,
	}

	category, err = category.ApplyCategory(categoryChanges)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return svc.repo.Update(ctx, category)
}
