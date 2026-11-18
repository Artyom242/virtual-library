package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepository struct {
	pool *pgxpool.Pool
}

func NewAuthorRepository(pool *pgxpool.Pool) *AuthorRepository {
	return &AuthorRepository{pool: pool}
}

var _ domain.AuthorRepository = (*AuthorRepository)(nil)

func (a *AuthorRepository) Create(c context.Context, author domain.Author) (domain.Author, error) {
	db := txFromContext(c, a.pool)
	query := "INSERT INTO authors(name) VALUES($1) returning id, name"
	err := db.QueryRow(c, query, author.Name).Scan(&author.ID, &author.Name)
	if err != nil {
		return domain.Author{}, fmt.Errorf("create author: %w", err)
	}

	return author, nil
}

func (a *AuthorRepository) GetByID(c context.Context, authorID int) (*domain.Author, error) {
	var author domain.Author

	db := txFromContext(c, a.pool)
	query := "SELECT authors.id, authors.name FROM authors WHERE id = $1"
	err := db.QueryRow(c, query, authorID).Scan(&author.ID, &author.Name)
	if err != nil {
		return &author, fmt.Errorf("create author: %w", err)
	}

	return &author, nil
}

func (a *AuthorRepository) GetAll(c context.Context) ([]domain.Author, error) {
	db := txFromContext(c, a.pool)
	query := "SELECT authors.id, authors.name FROM authors"
	rows, err := db.Query(c, query)
	if err != nil {
		return []domain.Author{}, fmt.Errorf("get all authors: %w", err)
	}

	authors, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Author])
	if err != nil {
		return []domain.Author{}, fmt.Errorf("collect authors: %w", err)
	}

	return authors, nil
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
	var author domain.Author

	query := `UPDATE authors 
		SET name = COALESCE($1, name) 
		WHERE id = $2 
		RETURNING id, name`
	err := db.QueryRow(c, query, patch.Name, id).Scan(&author.ID, &author.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Author{}, domain.ErrAuthorNotFound
		}
		return domain.Author{}, fmt.Errorf("update author: %w", err)
	}
	return author, nil
}
