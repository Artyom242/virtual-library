package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/domain/repositories"
	"virtual-library/catalog-service/internal/httpapi/pagination"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

type Category struct {
	ID   int
	Name string
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	if pool == nil {
		return &CategoryRepository{}
	}
	return &CategoryRepository{pool: pool}
}

var _ repositories.CategoryRepository = (*CategoryRepository)(nil)

func (c *CategoryRepository) Create(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	db := txFromContext(ctx, c.pool)

	query := "INSERT INTO categories(name) VALUES($1) returning id, name"
	rows, err := db.Query(ctx, query, category.Name())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, entities.ErrCategoryAlreadyExists
		}
		return nil, fmt.Errorf("create category: %w", err)
	}

	categoryRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Category])
	if err != nil {
		return nil, fmt.Errorf("scan category: %w", err)
	}

	newCategory := toDomainCategory(&categoryRow)
	return newCategory, nil
}

func (c *CategoryRepository) GetByID(ctx context.Context, categoryID int) (*entities.Category, error) {
	db := txFromContext(ctx, c.pool)

	query := "SELECT categories.id, categories.name FROM categories WHERE id = $1"
	row, err := db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}

	categoryRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("scan category: %w", err)
	}

	category := toDomainCategory(&categoryRow)
	return category, nil
}

func (c *CategoryRepository) GetAll(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Category, int, error) {
	db := txFromContext(ctx, c.pool)

	query := `
		SELECT categories.id, categories.name 
		FROM categories
		ORDER BY categories.id
		LIMIT $1
		OFFSET $2
		`
	rows, err := db.Query(ctx, query, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("get all categories: %w", err)
	}

	categoryRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Category])
	if err != nil {
		return nil, 0, fmt.Errorf("scan category: %w", err)
	}

	categories := toDomainCategories(categoryRows)

	var total int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM categories").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count categories: %w", err)
	}

	return categories, total, nil
}

func (c *CategoryRepository) Delete(ctx context.Context, categoryID int) error {
	db := txFromContext(ctx, c.pool)
	query := "DELETE FROM categories WHERE id = $1"
	tag, err := db.Exec(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.ErrCategoryNotFound
	}

	return nil
}

func (c *CategoryRepository) Update(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	db := txFromContext(ctx, c.pool)

	query := `UPDATE categories 
		SET name = COALESCE($1, name) 
		WHERE id = $2 
		RETURNING id, name`
	row, err := db.Query(ctx, query, category.Name(), category.ID())
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	categoryRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Category])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, entities.ErrCategoryAlreadyExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("update category: %w", err)
	}

	newCategory := toDomainCategory(&categoryRow)
	return newCategory, nil
}

func toDomainCategories(categoryRows []Category) []*entities.Category {
	categories := make([]*entities.Category, 0, len(categoryRows))

	for _, row := range categoryRows {
		category := toDomainCategory(&row)
		categories = append(categories, category)
	}

	return categories
}

func toDomainCategory(c *Category) *entities.Category {
	return entities.HydrateCategory(
		c.ID,
		c.Name,
	)
}
