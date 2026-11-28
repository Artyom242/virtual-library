package service_test

import (
	"context"
	"errors"
	"testing"

	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/service"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	infraRepo "virtual-library/catalog-service/internal/infra/db/repository"
)

var errBoom = errors.New("boom")

func TestAuthorService_CreateAuthor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockAuthorRepo{
			getOrCreateFn: func(_ context.Context, author infraRepo.IAuthor) (*entity.Author, error) {
				return entity.HydrateAuthor(1, author.Name()), nil
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.CreateAuthor(context.Background(), &command.CreateAuthorCommand{Name: "John"})
		if err != nil {
			t.Fatalf("CreateAuthor() unexpected error: %v", err)
		}
		if got == nil || got.Name() != "John" {
			t.Fatalf("CreateAuthor() = %+v, want name John", got)
		}
		if len(repo.getOrCreateCalls) != 1 {
			t.Fatalf("GetOrCreate() calls = %d, want 1", len(repo.getOrCreateCalls))
		}
		if repo.getOrCreateCalls[0].Name() != "John" {
			t.Errorf("GetOrCreate() name = %q, want John", repo.getOrCreateCalls[0].Name())
		}
	})

	t.Run("invalid name stops before repo", func(t *testing.T) {
		repo := &mockAuthorRepo{}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.CreateAuthor(context.Background(), &command.CreateAuthorCommand{Name: "  "})
		if err == nil {
			t.Fatal("CreateAuthor() expected error, got nil")
		}
		if got != nil {
			t.Errorf("CreateAuthor() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrValidationError) {
			t.Errorf("CreateAuthor() error = %v, want ErrValidationError", err)
		}
		if len(repo.getOrCreateCalls) != 0 {
			t.Errorf("GetOrCreate() calls = %d, want 0", len(repo.getOrCreateCalls))
		}
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := &mockAuthorRepo{
			getOrCreateFn: func(_ context.Context, _ infraRepo.IAuthor) (*entity.Author, error) {
				return nil, errBoom
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.CreateAuthor(context.Background(), &command.CreateAuthorCommand{Name: "John"})
		if got != nil {
			t.Errorf("CreateAuthor() = %+v, want nil", got)
		}
		if !errors.Is(err, errBoom) {
			t.Errorf("CreateAuthor() error = %v, want errBoom", err)
		}
	})
}

func TestAuthorService_DeleteAuthor(t *testing.T) {
	tests := []struct {
		name        string
		hasBooks    bool
		hasBooksErr error
		deleteErr   error
		wantErr     error
		wantDelete  bool
	}{
		{
			name:       "has books returns conflict",
			hasBooks:   true,
			wantErr:    app.ErrAuthorHasBooks,
			wantDelete: false,
		},
		{
			name:        "has books check error",
			hasBooksErr: errBoom,
			wantErr:     errBoom,
			wantDelete:  false,
		},
		{
			name:       "success deletes",
			hasBooks:   false,
			wantDelete: true,
		},
		{
			name:       "delete error propagated",
			hasBooks:   false,
			deleteErr:  errBoom,
			wantErr:    errBoom,
			wantDelete: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAuthorRepo{
				hasBooksFn: func(_ context.Context, _ int) (bool, error) {
					return tt.hasBooks, tt.hasBooksErr
				},
				deleteFn: func(_ context.Context, _ int) error {
					return tt.deleteErr
				},
			}
			svc := service.NewAuthorService(&mockTransactor{}, repo)

			err := svc.DeleteAuthor(context.Background(), 5)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("DeleteAuthor() error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("DeleteAuthor() unexpected error: %v", err)
			}

			if tt.wantDelete && len(repo.deleteCalls) != 1 {
				t.Errorf("Delete() calls = %d, want 1", len(repo.deleteCalls))
			}
			if !tt.wantDelete && len(repo.deleteCalls) != 0 {
				t.Errorf("Delete() calls = %d, want 0", len(repo.deleteCalls))
			}
		})
	}
}

func TestAuthorService_UpdateAuthor(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		repo := &mockAuthorRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
				return nil, entity.ErrAuthorNotFound
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.UpdateAuthor(context.Background(), &command.UpdateAuthorCommand{ID: 5, Name: strPtr("Jane")})
		if got != nil {
			t.Errorf("UpdateAuthor() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrAuthorNotFound) {
			t.Errorf("UpdateAuthor() error = %v, want ErrAuthorNotFound", err)
		}
		if len(repo.updateCalls) != 0 {
			t.Errorf("Update() calls = %d, want 0", len(repo.updateCalls))
		}
	})

	t.Run("invalid name stops before update", func(t *testing.T) {
		repo := &mockAuthorRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
				return entity.HydrateAuthor(5, "John"), nil
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.UpdateAuthor(context.Background(), &command.UpdateAuthorCommand{ID: 5, Name: strPtr("  ")})
		if got != nil {
			t.Errorf("UpdateAuthor() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrValidationError) {
			t.Errorf("UpdateAuthor() error = %v, want ErrValidationError", err)
		}
		if len(repo.updateCalls) != 0 {
			t.Errorf("Update() calls = %d, want 0", len(repo.updateCalls))
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &mockAuthorRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
				return entity.HydrateAuthor(5, "John"), nil
			},
			updateFn: func(_ context.Context, author infraRepo.IAuthor) (*entity.Author, error) {
				return entity.HydrateAuthor(5, author.Name()), nil
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.UpdateAuthor(context.Background(), &command.UpdateAuthorCommand{ID: 5, Name: strPtr("Jane")})
		if err != nil {
			t.Fatalf("UpdateAuthor() unexpected error: %v", err)
		}
		if got == nil || got.Name() != "Jane" {
			t.Fatalf("UpdateAuthor() = %+v, want name Jane", got)
		}
		if len(repo.updateCalls) != 1 {
			t.Errorf("Update() calls = %d, want 1", len(repo.updateCalls))
		}
	})
}

func TestAuthorService_Delegation(t *testing.T) {
	t.Run("GetAuthorByID", func(t *testing.T) {
		want := entity.HydrateAuthor(5, "John")
		repo := &mockAuthorRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
				return want, nil
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, err := svc.GetAuthorByID(context.Background(), 5)
		if err != nil {
			t.Fatalf("GetAuthorByID() unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("GetAuthorByID() = %+v, want %+v", got, want)
		}
	})

	t.Run("GetAllAuthors", func(t *testing.T) {
		want := []*entity.Author{entity.HydrateAuthor(5, "John")}
		repo := &mockAuthorRepo{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Author, int, error) {
				return want, 1, nil
			},
		}
		svc := service.NewAuthorService(&mockTransactor{}, repo)

		got, total, err := svc.GetAllAuthors(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetAllAuthors() unexpected error: %v", err)
		}
		if total != 1 {
			t.Errorf("GetAllAuthors() total = %d, want 1", total)
		}
		if len(got) != 1 || got[0] != want[0] {
			t.Errorf("GetAllAuthors() = %+v, want %+v", got, want)
		}
	})
}
