package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

var _ domain.BookRepository = (*BookRepository)(nil)

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

func (b *BookRepository) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := `
		INSERT INTO books(title, year, author_id, description) 
		VALUES($1, $2, $3, $4) 
		RETURNING id, title, year, author_id, description
	`
	row, err := db.Query(ctx, query, book.Title(), book.Year(), book.AuthorID(), book.Description())
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		return domain.Book{}, mapPostgresError(err)
	}

	newBook, err := toDomainBook(bookRow)
	if err != nil {
		return domain.Book{}, domain.ErrMapBookToDomain
	}
	return newBook, nil
}

func (b *BookRepository) GetByID(ctx context.Context, bookID int) (*domain.Book, error) {
	db := txFromContext(ctx, b.pool)
	query := `SELECT id, title, year, author_id, description FROM books WHERE id = $1`
	row, err := db.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book by id: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		return nil, fmt.Errorf("scan book: %w", err)
	}

	newBook, err := toDomainBook(bookRow)
	if err != nil {
		return nil, domain.ErrMapBookToDomain
	}

	return &newBook, nil
}

func (b *BookRepository) GetAll(ctx context.Context, page int, limit int) ([]domain.Book, int, error) {
	db := txFromContext(ctx, b.pool)
	offset := (page - 1) * limit
	query := `
		SELECT id, title, year, author_id, description
		FROM books 
		ORDER BY id
		LIMIT $1
		OFFSET $2
`

	rows, err := db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get all books: %w", err)
	}

	bookRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, 0, fmt.Errorf("scan books: %w", err)
	}
	books, err := toDomainBooks(bookRows)
	if err != nil {
		return nil, 0, domain.ErrMapBookToDomain
	}

	var total int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return []domain.Book{}, 0, fmt.Errorf("count authors: %w", err)
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
		return domain.ErrBookNotFound
	}

	return nil
}

func (b *BookRepository) Update(ctx context.Context, id int, patch domain.BookPatch) (domain.Book, error) {
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
	row, err := db.Query(ctx, query, patch.Title, patch.Year, patch.AuthorID, patch.Description, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("update book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, domain.ErrBookNotFound
		}

		return domain.Book{}, fmt.Errorf(
			"update book: %w",
			mapPostgresError(err),
		)
	}

	newBook, err := toDomainBook(bookRow)
	if err != nil {
		return domain.Book{}, domain.ErrMapBookToDomain
	}

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
		return domain.ErrAuthorNotFound
	}
	// add fk_books_categories case

	return err
}

func toDomainBooks(bookRows []Book) ([]domain.Book, error) {
	books := make([]domain.Book, 0, len(bookRows))
	for _, row := range bookRows {
		book, err := toDomainBook(row)
		if err != nil {
			return nil, domain.ErrMapBookToDomain
		}
		books = append(books, book)
	}
	return books, nil
}

func toDomainBook(b Book) (domain.Book, error) {
	var description *string

	if b.Description.Valid {
		description = &b.Description.String
	}

	return domain.HydrateBook(
		b.ID,
		b.Title,
		b.Year,
		b.AuthorID,
		description,
	), nil
}
