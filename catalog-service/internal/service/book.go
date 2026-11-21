package service

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/domain"
)

type BookService struct {
	tx         domain.Transactor
	repoBook   domain.BookRepository
	repoAuthor domain.AuthorRepository
}

func NewService(
	tx domain.Transactor,
	repoBook domain.BookRepository,
	repoAuthor domain.AuthorRepository,
) *BookService {
	return &BookService{
		tx:         tx,
		repoBook:   repoBook,
		repoAuthor: repoAuthor,
	}
}

func (s *BookService) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	_, err := s.repoAuthor.GetByID(ctx, book.AuthorID())
	if errors.Is(err, domain.ErrAuthorNotFound) {
		return domain.Book{}, domain.ErrAuthorNotFound
	}
	if err != nil {
		return domain.Book{}, err
	}

	return s.repoBook.Create(ctx, book)
}

func (s *BookService) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	return s.repoBook.GetByID(ctx, id)
}

func (s *BookService) GetAll(ctx context.Context, page, limit int) ([]domain.Book, int, error) {
	return s.repoBook.GetAll(ctx, page, limit)
}

func (s *BookService) Delete(ctx context.Context, id int) error {
	return s.repoBook.Delete(ctx, id)
}

func (s *BookService) Update(ctx context.Context, id int, patch domain.BookPatch) (domain.Book, error) {
	if patch.IsEmpty() {
		return domain.Book{}, domain.ErrEmptyBookPatch
	}

	var updatedBook domain.Book
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		_, err := s.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if patch.AuthorID != nil {
			_, err = s.repoAuthor.GetByID(ctx, *patch.AuthorID)
			if err != nil {
				return err
			}
		}

		if err = patch.Validate(); err != nil {
			return err
		}

		updatedBook, err = s.repoBook.Update(ctx, id, patch)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return updatedBook, err
	}

	return updatedBook, nil
}

func (s *BookService) CreateBookWithNewAuthor(ctx context.Context, book domain.CreateBookWithAuthor, author domain.Author) (domain.Book, error) {
	var createdBook domain.Book
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		createdAuthor, err := s.repoAuthor.GetOrCreate(ctx, author)
		if err != nil {
			return fmt.Errorf("create author: %w", err)
		}

		newBook, err := domain.NewBook(
			book.Title,
			book.Year,
			createdAuthor.ID(),
			book.Description,
		)
		if err != nil {
			return fmt.Errorf("create book entity: %w", err)
		}

		createdBook, err = s.repoBook.Create(ctx, newBook)
		if err != nil {

			return fmt.Errorf("create book: %w", err)
		}
		return nil
	})
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book with new author: %w", err)
	}
	return createdBook, nil
}
