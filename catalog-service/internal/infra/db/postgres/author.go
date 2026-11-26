package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/domain/repository"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	"virtual-library/catalog-service/internal/infra/db/postgres/query"
	IRepository "virtual-library/catalog-service/internal/infra/db/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

var _ repository.AuthorRepository = (*AuthorRepository)(nil)

func (a *AuthorRepository) Create(c context.Context, author IRepository.IAuthor) (*entity.Author, error) {
	db := txFromContext(c, a.pool)

	rows, err := db.Query(c, query.CreateAuthor, author.Name())
	if err != nil {
		return nil, err
	}

	authorRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Author])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, app.ErrAuthorAlreadyExists
		}

		return nil, fmt.Errorf("create author: %w", err)
	}

	newAuthor := toDomainAuthor(&authorRow)

	return newAuthor, nil
}

func (a *AuthorRepository) GetByID(c context.Context, authorID int) (*entity.Author, error) {
	db := txFromContext(c, a.pool)

	row, err := db.Query(c, query.GetAuthorByID, authorID)
	if err != nil {
		return nil, fmt.Errorf("get author by id: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrAuthorNotFound
		}
		return nil, fmt.Errorf("scan author: %w", err)
	}

	author := toDomainAuthor(&authorRow)

	return author, nil
}

func (a *AuthorRepository) GetAll(c context.Context, pagination *pagination.Pagination) ([]*entity.Author, int, error) {
	db := txFromContext(c, a.pool)

	rows, err := db.Query(c, query.GetAllAuthors, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("get all authors: %w", err)
	}

	authorRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Author])
	if err != nil {
		return nil, 0, fmt.Errorf("scan author: %w", err)
	}

	authors := toDomainAuthors(authorRows)

	var total int
	err = db.QueryRow(c, query.GetTotalAuthors).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count authors: %w", err)
	}

	return authors, total, nil
}

func (a *AuthorRepository) GetOrCreate(ctx context.Context, author IRepository.IAuthor) (*entity.Author, error) {
	db := txFromContext(ctx, a.pool)

	row, err := db.Query(ctx, query.GetOrCreateAuthor, author.Name())
	if err != nil {
		return nil, fmt.Errorf("create author: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("scan author: %w", err)
		}

		row, err = db.Query(ctx, query.GetAuthorsByName, author.Name())
		if err != nil {
			return nil, fmt.Errorf("get author: %w", err)
		}
		defer row.Close()

		authorRow, err = pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
		if err != nil {
			return nil, fmt.Errorf("scan author: %w", err)
		}
	}

	newAuthor := toDomainAuthor(&authorRow)

	return newAuthor, nil
}

func (a *AuthorRepository) HasBooks(c context.Context, authorID int) (bool, error) {
	db := txFromContext(c, a.pool)

	var hasBooks bool
	err := db.QueryRow(c, query.ExistsBooksByAuthorID, authorID).Scan(&hasBooks)
	if err != nil {
		return false, fmt.Errorf("check author books: %w", err)
	}

	return hasBooks, nil
}

func (a *AuthorRepository) Delete(c context.Context, authorID int) error {
	db := txFromContext(c, a.pool)

	tag, err := db.Exec(c, query.DeleteAuthor, authorID)
	if err != nil {
		return fmt.Errorf("delete author: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entity.ErrAuthorNotFound
	}

	return nil
}

func (a *AuthorRepository) Update(c context.Context, author IRepository.IAuthor) (*entity.Author, error) {
	db := txFromContext(c, a.pool)

	row, err := db.Query(c, query.UpdateAuthor, author.Name(), author.ID())
	if err != nil {
		return nil, fmt.Errorf("update author: %w", err)
	}

	authorRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Author])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, app.ErrAuthorAlreadyExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrAuthorNotFound
		}
		return nil, fmt.Errorf("update author: %w", err)
	}

	newAuthor := toDomainAuthor(&authorRow)

	return newAuthor, nil
}

func toDomainAuthors(authorRows []Author) []*entity.Author {
	authors := make([]*entity.Author, 0, len(authorRows))

	for _, row := range authorRows {
		author := toDomainAuthor(&row)
		authors = append(authors, author)
	}

	return authors
}

func toDomainAuthor(a *Author) *entity.Author {
	return entity.HydrateAuthor(
		a.ID,
		a.Name,
	)
}
