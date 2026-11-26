package service

import (
	"context"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interface"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/domain/repository"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type CategoryService struct {
	tx   app.Transactor
	repo repository.CategoryRepository
}

func NewCategoryService(tx app.Transactor, repo repository.CategoryRepository) _interface.CategoryService {
	return &CategoryService{tx: tx, repo: repo}
}

func (svc *CategoryService) CreateCategory(ctx context.Context, command *command.CreateCategoryCommand) (*entity.Category, error) {
	category, err := entity.NewCategory(command.Name)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return svc.repo.Create(ctx, category)
}

func (svc *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return svc.repo.Delete(ctx, id)
}

func (svc *CategoryService) GetAllCategories(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Category, int, error) {
	return svc.repo.GetAll(ctx, pagination)
}

func (svc *CategoryService) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *CategoryService) UpdateCategory(ctx context.Context, command *command.UpdateCategoryCommand) (*entity.Category, error) {
	category, err := svc.repo.GetByID(ctx, command.ID)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	categoryChanges := entity.CategoryChanges{
		Name: command.Name,
	}

	category, err = category.ApplyCategory(categoryChanges)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return svc.repo.Update(ctx, category)
}
