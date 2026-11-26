package postgres

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/domain/repository"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	"virtual-library/catalog-service/internal/infra/db/postgres/query"
	IRepository "virtual-library/catalog-service/internal/infra/db/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

var _ repository.BookRepository = (*BookRepository)(nil)

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

type BookWithCategories struct {
	Book
	CategoryIDs []int
}

func (b *BookRepository) Create(ctx context.Context, book IRepository.IBook) (*entity.Book, error) {
	db := txFromContext(ctx, b.pool)

	row, err := db.Query(ctx, query.CreateBook, book.Title(), book.Year(), book.AuthorID(), book.Description())
	if err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if err = mapPostgresError(err); err != nil {
			return nil, fmt.Errorf("error map pgs error: %w", err)
		}
		return nil, err
	}

	err = b.SetCategories(ctx, bookRow.ID, book.Categories())
	if err != nil {
		return nil, fmt.Errorf("create book categories: %w", err)
	}

	bookWithCategories := BookWithCategories{
		Book:        bookRow,
		CategoryIDs: book.Categories(),
	}
	newBook := toDomainBook(&bookWithCategories)

	return newBook, nil
}

func (b *BookRepository) GetByID(ctx context.Context, bookID int) (*entity.Book, error) {
	db := txFromContext(ctx, b.pool)

	row, err := db.Query(ctx, query.GetBookById, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book by id: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrBookNotFound
		}
		return nil, fmt.Errorf("scan book: %w", err)
	}

	categoryIDs, err := b.GetCategoriesByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book category ids: %w", err)
	}

	bookWithCategories := BookWithCategories{
		Book:        bookRow,
		CategoryIDs: categoryIDs,
	}

	newBook := toDomainBook(&bookWithCategories)

	return newBook, nil
}

func (b *BookRepository) GetAll(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Book, int, error) {
	db := txFromContext(ctx, b.pool)

	rows, err := db.Query(ctx, query.GetAllBooks, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("get all books: %w", err)
	}

	bookRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, 0, fmt.Errorf("scan books: %w", err)
	}

	bookIDs := make([]int, 0, len(bookRows))
	for _, book := range bookRows {
		bookIDs = append(bookIDs, book.ID)
	}

	categoriesByBook, err := b.GetCategoriesByBookIDs(ctx, bookIDs)
	if err != nil {
		return nil, 0, err
	}

	books := make([]*entity.Book, 0, len(bookRows))
	for _, book := range bookRows {
		bookWithCategories := BookWithCategories{
			Book:        book,
			CategoryIDs: categoriesByBook[book.ID],
		}
		books = append(books, toDomainBook(&bookWithCategories))
	}

	var total int
	err = db.QueryRow(ctx, query.GetTotalBooks).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count books: %w", err)
	}

	return books, total, nil
}

func (b *BookRepository) Delete(ctx context.Context, bookID int) error {
	db := txFromContext(ctx, b.pool)

	tag, err := db.Exec(ctx, query.DeleteBook, bookID)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entity.ErrBookNotFound
	}

	err = b.deleteBookCategories(ctx, bookID)
	if err != nil {
		return fmt.Errorf("delete book categories: %w", err)
	}

	return nil
}

func (b *BookRepository) Update(ctx context.Context, changes IRepository.IBook) (*entity.Book, error) {
	db := txFromContext(ctx, b.pool)

	row, err := db.Query(ctx, query.UpdateBook, changes.Title(), changes.Year(), changes.AuthorID(), changes.Description(), changes.ID())
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	bookRow, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrBookNotFound
		}
		if err = mapPostgresError(err); err != nil {
			return nil, fmt.Errorf("error map pgs error: %w", err)
		}

		return nil, fmt.Errorf("update book: %w", err)
	}

	// Category
	if changes.Categories() != nil {
		if err = b.deleteBookCategories(ctx, bookRow.ID); err != nil {
			return nil, fmt.Errorf("update book categories: %w", err)
		}
		if err = b.SetCategories(ctx, bookRow.ID, changes.Categories()); err != nil {
			return nil, fmt.Errorf("update book categories: %w", err)
		}
	}

	bookWithCategories := BookWithCategories{
		Book:        bookRow,
		CategoryIDs: changes.Categories(),
	}

	newBook := toDomainBook(&bookWithCategories)

	return newBook, nil
}

func (b *BookRepository) deleteBookCategories(ctx context.Context, bookID int) error {
	db := txFromContext(ctx, b.pool)

	_, err := db.Exec(ctx, query.DeleteBookCategories, bookID)
	if err != nil {
		return fmt.Errorf("delete book categories: %w", err)
	}

	return nil
}

func (b *BookRepository) SetCategories(ctx context.Context, bookID int, categories []int) error {
	db := txFromContext(ctx, b.pool)
	batch := &pgx.Batch{}

	for _, categoryID := range categories {
		batch.Queue(query.SetCategories, bookID, categoryID)
	}

	batchLen := batch.Len()
	if batchLen > 0 {
		br := db.SendBatch(ctx, batch)

		for range batchLen {
			_, err := br.Exec()
			if err != nil {
				if err = mapPostgresError(err); err != nil {
					return fmt.Errorf("error map pgs error: %w", err)
				}

				closeErr := br.Close()
				return errors.Join(
					fmt.Errorf("create book category relation: %w", err),
					closeErr,
				)
			}
		}

		if err := br.Close(); err != nil {
			return fmt.Errorf("close book categories batch: %w", err)
		}
	}

	return nil
}

func (b *BookRepository) GetCategoriesByBookID(ctx context.Context, bookID int) ([]int, error) {
	db := txFromContext(ctx, b.pool)

	rows, err := db.Query(ctx, query.GetCategoriesByBookID, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book_categories: %w", err)
	}

	var categories []int

	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			return nil, fmt.Errorf("scan book_categories: %w", err)
		}

		categories = append(categories, categoryID)
	}

	return categories, nil
}

func (b *BookRepository) GetCategoriesByBookIDs(ctx context.Context, bookIDs []int) (map[int][]int, error) {
	result := make(map[int][]int, len(bookIDs))

	if len(bookIDs) == 0 {
		return result, nil
	}

	for _, id := range bookIDs {
		result[id] = make([]int, 0)
	}

	db := txFromContext(ctx, b.pool)
	rows, err := db.Query(ctx, query.GetCategoriesByBookIDs, bookIDs)
	if err != nil {
		return nil, fmt.Errorf("get categories by book ids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bookID, categoryID int

		if err := rows.Scan(&bookID, &categoryID); err != nil {
			return nil, fmt.Errorf("scan book category: %w", err)
		}

		result[bookID] = append(result[bookID], categoryID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func toDomainBook(b *BookWithCategories) *entity.Book {
	var description *string

	if b.Description.Valid {
		description = &b.Description.String
	}

	return entity.HydrateBook(
		b.ID,
		b.Title,
		b.Year,
		b.AuthorID,
		description,
		&b.CategoryIDs,
	)
}
