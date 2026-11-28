package service_test

import (
	"context"
	"errors"
	"testing"

	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/service"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	infraRepo "virtual-library/catalog-service/internal/infra/db/repository"
)

func TestCategoryService_CreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockCategoryRepo{
			createFn: func(_ context.Context, category infraRepo.ICategory) (*entity.Category, error) {
				return entity.HydrateCategory(1, category.Name()), nil
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.CreateCategory(context.Background(), &command.CreateCategoryCommand{Name: "Fantasy"})
		if err != nil {
			t.Fatalf("CreateCategory() unexpected error: %v", err)
		}
		if got == nil || got.Name() != "Fantasy" {
			t.Fatalf("CreateCategory() = %+v, want name Fantasy", got)
		}
		if len(repo.createCalls) != 1 {
			t.Fatalf("Create() calls = %d, want 1", len(repo.createCalls))
		}
	})

	t.Run("invalid name stops before repo", func(t *testing.T) {
		repo := &mockCategoryRepo{}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.CreateCategory(context.Background(), &command.CreateCategoryCommand{Name: "  "})
		if got != nil {
			t.Errorf("CreateCategory() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrValidationError) {
			t.Errorf("CreateCategory() error = %v, want ErrValidationError", err)
		}
		if len(repo.createCalls) != 0 {
			t.Errorf("Create() calls = %d, want 0", len(repo.createCalls))
		}
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := &mockCategoryRepo{
			createFn: func(_ context.Context, _ infraRepo.ICategory) (*entity.Category, error) {
				return nil, errBoom
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.CreateCategory(context.Background(), &command.CreateCategoryCommand{Name: "Fantasy"})
		if got != nil {
			t.Errorf("CreateCategory() = %+v, want nil", got)
		}
		if !errors.Is(err, errBoom) {
			t.Errorf("CreateCategory() error = %v, want errBoom", err)
		}
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	t.Run("delegates to repo", func(t *testing.T) {
		repo := &mockCategoryRepo{}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		if err := svc.DeleteCategory(context.Background(), 5); err != nil {
			t.Fatalf("DeleteCategory() unexpected error: %v", err)
		}
		if len(repo.deleteCalls) != 1 || repo.deleteCalls[0] != 5 {
			t.Errorf("Delete() calls = %v, want [5]", repo.deleteCalls)
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repo := &mockCategoryRepo{
			deleteFn: func(_ context.Context, _ int) error {
				return errBoom
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		if err := svc.DeleteCategory(context.Background(), 5); !errors.Is(err, errBoom) {
			t.Errorf("DeleteCategory() error = %v, want errBoom", err)
		}
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		repo := &mockCategoryRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Category, error) {
				return nil, entity.ErrCategoryNotFound
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.UpdateCategory(context.Background(), &command.UpdateCategoryCommand{ID: 5, Name: strPtr("Sci-Fi")})
		if got != nil {
			t.Errorf("UpdateCategory() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrCategoryNotFound) {
			t.Errorf("UpdateCategory() error = %v, want ErrCategoryNotFound", err)
		}
		if len(repo.updateCalls) != 0 {
			t.Errorf("Update() calls = %d, want 0", len(repo.updateCalls))
		}
	})

	t.Run("invalid name stops before update", func(t *testing.T) {
		repo := &mockCategoryRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Category, error) {
				return entity.HydrateCategory(5, "Fantasy"), nil
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.UpdateCategory(context.Background(), &command.UpdateCategoryCommand{ID: 5, Name: strPtr("  ")})
		if got != nil {
			t.Errorf("UpdateCategory() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrValidationError) {
			t.Errorf("UpdateCategory() error = %v, want ErrValidationError", err)
		}
		if len(repo.updateCalls) != 0 {
			t.Errorf("Update() calls = %d, want 0", len(repo.updateCalls))
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &mockCategoryRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Category, error) {
				return entity.HydrateCategory(5, "Fantasy"), nil
			},
			updateFn: func(_ context.Context, category infraRepo.ICategory) (*entity.Category, error) {
				return entity.HydrateCategory(5, category.Name()), nil
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.UpdateCategory(context.Background(), &command.UpdateCategoryCommand{ID: 5, Name: strPtr("Sci-Fi")})
		if err != nil {
			t.Fatalf("UpdateCategory() unexpected error: %v", err)
		}
		if got == nil || got.Name() != "Sci-Fi" {
			t.Fatalf("UpdateCategory() = %+v, want name Sci-Fi", got)
		}
		if len(repo.updateCalls) != 1 {
			t.Errorf("Update() calls = %d, want 1", len(repo.updateCalls))
		}
	})
}

func TestCategoryService_Delegation(t *testing.T) {
	t.Run("GetCategoryByID", func(t *testing.T) {
		want := entity.HydrateCategory(5, "Fantasy")
		repo := &mockCategoryRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Category, error) {
				return want, nil
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, err := svc.GetCategoryByID(context.Background(), 5)
		if err != nil {
			t.Fatalf("GetCategoryByID() unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("GetCategoryByID() = %+v, want %+v", got, want)
		}
	})

	t.Run("GetAllCategories", func(t *testing.T) {
		want := []*entity.Category{entity.HydrateCategory(5, "Fantasy")}
		repo := &mockCategoryRepo{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Category, int, error) {
				return want, 1, nil
			},
		}
		svc := service.NewCategoryService(&mockTransactor{}, repo)

		got, total, err := svc.GetAllCategories(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetAllCategories() unexpected error: %v", err)
		}
		if total != 1 {
			t.Errorf("GetAllCategories() total = %d, want 1", total)
		}
		if len(got) != 1 || got[0] != want[0] {
			t.Errorf("GetAllCategories() = %+v, want %+v", got, want)
		}
	})
}
