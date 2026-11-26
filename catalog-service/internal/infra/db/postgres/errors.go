package postgres

import (
	"errors"
	"virtual-library/catalog-service/internal/domain/entities"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL unique_violation code
const pgUniqueViolation = "23505"

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
	case "fk_book_categories_books":
		return entities.ErrBookNotFound
	case "fk_book_categories_categories":
		return entities.ErrCategoryNotFound
	}

	return nil
}
