package service

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interface"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/domain/repository"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type BookService struct {
	tx         app.Transactor
	repoBook   repository.BookRepository
	repoAuthor repository.AuthorRepository
}

func NewService(
	tx app.Transactor,
	repoBook repository.BookRepository,
	repoAuthor repository.AuthorRepository,
) _interface.BookService {
	return &BookService{
		tx:         tx,
		repoBook:   repoBook,
		repoAuthor: repoAuthor,
	}
}

func (s *BookService) CreateBook(ctx context.Context, bookCommand *command.CreateBookCommand) (*entity.Book, error) {
	var newBook *entity.Book
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		_, err := s.repoAuthor.GetByID(ctx, bookCommand.AuthorID)
		if err != nil {
			if errors.Is(err, entity.ErrAuthorNotFound) {
				return entity.ErrAuthorNotFound
			}
			return err
		}

		book, err := entity.NewBook(
			bookCommand.Title,
			bookCommand.Year,
			bookCommand.AuthorID,
			bookCommand.Description,
			bookCommand.CategoryIDs,
		)
		if err != nil {
			return err
		}

		newBook, err = s.repoBook.Create(ctx, book)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newBook, nil
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*entity.Book, error) {
	return s.repoBook.GetByID(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Book, int, error) {
	return s.repoBook.GetAll(ctx, pagination)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.repoBook.Delete(ctx, id)
}

func (s *BookService) UpdateBook(ctx context.Context, bookCommand *command.UpdateBookCommand) (*entity.Book, error) {
	var updatedBook *entity.Book
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		book, err := s.GetBookByID(ctx, bookCommand.ID)
		if err != nil {
			return err
		}

		if bookCommand.AuthorID != nil {
			_, err = s.repoAuthor.GetByID(ctx, *bookCommand.AuthorID)
			if err != nil {
				return err
			}
		}

		bookChanges := entity.BookChanges{
			Title:         bookCommand.Title,
			Year:          bookCommand.Year,
			AuthorID:      bookCommand.AuthorID,
			Description:   bookCommand.Description,
			CategoriesIDs: bookCommand.CategoryIDs,
		}

		err = book.ApplyChanges(bookChanges)
		if err != nil {
			return err
		}

		updatedBook, err = s.repoBook.Update(ctx, book)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedBook, nil
}

func (s *BookService) CreateBookWithNewAuthor(ctx context.Context, commandBook *command.CreateBookCommand, commandAuthor *command.CreateAuthorCommand) (*entity.Book, error) {
	var createdBook *entity.Book

	err := s.tx.Do(ctx, func(ctx context.Context) error {
		createAuthor, err := entity.NewAuthor(commandAuthor.Name)
		if err != nil {
			return fmt.Errorf("create author: %w", err)
		}

		createdAuthor, err := s.repoAuthor.GetOrCreate(ctx, createAuthor)
		if err != nil {
			return fmt.Errorf("created author: %w", err)
		}

		newBook, err := entity.NewBook(
			commandBook.Title,
			commandBook.Year,
			createdAuthor.ID(),
			commandBook.Description,
			commandBook.CategoryIDs,
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
		return nil, fmt.Errorf("create book with new author: %w", err)
	}

	return createdBook, nil
}
