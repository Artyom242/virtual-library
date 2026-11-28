package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"

	"github.com/gin-gonic/gin"
)

func newTestContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req

	return c, w
}

func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()

	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		t.Fatalf("decode response: %v (body=%q)", err, w.Body.String())
	}
}

type mockBookService struct {
	createFn           func(ctx context.Context, cmd *command.CreateBookCommand) (*entity.Book, error)
	updateFn           func(ctx context.Context, cmd *command.UpdateBookCommand) (*entity.Book, error)
	deleteFn           func(ctx context.Context, id int) error
	getByIDFn          func(ctx context.Context, id int) (*entity.Book, error)
	getAllFn           func(ctx context.Context, p *pagination.Pagination) ([]*entity.Book, int, error)
	createWithAuthorFn func(ctx context.Context, book *command.CreateBookCommand, author *command.CreateAuthorCommand) (*entity.Book, error)

	createCmd           *command.CreateBookCommand
	updateCmd           *command.UpdateBookCommand
	createWithAuthorCmd *command.CreateBookCommand
	createWithAuthorAut *command.CreateAuthorCommand
	deleteCalled        bool
	getByIDCalled       bool
}

func (m *mockBookService) CreateBook(ctx context.Context, cmd *command.CreateBookCommand) (*entity.Book, error) {
	m.createCmd = cmd
	if m.createFn != nil {
		return m.createFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockBookService) UpdateBook(ctx context.Context, cmd *command.UpdateBookCommand) (*entity.Book, error) {
	m.updateCmd = cmd
	if m.updateFn != nil {
		return m.updateFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockBookService) DeleteBook(ctx context.Context, id int) error {
	m.deleteCalled = true
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockBookService) GetBookByID(ctx context.Context, id int) (*entity.Book, error) {
	m.getByIDCalled = true
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockBookService) GetAllBooks(ctx context.Context, p *pagination.Pagination) ([]*entity.Book, int, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}

func (m *mockBookService) CreateBookWithNewAuthor(ctx context.Context, book *command.CreateBookCommand, author *command.CreateAuthorCommand) (*entity.Book, error) {
	m.createWithAuthorCmd = book
	m.createWithAuthorAut = author
	if m.createWithAuthorFn != nil {
		return m.createWithAuthorFn(ctx, book, author)
	}

	return nil, nil
}

type mockAuthorService struct {
	createFn     func(ctx context.Context, cmd *command.CreateAuthorCommand) (*entity.Author, error)
	updateFn     func(ctx context.Context, cmd *command.UpdateAuthorCommand) (*entity.Author, error)
	deleteFn     func(ctx context.Context, id int) error
	getByIDFn    func(ctx context.Context, id int) (*entity.Author, error)
	getAllFn     func(ctx context.Context, p *pagination.Pagination) ([]*entity.Author, int, error)
	createCmd    *command.CreateAuthorCommand
	updateCmd    *command.UpdateAuthorCommand
	deleteCalled bool
}

func (m *mockAuthorService) CreateAuthor(ctx context.Context, cmd *command.CreateAuthorCommand) (*entity.Author, error) {
	m.createCmd = cmd
	if m.createFn != nil {
		return m.createFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockAuthorService) UpdateAuthor(ctx context.Context, cmd *command.UpdateAuthorCommand) (*entity.Author, error) {
	m.updateCmd = cmd
	if m.updateFn != nil {
		return m.updateFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockAuthorService) DeleteAuthor(ctx context.Context, id int) error {
	m.deleteCalled = true
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockAuthorService) GetAuthorByID(ctx context.Context, id int) (*entity.Author, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockAuthorService) GetAllAuthors(ctx context.Context, p *pagination.Pagination) ([]*entity.Author, int, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}

type mockCategoryService struct {
	createFn     func(ctx context.Context, cmd *command.CreateCategoryCommand) (*entity.Category, error)
	updateFn     func(ctx context.Context, cmd *command.UpdateCategoryCommand) (*entity.Category, error)
	deleteFn     func(ctx context.Context, id int) error
	getByIDFn    func(ctx context.Context, id int) (*entity.Category, error)
	getAllFn     func(ctx context.Context, p *pagination.Pagination) ([]*entity.Category, int, error)
	createCmd    *command.CreateCategoryCommand
	updateCmd    *command.UpdateCategoryCommand
	deleteCalled bool
}

func (m *mockCategoryService) CreateCategory(ctx context.Context, cmd *command.CreateCategoryCommand) (*entity.Category, error) {
	m.createCmd = cmd
	if m.createFn != nil {
		return m.createFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockCategoryService) UpdateCategory(ctx context.Context, cmd *command.UpdateCategoryCommand) (*entity.Category, error) {
	m.updateCmd = cmd
	if m.updateFn != nil {
		return m.updateFn(ctx, cmd)
	}

	return nil, nil
}

func (m *mockCategoryService) DeleteCategory(ctx context.Context, id int) error {
	m.deleteCalled = true
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockCategoryService) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockCategoryService) GetAllCategories(ctx context.Context, p *pagination.Pagination) ([]*entity.Category, int, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}
