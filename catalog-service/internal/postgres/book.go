package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

var _ domain.BookRepository = (*BookRepository)(nil)

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		pool: pool,
	}
}

func (b *BookRepository) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := "INSERT INTO books(title, author_id, year) VALUES($1, $2, $3) returning id, title, author_id, year;"
	err := db.QueryRow(ctx, query, book.ID, book.Title, book.AuthorID, book.Year).Scan(&book.ID, &book.Title, &book.AuthorID, &book.Year)

	if err != nil {
		return book, fmt.Errorf("create book: %w", err)
	}

	return book, nil
}

func (b *BookRepository) GetByID(ctx context.Context, bookID int) (*domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := "SELECT * FROM books WHERE id = $1"
	var book domain.Book

	err := db.QueryRow(ctx, query, bookID).Scan()
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBookNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get book by id: %w", err)
	}

	return &book, nil
}

func (b *BookRepository) GetAll(ctx context.Context) ([]domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := "SELECT books.id, books.title, books.author_id, books.year FROM books"
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all books: %w", err)
	}

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Book])
	if err != nil {
		return nil, fmt.Errorf("collect books: %w", err)
	}
	return books, nil
}

func (b *BookRepository) Delete(ctx context.Context, bookID int) error {
	db := txFromContext(ctx, b.pool)
	query := "DELETE FROM books WHERE id = $1"
	tag, err := db.Exec(ctx, query, bookID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrBookNotFound
	}
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}

	return nil
}

func (b *BookRepository) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := `
			UPDATE books
			SET title = $1, author_id = $2, year = $3
			WHERE id = $4
			RETURNING id, title, author_id, year
		`
	err := db.QueryRow(ctx, query, book.Title, book.AuthorID, book.Year, book.ID).
		Scan(&book.ID, &book.Title, &book.AuthorID, &book.Year)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return book, fmt.Errorf("update book: %w", domain.ErrBookNotFound)
		}
		return book, fmt.Errorf("update book: %w", err)
	}

	return book, nil
}
