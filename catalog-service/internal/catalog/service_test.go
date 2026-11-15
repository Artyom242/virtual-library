package catalog

import (
	"errors"
	"reflect"
	"testing"
)

type fakeBookRepository struct {
	book        Book
	books       []Book
	err         error
	gotID       int
	lastGetByID int
}

var _ BookRepository = (*fakeBookRepository)(nil)

func (f *fakeBookRepository) Create(book Book) (Book, error) {
	f.book = book
	return f.book, f.err
}
func (f *fakeBookRepository) GetByID(bookID int) (*Book, error) {
	f.lastGetByID = bookID
	if f.err != nil {
		return nil, f.err
	}
	return &f.book, f.err
}
func (f *fakeBookRepository) GetAll() ([]Book, error) {
	return f.books, f.err
}
func (f *fakeBookRepository) Delete(bookID int) error {
	f.gotID = bookID
	return f.err
}

func TestCatalogService_Create(t *testing.T) {
	expected := Book{
		ID:     1,
		Title:  "test book",
		Author: "test author",
		Year:   2021,
	}

	fakeRepo := &fakeBookRepository{}
	s := NewService(fakeRepo)

	got, err := s.CreateBook(expected)

	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected: %v, got: %v", expected, got)
	}
}

func TestCatalogService_GetAll(t *testing.T) {
	expected := []Book{
		{
			ID:     1,
			Title:  "test book",
			Author: "test author",
			Year:   2021,
		},
		{
			ID:     2,
			Title:  "second book",
			Author: "second author",
			Year:   2022,
		},
	}

	fakeRepo := &fakeBookRepository{
		books: expected,
	}
	s := NewService(fakeRepo)
	books, err := s.GetAll()

	if err != nil {
		t.Fatalf("GetAll() returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(books, expected) {
		t.Fatalf("Expected: %v, got: %v", expected, books)
	}
}
func TestCatalogService_GetAll_Error(t *testing.T) {
	fakeRepo := &fakeBookRepository{
		err: errors.New("some error"),
	}
	s := NewService(fakeRepo)
	books, err := s.GetAll()
	if err == nil {
		t.Fatalf("GetAll() expected an error, got nil")
	}
	if books != nil {
		t.Fatalf("GetAll() returned books %v, want nil", books)
	}
}

func TestCatalogService_GetByID(t *testing.T) {
	expected := Book{
		ID:     1,
		Title:  "test book",
		Author: "test author",
		Year:   2021,
	}
	fakeRepo := &fakeBookRepository{
		book: expected,
	}
	s := NewService(fakeRepo)
	got, err := s.GetByID(expected.ID)

	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}
	if got == nil {
		t.Fatal("GetByID() returned nil, want book")
	}
	if fakeRepo.lastGetByID != expected.ID {
		t.Fatalf("GetByID() returned book id = %d, want %d", got.ID, expected.ID)
	}
	if !reflect.DeepEqual(*got, expected) {
		t.Fatalf("Expected: %v, got: %v", expected, got)
	}
}
func TestCatalogService_GetByID_Error(t *testing.T) {
	fakeRepo := &fakeBookRepository{
		err: errors.New("some error"),
	}
	s := NewService(fakeRepo)
	got, err := s.GetByID(1)
	if err == nil {
		t.Fatalf("GetByID() expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("GetByID() returned books %v, want nil", got)
	}
}

func TestCatalogService_Delete(t *testing.T) {
	expectedID := 1
	fakeRepo := &fakeBookRepository{
		books: []Book{
			{
				ID:     1,
				Title:  "test book",
				Author: "test author",
				Year:   2021,
			},
		},
	}
	s := NewService(fakeRepo)
	err := s.Delete(expectedID)
	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
	if fakeRepo.gotID != expectedID {
		t.Fatalf("Delete() called with unexpected ID %v, want %v", fakeRepo.gotID, expectedID)
	}
}
func TestCatalogService_Delete_Error(t *testing.T) {
	expectedID := 1
	fakeRepo := &fakeBookRepository{
		err: errors.New("some error"),
	}
	s := NewService(fakeRepo)
	err := s.Delete(expectedID)

	if err == nil {
		t.Fatalf("Delete() expected an error, got nil")
	}
	if fakeRepo.gotID != expectedID {
		t.Fatalf("Delete() called with unexpected ID %v, want %v", fakeRepo.gotID, expectedID)
	}
}
