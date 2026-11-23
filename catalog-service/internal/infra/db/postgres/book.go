package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/domain/repositories"
	"virtual-library/catalog-service/internal/httpapi/pagination"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

var _ repositories.BookRepository = (*BookRepository)(nil)

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	if pool == nil {
		return &BookRepository{}
	}
	return &BookRepository{pool: pool}
}

type Book struct {
	ID          int
	Title       string
	Year        int
	AuthorID    int
	Description pgtype.Text
}

func (b *BookRepository) Create(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := `
		INSERT INTO books(title, year, author_id, description) 
		VALUES($1, $2, $3, $4) 
		RETURNING id, title, year, author_id, description
	`
	row, err := db.Query(ctx, query, book.Title(), book.Year(), book.AuthorID(), book.Description())
	if err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, mapPostgresError(err)
	}

	newBook := toDomainBook(&bookRow)

	return newBook, nil
}

func (b *BookRepository) GetByID(ctx context.Context, bookID int) (*entities.Book, error) {
	db := txFromContext(ctx, b.pool)

	query := `SELECT id, title, year, author_id, description FROM books WHERE id = $1`
	row, err := db.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book by id: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrBookNotFound
		}
		return nil, fmt.Errorf("scan book: %w", err)
	}

	newBook := toDomainBook(&bookRow)

	return newBook, nil
}

func (b *BookRepository) GetAll(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Book, int, error) {
	db := txFromContext(ctx, b.pool)
	query := `
		SELECT id, title, year, author_id, description
		FROM books 
		ORDER BY id
		LIMIT $1
		OFFSET $2
`

	rows, err := db.Query(ctx, query, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("get all books: %w", err)
	}

	bookRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, 0, fmt.Errorf("scan books: %w", err)
	}

	books := toDomainBooks(bookRows)

	var total int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count authors: %w", err)
	}

	return books, total, nil
}

func (b *BookRepository) Delete(ctx context.Context, bookID int) error {
	db := txFromContext(ctx, b.pool)
	query := "DELETE FROM books WHERE id = $1"
	tag, err := db.Exec(ctx, query, bookID)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.ErrBookNotFound
	}

	return nil
}

func (b *BookRepository) Update(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := `
		UPDATE books
		SET title     = COALESCE($1, title),
		    year      = COALESCE($2, year),
		    author_id = COALESCE($3, author_id),
		    description = COALESCE($4, description)
		WHERE id = $5
		RETURNING id, title, author_id, year, description
	`
	row, err := db.Query(ctx, query, book.Title(), book.Year(), book.AuthorID(), book.Description(), book.ID())
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrBookNotFound
		}

		return nil, fmt.Errorf("update book: %w", mapPostgresError(err))
	}

	newBook := toDomainBook(&bookRow)

	return newBook, nil
}

func mapPostgresError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code != pgerrcode.ForeignKeyViolation {
		return err
	}

	switch pgErr.ConstraintName {
	case "fk_books_authors":
		return entities.ErrAuthorNotFound
	}
	// add fk_books_categories case

	return err
}

func toDomainBooks(bookRows []Book) []*entities.Book {
	books := make([]*entities.Book, 0, len(bookRows))
	for _, row := range bookRows {
		book := toDomainBook(&row)
		books = append(books, book)
	}
	return books
}

func toDomainBook(b *Book) *entities.Book {
	var description *string

	if b.Description.Valid {
		description = &b.Description.String
	}

	return entities.HydrateBook(
		b.ID,
		b.Title,
		b.Year,
		b.AuthorID,
		description,
	)
}
