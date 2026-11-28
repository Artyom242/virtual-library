package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/httpapi/handler"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

var errBoom = errors.New("boom")

func TestBookHandler_GetBooks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockBookService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Book, int, error) {
				return []*entity.Book{
					entity.HydrateBook(1, "Go", 2020, 2, nil, nil),
					entity.HydrateBook(2, "Rust", 2021, 2, nil, nil),
				}, 2, nil
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodGet, "/books", "")

		h.GetBooks(c, generated.GetBooksParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body=%q)", w.Code, w.Body.String())
		}
		var resp generated.GetBooks200JSONResponse
		decodeJSON(t, w, &resp)
		if len(resp.Data) != 2 || resp.Total != 2 {
			t.Errorf("response = %+v, want 2 items total 2", resp)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &mockBookService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Book, int, error) {
				return []*entity.Book{}, 0, nil
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodGet, "/books", "")

		h.GetBooks(c, generated.GetBooksParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp generated.GetBooks200JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data == nil || len(resp.Data) != 0 {
			t.Errorf("data = %v, want empty non-nil slice", resp.Data)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockBookService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Book, int, error) {
				return nil, 0, errBoom
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodGet, "/books", "")

		h.GetBooks(c, generated.GetBooksParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
		var resp generated.ErrorResponse
		decodeJSON(t, w, &resp)
		if resp.Message != "Internal server error" {
			t.Errorf("message = %q, want %q", resp.Message, "Internal server error")
		}
	})
}

func TestBookHandler_GetBookById(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		book       *entity.Book
		serviceErr error
		wantStatus int
		wantMsg    string
		wantCall   bool
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantMsg: "invalid book id", wantCall: false},
		{name: "not found", id: 5, serviceErr: entity.ErrBookNotFound, wantStatus: http.StatusNotFound, wantMsg: "book not found", wantCall: true},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantMsg: "Internal server error", wantCall: true},
		{name: "success", id: 5, book: entity.HydrateBook(5, "Go", 2020, 2, nil, nil), wantStatus: http.StatusOK, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockBookService{
				getByIDFn: func(_ context.Context, _ int) (*entity.Book, error) {
					return tt.book, tt.serviceErr
				},
			}
			h := handler.NewBookHandler(svc)
			c, w := newTestContext(http.MethodGet, "/books/5", "")

			h.GetBookById(c, tt.id)
			c.Writer.WriteHeaderNow()

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", w.Code, tt.wantStatus, w.Body.String())
			}
			if svc.getByIDCalled != tt.wantCall {
				t.Errorf("GetBookByID called = %v, want %v", svc.getByIDCalled, tt.wantCall)
			}
			if tt.wantMsg != "" {
				var resp generated.ErrorResponse
				decodeJSON(t, w, &resp)
				if resp.Message != tt.wantMsg {
					t.Errorf("message = %q, want %q", resp.Message, tt.wantMsg)
				}
			}
			if tt.wantStatus == http.StatusOK {
				var resp generated.GetBookById200JSONResponse
				decodeJSON(t, w, &resp)
				if resp.Data.Id != tt.book.ID() {
					t.Errorf("data.id = %d, want %d", resp.Data.Id, tt.book.ID())
				}
			}
		})
	}
}

func TestBookHandler_CreateBook(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		svc := &mockBookService{}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books", "{\"title\":}")

		h.CreateBook(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
		var resp generated.ErrorResponse
		decodeJSON(t, w, &resp)
		if resp.Message != "Invalid request body" {
			t.Errorf("message = %q, want %q", resp.Message, "Invalid request body")
		}
		if svc.createCmd != nil {
			t.Error("service called on invalid body, want not called")
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockBookService{
			createFn: func(_ context.Context, _ *command.CreateBookCommand) (*entity.Book, error) {
				return entity.HydrateBook(10, "Go", 2020, 2, nil, nil), nil
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books", `{"title":"Go","year":2020,"author_id":2}`)

		h.CreateBook(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201 (body=%q)", w.Code, w.Body.String())
		}
		var resp generated.CreateBook201JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data.Title != "Go" {
			t.Errorf("data.title = %q, want Go", resp.Data.Title)
		}
		if svc.createCmd == nil || svc.createCmd.Title != "Go" || svc.createCmd.Year != 2020 || svc.createCmd.AuthorID != 2 {
			t.Errorf("service command = %+v, want {Go 2020 2}", svc.createCmd)
		}
	})

	t.Run("author not found", func(t *testing.T) {
		svc := &mockBookService{
			createFn: func(_ context.Context, _ *command.CreateBookCommand) (*entity.Book, error) {
				return nil, entity.ErrAuthorNotFound
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books", `{"title":"Go","year":2020,"author_id":2}`)

		h.CreateBook(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &mockBookService{
			createFn: func(_ context.Context, _ *command.CreateBookCommand) (*entity.Book, error) {
				return nil, entity.ErrValidationError
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books", `{"title":"","year":2020,"author_id":2}`)

		h.CreateBook(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})
}

func TestBookHandler_CreateBookWithAuthor(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		svc := &mockBookService{}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books/with-author", "{\"title\":}")

		h.CreateBookWithAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockBookService{
			createWithAuthorFn: func(_ context.Context, _ *command.CreateBookCommand, _ *command.CreateAuthorCommand) (*entity.Book, error) {
				return entity.HydrateBook(10, "Go", 2020, 5, nil, nil), nil
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books/with-author", `{"title":"Go","year":2020,"author_name":"New"}`)

		h.CreateBookWithAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", w.Code)
		}
		if svc.createWithAuthorCmd == nil || svc.createWithAuthorCmd.Title != "Go" {
			t.Errorf("book command = %+v, want title Go", svc.createWithAuthorCmd)
		}
		if svc.createWithAuthorAut == nil || svc.createWithAuthorAut.Name != "New" {
			t.Errorf("author command = %+v, want name New", svc.createWithAuthorAut)
		}
	})

	t.Run("author already exists", func(t *testing.T) {
		svc := &mockBookService{
			createWithAuthorFn: func(_ context.Context, _ *command.CreateBookCommand, _ *command.CreateAuthorCommand) (*entity.Book, error) {
				return nil, app.ErrAuthorAlreadyExists
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPost, "/books/with-author", `{"title":"Go","year":2020,"author_name":"New"}`)

		h.CreateBookWithAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", w.Code)
		}
	})
}

func TestBookHandler_DeleteBook(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		serviceErr error
		wantStatus int
		wantCall   bool
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantCall: false},
		{name: "not found", id: 5, serviceErr: entity.ErrBookNotFound, wantStatus: http.StatusNotFound, wantCall: true},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantCall: true},
		{name: "success", id: 5, wantStatus: http.StatusNoContent, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockBookService{
				deleteFn: func(_ context.Context, _ int) error {
					return tt.serviceErr
				},
			}
			h := handler.NewBookHandler(svc)
			c, w := newTestContext(http.MethodDelete, "/books/5", "")

			h.DeleteBook(c, tt.id)
			c.Writer.WriteHeaderNow()

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if svc.deleteCalled != tt.wantCall {
				t.Errorf("DeleteBook called = %v, want %v", svc.deleteCalled, tt.wantCall)
			}
		})
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		svc := &mockBookService{}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/0", `{"title":"New"}`)

		h.UpdateBook(c, 0)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
		var resp generated.ErrorResponse
		decodeJSON(t, w, &resp)
		if resp.Message != "invalid book id" {
			t.Errorf("message = %q, want invalid book id", resp.Message)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockBookService{}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/5", "{\"title\":}")

		h.UpdateBook(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockBookService{
			updateFn: func(_ context.Context, _ *command.UpdateBookCommand) (*entity.Book, error) {
				return entity.HydrateBook(5, "New", 2020, 2, nil, nil), nil
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/5", `{"title":"New"}`)

		h.UpdateBook(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		if svc.updateCmd == nil || svc.updateCmd.ID != 5 || svc.updateCmd.Title == nil || *svc.updateCmd.Title != "New" {
			t.Errorf("update command = %+v, want id 5 title New", svc.updateCmd)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockBookService{
			updateFn: func(_ context.Context, _ *command.UpdateBookCommand) (*entity.Book, error) {
				return nil, entity.ErrBookNotFound
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/5", `{"title":"New"}`)

		h.UpdateBook(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("empty patch", func(t *testing.T) {
		svc := &mockBookService{
			updateFn: func(_ context.Context, _ *command.UpdateBookCommand) (*entity.Book, error) {
				return nil, app.ErrEmptyBookPatch
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/5", `{}`)

		h.UpdateBook(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockBookService{
			updateFn: func(_ context.Context, _ *command.UpdateBookCommand) (*entity.Book, error) {
				return nil, errBoom
			},
		}
		h := handler.NewBookHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/books/5", `{"title":"New"}`)

		h.UpdateBook(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}
