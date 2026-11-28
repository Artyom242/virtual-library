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

func TestCategoryHandler_GetCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockCategoryService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Category, int, error) {
				return []*entity.Category{entity.HydrateCategory(1, "Fantasy")}, 1, nil
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodGet, "/categories", "")

		h.GetCategories(c, generated.GetCategoriesParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp generated.GetCategories200JSONResponse
		decodeJSON(t, w, &resp)
		if len(resp.Data) != 1 || resp.Total != 1 {
			t.Errorf("response = %+v, want 1 item total 1", resp)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &mockCategoryService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Category, int, error) {
				return []*entity.Category{}, 0, nil
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodGet, "/categories", "")

		h.GetCategories(c, generated.GetCategoriesParams{})
		c.Writer.WriteHeaderNow()

		var resp generated.GetCategories200JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data == nil || len(resp.Data) != 0 {
			t.Errorf("data = %v, want empty non-nil slice", resp.Data)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockCategoryService{
			getAllFn: func(_ context.Context, _ *pagination.Pagination) ([]*entity.Category, int, error) {
				return nil, 0, errBoom
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodGet, "/categories", "")

		h.GetCategories(c, generated.GetCategoriesParams{})
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}

func TestCategoryHandler_GetCategoryById(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		category   *entity.Category
		serviceErr error
		wantStatus int
		wantMsg    string
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantMsg: "invalid category id"},
		{name: "not found", id: 5, serviceErr: entity.ErrCategoryNotFound, wantStatus: http.StatusNotFound, wantMsg: "category not found"},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantMsg: "Internal server error"},
		{name: "success", id: 5, category: entity.HydrateCategory(5, "Fantasy"), wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoryService{
				getByIDFn: func(_ context.Context, _ int) (*entity.Category, error) {
					return tt.category, tt.serviceErr
				},
			}
			h := handler.NewCategoryHandler(svc)
			c, w := newTestContext(http.MethodGet, "/categories/5", "")

			h.GetCategoryById(c, tt.id)
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

func TestCategoryHandler_CreateCategory(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		svc := &mockCategoryService{}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPost, "/categories", "{\"title\":}")

		h.CreateCategory(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockCategoryService{
			createFn: func(_ context.Context, _ *command.CreateCategoryCommand) (*entity.Category, error) {
				return entity.HydrateCategory(1, "Fantasy"), nil
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPost, "/categories", `{"name":"Fantasy"}`)

		h.CreateCategory(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201 (body=%q)", w.Code, w.Body.String())
		}
		var resp generated.CreateCategory201JSONResponse
		decodeJSON(t, w, &resp)
		if resp.Data.Name != "Fantasy" {
			t.Errorf("data.name = %q, want Fantasy", resp.Data.Name)
		}
		if svc.createCmd == nil || svc.createCmd.Name != "Fantasy" {
			t.Errorf("service command = %+v, want name Fantasy", svc.createCmd)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &mockCategoryService{
			createFn: func(_ context.Context, _ *command.CreateCategoryCommand) (*entity.Category, error) {
				return nil, entity.ErrValidationError
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPost, "/categories", `{"name":""}`)

		h.CreateCategory(c)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		svc := &mockCategoryService{}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/categories/0", `{"name":"Sci-Fi"}`)

		h.UpdateCategory(c, 0)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
		var resp generated.ErrorResponse
		decodeJSON(t, w, &resp)
		if resp.Message != "invalid category id" {
			t.Errorf("message = %q, want invalid category id", resp.Message)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockCategoryService{}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/categories/5", "{\"title\":}")

		h.UpdateCategory(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockCategoryService{
			updateFn: func(_ context.Context, _ *command.UpdateCategoryCommand) (*entity.Category, error) {
				return entity.HydrateCategory(5, "Sci-Fi"), nil
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/categories/5", `{"name":"Sci-Fi"}`)

		h.UpdateCategory(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		if svc.updateCmd == nil || svc.updateCmd.ID != 5 || svc.updateCmd.Name == nil || *svc.updateCmd.Name != "Sci-Fi" {
			t.Errorf("update command = %+v, want id 5 name Sci-Fi", svc.updateCmd)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockCategoryService{
			updateFn: func(_ context.Context, _ *command.UpdateCategoryCommand) (*entity.Category, error) {
				return nil, entity.ErrCategoryNotFound
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/categories/5", `{"name":"Sci-Fi"}`)

		h.UpdateCategory(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", w.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockCategoryService{
			updateFn: func(_ context.Context, _ *command.UpdateCategoryCommand) (*entity.Category, error) {
				return nil, errBoom
			},
		}
		h := handler.NewCategoryHandler(svc)
		c, w := newTestContext(http.MethodPatch, "/categories/5", `{"name":"Sci-Fi"}`)

		h.UpdateCategory(c, 5)
		c.Writer.WriteHeaderNow()

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		serviceErr error
		wantStatus int
		wantCall   bool
	}{
		{name: "invalid id", id: 0, wantStatus: http.StatusBadRequest, wantCall: false},
		{name: "not found", id: 5, serviceErr: entity.ErrCategoryNotFound, wantStatus: http.StatusNotFound, wantCall: true},
		{name: "internal error", id: 5, serviceErr: errBoom, wantStatus: http.StatusInternalServerError, wantCall: true},
		{name: "success", id: 5, wantStatus: http.StatusNoContent, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoryService{
				deleteFn: func(_ context.Context, _ int) error {
					return tt.serviceErr
				},
			}
			h := handler.NewCategoryHandler(svc)
			c, w := newTestContext(http.MethodDelete, "/categories/5", "")

			h.DeleteCategory(c, tt.id)
			c.Writer.WriteHeaderNow()

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if svc.deleteCalled != tt.wantCall {
				t.Errorf("DeleteCategory called = %v, want %v", svc.deleteCalled, tt.wantCall)
			}
		})
	}
}
