package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgUniqueViolation = "23505"

type AuthorRepository struct {
	pool *pgxpool.Pool
}

type Author struct {
	ID   int
	Name string
}

func NewAuthorRepository(pool *pgxpool.Pool) *AuthorRepository {
	if pool == nil {
		return &AuthorRepository{}
	}
	return &AuthorRepository{pool: pool}
}

var _ domain.AuthorRepository = (*AuthorRepository)(nil)

func (a *AuthorRepository) Create(c context.Context, author domain.Author) (domain.Author, error) {
	db := txFromContext(c, a.pool)

	query := "INSERT INTO authors(name) VALUES($1) returning id, name"
	rows, err := db.Query(c, query, author.Name())
	if err != nil {
		return domain.Author{}, err
	}

	authorRow, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[Author],
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.Author{}, domain.ErrAuthorAlreadyExists
		}

		return domain.Author{}, fmt.Errorf("create author: %w", err)
	}

	newAuthor, err := toDomainAuthor(authorRow)
	if err != nil {
		return domain.Author{}, domain.ErrMapAuthorToDomain
	}
	return newAuthor, nil
}

func (a *AuthorRepository) GetByID(c context.Context, authorID int) (*domain.Author, error) {
	db := txFromContext(c, a.pool)
	query := "SELECT authors.id, authors.name FROM authors WHERE id = $1"
	row, err := db.Query(c, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("get author by id: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(
		row,
		pgx.RowToStructByName[Author],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAuthorNotFound
		}
		return nil, fmt.Errorf("scan author: %w", err)
	}

	author, err := toDomainAuthor(authorRow)
	if err != nil {
		return nil, fmt.Errorf("map author to domain: %w", err)
	}

	return &author, nil
}

func (a *AuthorRepository) GetAll(c context.Context, page int, limit int) ([]domain.Author, int, error) {
	db := txFromContext(c, a.pool)
	offset := (page - 1) * limit

	query := `
		SELECT authors.id, authors.name 
		FROM authors
		ORDER BY authors.id
		LIMIT $1
		OFFSET $2
		`
	rows, err := db.Query(c, query, limit, offset)
	if err != nil {
		return []domain.Author{}, 0, fmt.Errorf("get all authors: %w", err)
	}

	authorRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Author])
	if err != nil {
		return []domain.Author{}, 0, fmt.Errorf("scan author: %w", err)
	}

	authors, err := toDomainAuthors(authorRows)
	if err != nil {
		return []domain.Author{}, 0, fmt.Errorf("map author to domain: %w", err)
	}

	var total int
	err = db.QueryRow(c, "SELECT COUNT(*) FROM authors").Scan(&total)
	if err != nil {
		return []domain.Author{}, 0, fmt.Errorf("count authors: %w", err)
	}

	return authors, total, nil
}

func (a *AuthorRepository) GetOrCreate(ctx context.Context, author domain.Author) (domain.Author, error) {
	db := txFromContext(ctx, a.pool)

	query := `
		INSERT INTO authors(name) 
		VALUES($1)  
		ON CONFLICT (name) 
		    DO NOTHING 
			RETURNING id, name
	`

	row, err := db.Query(ctx, query, author.Name())
	if err != nil {
		return domain.Author{}, fmt.Errorf("create author: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			query = `SELECT id, name FROM authors WHERE name = $1`

			row, err = db.Query(ctx, query, author.Name())
			if err != nil {
				return domain.Author{}, fmt.Errorf("get author: %w", err)
			}

			authorRow, err = pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
			if err != nil {
				return domain.Author{}, fmt.Errorf("scan author: %w", err)
			}
		}

		return domain.Author{}, fmt.Errorf("get or create author: %w", err)
	}

	newAuthor, err := toDomainAuthor(authorRow)
	if err != nil {
		return domain.Author{}, fmt.Errorf("map author to domain: %w", err)
	}
	return newAuthor, nil
}

func (a *AuthorRepository) HasBooks(c context.Context, authorID int) (bool, error) {
	db := txFromContext(c, a.pool)
	query := "SELECT EXISTS(SELECT 1 FROM books WHERE author_id = $1)"
	var hasBooks bool
	err := db.QueryRow(c, query, authorID).Scan(&hasBooks)
	if err != nil {
		return false, fmt.Errorf("check author books: %w", err)
	}
	return hasBooks, nil
}

func (a *AuthorRepository) Delete(c context.Context, authorID int) error {
	db := txFromContext(c, a.pool)
	query := "DELETE FROM authors WHERE id = $1"
	tag, err := db.Exec(c, query, authorID)
	if err != nil {
		return fmt.Errorf("delete author: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrAuthorNotFound
	}

	return nil
}

func (a *AuthorRepository) Update(c context.Context, id int, patch domain.AuthorPatch) (domain.Author, error) {
	db := txFromContext(c, a.pool)

	query := `UPDATE authors 
		SET name = COALESCE($1, name) 
		WHERE id = $2 
		RETURNING id, name`
	row, err := db.Query(c, query, patch.Name, id)
	if err != nil {
		return domain.Author{}, fmt.Errorf("update author: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.Author{}, domain.ErrAuthorAlreadyExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Author{}, domain.ErrAuthorNotFound
		}
		return domain.Author{}, fmt.Errorf("update author: %w", err)
	}

	author, err := toDomainAuthor(authorRow)
	if err != nil {
		return domain.Author{}, fmt.Errorf("map author to domain: %w", err)
	}

	return author, nil
}

func toDomainAuthors(authorRows []Author) ([]domain.Author, error) {
	authors := make([]domain.Author, 0, len(authorRows))

	for _, row := range authorRows {
		author, err := toDomainAuthor(row)
		if err != nil {
			return nil, fmt.Errorf("map author: %w", err)
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func toDomainAuthor(a Author) (domain.Author, error) {
	return domain.HydrateAuthor(a.ID, a.Name), nil
}
