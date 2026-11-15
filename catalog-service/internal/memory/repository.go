package memory

import (
	"errors"
	"virtual-library/catalog-service/internal/catalog"
)

type MemoryBookRepository struct {
	books  map[int]catalog.Book
	nextID int
}

var _ catalog.BookRepository = (*MemoryBookRepository)(nil)

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		books:  make(map[int]catalog.Book),
		nextID: 1,
	}
}

func (c *MemoryBookRepository) Create(book catalog.Book) (catalog.Book, error) {
	book.ID = c.nextID

	c.books[book.ID] = book
	c.nextID++

	return book, nil
}

func (c *MemoryBookRepository) GetByID(bookID int) (*catalog.Book, error) {
	book, ok := c.books[bookID]

	if !ok {
		return nil, errors.New("book not found")
	}

	return &book, nil
}

func (c *MemoryBookRepository) GetAll() ([]catalog.Book, error) {
	books := make([]catalog.Book, 0, len(c.books))

	for _, book := range c.books {
		books = append(books, book)
	}

	return books, nil
}

func (c *MemoryBookRepository) Delete(bookID int) error {
	_, ok := c.books[bookID]

	if !ok {
		return errors.New("book not found")
	}

	delete(c.books, bookID)
	return nil
}
