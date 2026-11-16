package service

import (
	"context"
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
	//author, err := s.repoAuthor.GetByID(ctx, book.AuthorID)
	//if errors.Is(err, domain.ErrAuthorNotFound) {
	//	author = &domain.Author{}
	//	return s.CreateBookWithNewAuthor(ctx, book, *author)
	//}
	//if err != nil {
	//	return domain.Book{}, err
	//}
	//
	return s.repoBook.Create(ctx, book)
}

func (s *BookService) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	return s.repoBook.GetByID(ctx, id)
}

func (s *BookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return s.repoBook.GetAll(ctx)
}

func (s *BookService) Delete(ctx context.Context, id int) error {
	return s.repoBook.Delete(ctx, id)
}

func (s *BookService) Update(ctx context.Context, book domain.Book) (domain.Book, error) {

	return s.repoBook.Update(ctx, book)
}

func (s *BookService) CreateBookWithNewAuthor(ctx context.Context, book domain.Book, author domain.Author) (domain.Book, error) {
	var createdBook domain.Book
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		createdAuthor, err := s.repoAuthor.Create(ctx, author)
		if err != nil {
			return fmt.Errorf("create author: %w", err)
		}

		book.AuthorID = createdAuthor.ID
		createdBook, err = s.repoBook.Create(ctx, book)
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
