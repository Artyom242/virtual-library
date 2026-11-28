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

func TestBookService_CreateBook(t *testing.T) {
	tests := []struct {
		name       string
		command    *command.CreateBookCommand
		author     *entity.Author
		authorErr  error
		createBook *entity.Book
		createErr  error
		wantErr    error
		wantCreate bool
		wantAuthor bool
	}{
		{
			name:       "invalid author id",
			command:    &command.CreateBookCommand{Title: "Go", Year: 2020, AuthorID: 0},
			wantErr:    entity.ErrInvalidAuthorID,
			wantAuthor: false,
		},
		{
			name:       "author not found",
			command:    &command.CreateBookCommand{Title: "Go", Year: 2020, AuthorID: 2},
			authorErr:  entity.ErrAuthorNotFound,
			wantErr:    entity.ErrAuthorNotFound,
			wantAuthor: true,
		},
		{
			name:       "author lookup error",
			command:    &command.CreateBookCommand{Title: "Go", Year: 2020, AuthorID: 2},
			authorErr:  errBoom,
			wantErr:    errBoom,
			wantAuthor: true,
		},
		{
			name:       "invalid book data",
			command:    &command.CreateBookCommand{Title: "  ", Year: 2020, AuthorID: 2},
			author:     entity.HydrateAuthor(2, "John"),
			wantErr:    entity.ErrValidationError,
			wantAuthor: true,
		},
		{
			name:       "repo create error",
			command:    &command.CreateBookCommand{Title: "Go", Year: 2020, AuthorID: 2},
			author:     entity.HydrateAuthor(2, "John"),
			createErr:  errBoom,
			wantErr:    errBoom,
			wantAuthor: true,
			wantCreate: true,
		},
		{
			name:       "success",
			command:    &command.CreateBookCommand{Title: "Go", Year: 2020, AuthorID: 2},
			author:     entity.HydrateAuthor(2, "John"),
			createBook: entity.HydrateBook(10, "Go", 2020, 2, nil, nil),
			wantAuthor: true,
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mockBookRepo{
				createFn: func(_ context.Context, _ infraRepo.IBook) (*entity.Book, error) {
					return tt.createBook, tt.createErr
				},
			}
			authorRepo := &mockAuthorRepo{
				getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
					return tt.author, tt.authorErr
				},
			}
			svc := service.NewService(&mockTransactor{}, bookRepo, authorRepo)

			got, err := svc.CreateBook(context.Background(), tt.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("CreateBook() error = %v, want %v", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("CreateBook() = %+v, want nil", got)
				}
			} else {
				if err != nil {
					t.Fatalf("CreateBook() unexpected error: %v", err)
				}
				if got != tt.createBook {
					t.Errorf("CreateBook() = %+v, want %+v", got, tt.createBook)
				}
			}

			if tt.wantAuthor && len(authorRepo.getByIDCalls) != 1 {
				t.Errorf("GetByID() author calls = %d, want 1", len(authorRepo.getByIDCalls))
			}
			if !tt.wantAuthor && len(authorRepo.getByIDCalls) != 0 {
				t.Errorf("GetByID() author calls = %d, want 0", len(authorRepo.getByIDCalls))
			}
			if tt.wantCreate && len(bookRepo.createCalls) != 1 {
				t.Errorf("Create() book calls = %d, want 1", len(bookRepo.createCalls))
			}
			if !tt.wantCreate && len(bookRepo.createCalls) != 0 {
				t.Errorf("Create() book calls = %d, want 0", len(bookRepo.createCalls))
			}
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	t.Run("book not found", func(t *testing.T) {
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return nil, entity.ErrBookNotFound
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5, Title: strPtr("New")})
		if got != nil {
			t.Errorf("UpdateBook() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrBookNotFound) {
			t.Errorf("UpdateBook() error = %v, want ErrBookNotFound", err)
		}
	})

	t.Run("invalid author id", func(t *testing.T) {
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return entity.HydrateBook(5, "Go", 2020, 2, nil, nil), nil
			},
		}
		authorRepo := &mockAuthorRepo{}
		svc := service.NewService(&mockTransactor{}, bookRepo, authorRepo)

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5, AuthorID: intPtr(0)})
		if got != nil {
			t.Errorf("UpdateBook() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrInvalidAuthorID) {
			t.Errorf("UpdateBook() error = %v, want ErrInvalidAuthorID", err)
		}
		if len(authorRepo.getByIDCalls) != 0 {
			t.Errorf("GetByID() author calls = %d, want 0", len(authorRepo.getByIDCalls))
		}
		if len(bookRepo.updateCalls) != 0 {
			t.Errorf("Update() book calls = %d, want 0", len(bookRepo.updateCalls))
		}
	})

	t.Run("author not found", func(t *testing.T) {
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return entity.HydrateBook(5, "Go", 2020, 2, nil, nil), nil
			},
		}
		authorRepo := &mockAuthorRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
				return nil, entity.ErrAuthorNotFound
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, authorRepo)

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5, AuthorID: intPtr(2)})
		if got != nil {
			t.Errorf("UpdateBook() = %+v, want nil", got)
		}
		if !errors.Is(err, entity.ErrAuthorNotFound) {
			t.Errorf("UpdateBook() error = %v, want ErrAuthorNotFound", err)
		}
	})

	t.Run("empty patch", func(t *testing.T) {
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return entity.HydrateBook(5, "Go", 2020, 2, nil, nil), nil
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5})
		if got != nil {
			t.Errorf("UpdateBook() = %+v, want nil", got)
		}
		if !errors.Is(err, app.ErrEmptyBookPatch) {
			t.Errorf("UpdateBook() error = %v, want ErrEmptyBookPatch", err)
		}
		if len(bookRepo.updateCalls) != 0 {
			t.Errorf("Update() book calls = %d, want 0", len(bookRepo.updateCalls))
		}
	})

	t.Run("update error", func(t *testing.T) {
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return entity.HydrateBook(5, "Go", 2020, 2, nil, nil), nil
			},
			updateFn: func(_ context.Context, _ infraRepo.IBook) (*entity.Book, error) {
				return nil, errBoom
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5, Title: strPtr("New")})
		if got != nil {
			t.Errorf("UpdateBook() = %+v, want nil", got)
		}
		if !errors.Is(err, errBoom) {
			t.Errorf("UpdateBook() error = %v, want errBoom", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		updated := entity.HydrateBook(5, "New", 2020, 2, nil, nil)
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return entity.HydrateBook(5, "Go", 2020, 2, nil, nil), nil
			},
			updateFn: func(_ context.Context, _ infraRepo.IBook) (*entity.Book, error) {
				return updated, nil
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, err := svc.UpdateBook(context.Background(), &command.UpdateBookCommand{ID: 5, Title: strPtr("New")})
		if err != nil {
			t.Fatalf("UpdateBook() unexpected error: %v", err)
		}
		if got != updated {
			t.Errorf("UpdateBook() = %+v, want %+v", got, updated)
		}
		if len(bookRepo.updateCalls) != 1 {
			t.Errorf("Update() book calls = %d, want 1", len(bookRepo.updateCalls))
		}
	})
}

func TestBookService_CreateBookWithNewAuthor(t *testing.T) {
	tests := []struct {
		name            string
		bookCommand     *command.CreateBookCommand
		authorCommand   *command.CreateAuthorCommand
		createdAuthor   *entity.Author
		getOrCreateErr  error
		createdBook     *entity.Book
		createErr       error
		wantErr         error
		wantGetOrCreate bool
		wantCreate      bool
	}{
		{
			name:            "invalid author name",
			bookCommand:     &command.CreateBookCommand{Title: "Go", Year: 2020},
			authorCommand:   &command.CreateAuthorCommand{Name: "  "},
			wantErr:         entity.ErrValidationError,
			wantGetOrCreate: false,
		},
		{
			name:            "get or create error",
			bookCommand:     &command.CreateBookCommand{Title: "Go", Year: 2020},
			authorCommand:   &command.CreateAuthorCommand{Name: "New"},
			getOrCreateErr:  errBoom,
			wantErr:         errBoom,
			wantGetOrCreate: true,
		},
		{
			name:            "invalid book data",
			bookCommand:     &command.CreateBookCommand{Title: "  ", Year: 2020},
			authorCommand:   &command.CreateAuthorCommand{Name: "New"},
			createdAuthor:   entity.HydrateAuthor(5, "New"),
			wantErr:         entity.ErrValidationError,
			wantGetOrCreate: true,
		},
		{
			name:            "create error",
			bookCommand:     &command.CreateBookCommand{Title: "Go", Year: 2020},
			authorCommand:   &command.CreateAuthorCommand{Name: "New"},
			createdAuthor:   entity.HydrateAuthor(5, "New"),
			createErr:       errBoom,
			wantErr:         errBoom,
			wantGetOrCreate: true,
			wantCreate:      true,
		},
		{
			name:            "success",
			bookCommand:     &command.CreateBookCommand{Title: "Go", Year: 2020},
			authorCommand:   &command.CreateAuthorCommand{Name: "New"},
			createdAuthor:   entity.HydrateAuthor(5, "New"),
			createdBook:     entity.HydrateBook(10, "Go", 2020, 5, nil, nil),
			wantGetOrCreate: true,
			wantCreate:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mockBookRepo{
				createFn: func(_ context.Context, _ infraRepo.IBook) (*entity.Book, error) {
					return tt.createdBook, tt.createErr
				},
			}
			authorRepo := &mockAuthorRepo{
				getOrCreateFn: func(_ context.Context, _ infraRepo.IAuthor) (*entity.Author, error) {
					return tt.createdAuthor, tt.getOrCreateErr
				},
			}
			svc := service.NewService(&mockTransactor{}, bookRepo, authorRepo)

			got, err := svc.CreateBookWithNewAuthor(context.Background(), tt.bookCommand, tt.authorCommand)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("CreateBookWithNewAuthor() error = %v, want %v", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("CreateBookWithNewAuthor() = %+v, want nil", got)
				}
			} else {
				if err != nil {
					t.Fatalf("CreateBookWithNewAuthor() unexpected error: %v", err)
				}
				if got != tt.createdBook {
					t.Errorf("CreateBookWithNewAuthor() = %+v, want %+v", got, tt.createdBook)
				}
			}

			if tt.wantGetOrCreate && len(authorRepo.getOrCreateCalls) != 1 {
				t.Errorf("GetOrCreate() calls = %d, want 1", len(authorRepo.getOrCreateCalls))
			}
			if !tt.wantGetOrCreate && len(authorRepo.getOrCreateCalls) != 0 {
				t.Errorf("GetOrCreate() calls = %d, want 0", len(authorRepo.getOrCreateCalls))
			}
			if tt.wantCreate && len(bookRepo.createCalls) != 1 {
				t.Errorf("Create() book calls = %d, want 1", len(bookRepo.createCalls))
			}
			if !tt.wantCreate && len(bookRepo.createCalls) != 0 {
				t.Errorf("Create() book calls = %d, want 0", len(bookRepo.createCalls))
			}
		})
	}
}

func TestBookService_Delegation(t *testing.T) {
	t.Run("GetBookByID", func(t *testing.T) {
		want := entity.HydrateBook(5, "Go", 2020, 2, nil, nil)
		bookRepo := &mockBookRepo{
			getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
				return want, nil
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, err := svc.GetBookByID(context.Background(), 5)
		if err != nil {
			t.Fatalf("GetBookByID() unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("GetBookByID() = %+v, want %+v", got, want)
		}
	})

	t.Run("GetAllBooks", func(t *testing.T) {
		want := []*entity.Book{entity.HydrateBook(5, "Go", 2020, 2, nil, nil)}
		bookRepo := &mockBookRepo{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Book, int, error) {
				return want, 1, nil
			},
		}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		got, total, err := svc.GetAllBooks(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetAllBooks() unexpected error: %v", err)
		}
		if total != 1 {
			t.Errorf("GetAllBooks() total = %d, want 1", total)
		}
		if len(got) != 1 || got[0] != want[0] {
			t.Errorf("GetAllBooks() = %+v, want %+v", got, want)
		}
	})

	t.Run("DeleteBook", func(t *testing.T) {
		bookRepo := &mockBookRepo{}
		svc := service.NewService(&mockTransactor{}, bookRepo, &mockAuthorRepo{})

		if err := svc.DeleteBook(context.Background(), 5); err != nil {
			t.Fatalf("DeleteBook() unexpected error: %v", err)
		}
		if len(bookRepo.deleteCalls) != 1 || bookRepo.deleteCalls[0] != 5 {
			t.Errorf("Delete() calls = %v, want [5]", bookRepo.deleteCalls)
		}
	})
}
