package services

import (
	"context"
	"errors"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/domain/repositories"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type BookService struct {
	tx         app.Transactor
	repoBook   repositories.BookRepository
	repoAuthor repositories.AuthorRepository
}

func NewService(
	tx app.Transactor,
	repoBook repositories.BookRepository,
	repoAuthor repositories.AuthorRepository,
) interfaces.BookService {
	return &BookService{
		tx:         tx,
		repoBook:   repoBook,
		repoAuthor: repoAuthor,
	}
}

func (s *BookService) CreateBook(ctx context.Context, bookCommand *command.CreateBookCommand) (*entities.Book, error) {
	_, err := s.repoAuthor.GetByID(ctx, bookCommand.AuthorID)
	if errors.Is(err, entities.ErrAuthorNotFound) {
		return nil, entities.ErrAuthorNotFound
	}
	if err != nil {
		return nil, err
	}

	newBook, err := entities.NewBook(
		bookCommand.Title,
		bookCommand.Year,
		bookCommand.AuthorID,
		bookCommand.Description,
	)
	if err != nil {
		return nil, err
	}

	return s.repoBook.Create(ctx, newBook)
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*entities.Book, error) {
	return s.repoBook.GetByID(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Book, int, error) {
	return s.repoBook.GetAll(ctx, pagination)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.repoBook.Delete(ctx, id)
}

func (s *BookService) UpdateBook(ctx context.Context, bookCommand *command.UpdateBookCommand) (*entities.Book, error) {
	var updatedBook *entities.Book
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

		bookChanges := entities.BookChanges{
			Title:       bookCommand.Title,
			Year:        bookCommand.Year,
			AuthorID:    bookCommand.AuthorID,
			Description: bookCommand.Description,
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

func (s *BookService) CreateBookWithNewAuthor(ctx context.Context, commandBook *command.CreateBookCommand, commandAuthor *command.CreateAuthorCommand) (*entities.Book, error) {
	var createdBook *entities.Book

	err := s.tx.Do(ctx, func(ctx context.Context) error {
		createAuthor, err := entities.NewAuthor(commandAuthor.Name)
		if err != nil {
			return fmt.Errorf("create author: %w", err)
		}

		createdAuthor, err := s.repoAuthor.GetOrCreate(ctx, createAuthor)
		if err != nil {
			return fmt.Errorf("created author: %w", err)
		}

		newBook, err := entities.NewBook(
			commandBook.Title,
			commandBook.Year,
			createdAuthor.ID(),
			commandBook.Description,
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
