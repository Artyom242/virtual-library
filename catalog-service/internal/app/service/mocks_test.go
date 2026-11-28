package service_test

import (
	"context"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	infraRepo "virtual-library/catalog-service/internal/infra/db/repository"
)

func strPtr(s string) *string { return &s }

func intPtr(v int) *int { return &v }

type mockTransactor struct {
	do func(ctx context.Context, fn func(context.Context) error) error
}

func (m *mockTransactor) Do(ctx context.Context, fn func(context.Context) error) error {
	if m.do != nil {
		return m.do(ctx, fn)
	}

	return fn(ctx)
}

type mockAuthorRepo struct {
	createFn      func(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error)
	getByIDFn     func(ctx context.Context, id int) (*entity.Author, error)
	getAllFn      func(ctx context.Context, p *pagination.Pagination) ([]*entity.Author, int, error)
	deleteFn      func(ctx context.Context, id int) error
	updateFn      func(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error)
	getOrCreateFn func(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error)
	hasBooksFn    func(ctx context.Context, id int) (bool, error)

	createCalls      []infraRepo.IAuthor
	getByIDCalls     []int
	getAllCalls      int
	deleteCalls      []int
	updateCalls      []infraRepo.IAuthor
	getOrCreateCalls []infraRepo.IAuthor
	hasBooksCalls    []int
}

func (m *mockAuthorRepo) Create(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error) {
	m.createCalls = append(m.createCalls, author)
	if m.createFn != nil {
		return m.createFn(ctx, author)
	}

	return nil, nil
}

func (m *mockAuthorRepo) GetByID(ctx context.Context, id int) (*entity.Author, error) {
	m.getByIDCalls = append(m.getByIDCalls, id)
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockAuthorRepo) GetAll(ctx context.Context, p *pagination.Pagination) ([]*entity.Author, int, error) {
	m.getAllCalls++
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}

func (m *mockAuthorRepo) Delete(ctx context.Context, id int) error {
	m.deleteCalls = append(m.deleteCalls, id)
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockAuthorRepo) Update(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error) {
	m.updateCalls = append(m.updateCalls, author)
	if m.updateFn != nil {
		return m.updateFn(ctx, author)
	}

	return nil, nil
}

func (m *mockAuthorRepo) GetOrCreate(ctx context.Context, author infraRepo.IAuthor) (*entity.Author, error) {
	m.getOrCreateCalls = append(m.getOrCreateCalls, author)
	if m.getOrCreateFn != nil {
		return m.getOrCreateFn(ctx, author)
	}

	return nil, nil
}

func (m *mockAuthorRepo) HasBooks(ctx context.Context, id int) (bool, error) {
	m.hasBooksCalls = append(m.hasBooksCalls, id)
	if m.hasBooksFn != nil {
		return m.hasBooksFn(ctx, id)
	}

	return false, nil
}

type mockBookRepo struct {
	createFn     func(ctx context.Context, book infraRepo.IBook) (*entity.Book, error)
	getByIDFn    func(ctx context.Context, id int) (*entity.Book, error)
	getAllFn     func(ctx context.Context, p *pagination.Pagination) ([]*entity.Book, int, error)
	deleteFn     func(ctx context.Context, id int) error
	updateFn     func(ctx context.Context, book infraRepo.IBook) (*entity.Book, error)
	createCalls  []infraRepo.IBook
	getByIDCalls []int
	getAllCalls  int
	deleteCalls  []int
	updateCalls  []infraRepo.IBook
}

func (m *mockBookRepo) Create(ctx context.Context, book infraRepo.IBook) (*entity.Book, error) {
	m.createCalls = append(m.createCalls, book)
	if m.createFn != nil {
		return m.createFn(ctx, book)
	}

	return nil, nil
}

func (m *mockBookRepo) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	m.getByIDCalls = append(m.getByIDCalls, id)
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockBookRepo) GetAll(ctx context.Context, p *pagination.Pagination) ([]*entity.Book, int, error) {
	m.getAllCalls++
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}

func (m *mockBookRepo) Delete(ctx context.Context, id int) error {
	m.deleteCalls = append(m.deleteCalls, id)
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockBookRepo) Update(ctx context.Context, book infraRepo.IBook) (*entity.Book, error) {
	m.updateCalls = append(m.updateCalls, book)
	if m.updateFn != nil {
		return m.updateFn(ctx, book)
	}

	return nil, nil
}

type mockCategoryRepo struct {
	createFn     func(ctx context.Context, category infraRepo.ICategory) (*entity.Category, error)
	getByIDFn    func(ctx context.Context, id int) (*entity.Category, error)
	getAllFn     func(ctx context.Context, p *pagination.Pagination) ([]*entity.Category, int, error)
	deleteFn     func(ctx context.Context, id int) error
	updateFn     func(ctx context.Context, category infraRepo.ICategory) (*entity.Category, error)
	createCalls  []infraRepo.ICategory
	getByIDCalls []int
	getAllCalls  int
	deleteCalls  []int
	updateCalls  []infraRepo.ICategory
}

func (m *mockCategoryRepo) Create(ctx context.Context, category infraRepo.ICategory) (*entity.Category, error) {
	m.createCalls = append(m.createCalls, category)
	if m.createFn != nil {
		return m.createFn(ctx, category)
	}

	return nil, nil
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id int) (*entity.Category, error) {
	m.getByIDCalls = append(m.getByIDCalls, id)
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return nil, nil
}

func (m *mockCategoryRepo) GetAll(ctx context.Context, p *pagination.Pagination) ([]*entity.Category, int, error) {
	m.getAllCalls++
	if m.getAllFn != nil {
		return m.getAllFn(ctx, p)
	}

	return nil, 0, nil
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id int) error {
	m.deleteCalls = append(m.deleteCalls, id)
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}

func (m *mockCategoryRepo) Update(ctx context.Context, category infraRepo.ICategory) (*entity.Category, error) {
	m.updateCalls = append(m.updateCalls, category)
	if m.updateFn != nil {
		return m.updateFn(ctx, category)
	}

	return nil, nil
}
