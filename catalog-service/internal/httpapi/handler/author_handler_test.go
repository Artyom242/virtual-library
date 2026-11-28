package handler_test

import (
	"context"
	"net/http"
	"testing"

	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/httpapi/handler"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

func TestAuthorHandler_GetAuthors(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthorService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Author, int, error) {
				return []*entity.Author{entity.HydrateAuthor(1, "John")}, 1, nil
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodGet, "/authors", "")

		h.GetAuthors(c, generated.GetAuthorsParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp generated.GetAuthors200JSONResponse
		decodeJSON(t, w, &resp)
		if len(resp.Data) != 1 || resp.Total != 1 {
			t.Errorf("response = %+v, want 1 item total 1", resp)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &mockAuthorService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Author, int, error) {
				return []*entity.Author{}, 0, nil
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodGet, "/authors", "")

		h.GetAuthors(c, generated.GetAuthorsParams{})
		c.Writer.WriteHeaderNow()

		var resp generated.GetAuthors200JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data == nil || len(resp.Data) != 0 {
			t.Errorf("data = %v, want empty non-nil slice", resp.Data)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockAuthorService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Author, int, error) {
				return nil, 0, errBoom
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodGet, "/authors", "")

		h.GetAuthors(c, generated.GetAuthorsParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}

func TestAuthorHandler_GetAuthorById(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		author     *entity.Author
		serviceErr error
		wantStatus int
		wantMsg    string
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantMsg: "invalid author id"},
		{name: "not found", id: 5, serviceErr: entity.ErrAuthorNotFound, wantStatus: http.StatusNotFound, wantMsg: "author not found"},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantMsg: "Internal server error"},
		{name: "success", id: 5, author: entity.HydrateAuthor(5, "John"), wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockAuthorService{
				getByIDFn: func(_ context.Context, _ int) (*entity.Author, error) {
					return tt.author, tt.serviceErr
				},
			}
			h := handler.NewAuthorHandler(svc)
			c, w := newTestContext(http.MethodGet, "/authors/5", "")

			h.GetAuthorById(c, tt.id)
			c.Writer.WriteHeaderNow()

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantMsg != "" {
				var resp generated.ErrorResponse
				decodeJSON(t, w, &resp)
				if resp.Message != tt.wantMsg {
					t.Errorf("message = %q, want %q", resp.Message, tt.wantMsg)
				}
			}
		})
	}
}

func TestAuthorHandler_CreateAuthor(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAuthorService{}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPost, "/authors", "{\"title\":}")

		h.CreateAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockAuthorService{
			createFn: func(_ context.Context, _ *command.CreateAuthorCommand) (*entity.Author, error) {
				return entity.HydrateAuthor(1, "John"), nil
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPost, "/authors", `{"name":"John"}`)

		h.CreateAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201 (body=%q)", w.Code, w.Body.String())
		}
		var resp generated.CreateAuthor201JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data.Name != "John" {
			t.Errorf("data.name = %q, want John", resp.Data.Name)
		}
		if svc.createCmd == nil || svc.createCmd.Name != "John" {
			t.Errorf("service command = %+v, want name John", svc.createCmd)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &mockAuthorService{
			createFn: func(_ context.Context, _ *command.CreateAuthorCommand) (*entity.Author, error) {
				return nil, entity.ErrValidationError
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPost, "/authors", `{"name":""}`)

		h.CreateAuthor(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})
}

func TestAuthorHandler_UpdateAuthor(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAuthorService{}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/authors/0", `{"name":"Jane"}`)

		h.UpdateAuthor(c, 0)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
		var resp generated.ErrorResponse
		decodeJSON(t, w, &resp)
		if resp.Message != "invalid author id" {
			t.Errorf("message = %q, want invalid author id", resp.Message)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAuthorService{}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/authors/5", "{\"title\":}")

		h.UpdateAuthor(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockAuthorService{
			updateFn: func(_ context.Context, _ *command.UpdateAuthorCommand) (*entity.Author, error) {
				return entity.HydrateAuthor(5, "Jane"), nil
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/authors/5", `{"name":"Jane"}`)

		h.UpdateAuthor(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		if svc.updateCmd == nil || svc.updateCmd.ID != 5 || svc.updateCmd.Name == nil || *svc.updateCmd.Name != "Jane" {
			t.Errorf("update command = %+v, want id 5 name Jane", svc.updateCmd)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockAuthorService{
			updateFn: func(_ context.Context, _ *command.UpdateAuthorCommand) (*entity.Author, error) {
				return nil, entity.ErrAuthorNotFound
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/authors/5", `{"name":"Jane"}`)

		h.UpdateAuthor(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockAuthorService{
			updateFn: func(_ context.Context, _ *command.UpdateAuthorCommand) (*entity.Author, error) {
				return nil, errBoom
			},
		}
		h := handler.NewAuthorHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/authors/5", `{"name":"Jane"}`)

		h.UpdateAuthor(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}

func TestAuthorHandler_DeleteAuthor(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		serviceErr error
		wantStatus int
		wantCall   bool
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantCall: false},
		{name: "not found", id: 5, serviceErr: entity.ErrAuthorNotFound, wantStatus: http.StatusNotFound, wantCall: true},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantCall: true},
		{name: "success", id: 5, wantStatus: http.StatusNoContent, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockAuthorService{
				deleteFn: func(_ context.Context, _ int) error {
					return tt.serviceErr
				},
			}
			h := handler.NewAuthorHandler(svc)
			c, w := newTestContext(http.MethodDelete, "/authors/5", "")

			h.DeleteAuthor(c, tt.id)
			c.Writer.WriteHeaderNow()

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if svc.deleteCalled != tt.wantCall {
				t.Errorf("DeleteAuthor called = %v, want %v", svc.deleteCalled, tt.wantCall)
			}
		})
	}
}
