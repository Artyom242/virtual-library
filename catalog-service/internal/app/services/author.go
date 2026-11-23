package services

import (
	"context"
	"fmt"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/domain/repositories"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type AuthorService struct {
	tx   app.Transactor
	repo repositories.AuthorRepository
}

func NewAuthorService(tx app.Transactor, repo repositories.AuthorRepository) interfaces.AuthorService {
	return &AuthorService{tx: tx, repo: repo}
}

func (svc *AuthorService) CreateAuthor(ctx context.Context, command *command.CreateAuthorCommand) (*entities.Author, error) {
	author, err := entities.NewAuthor(
		command.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("create author: %w", err)
	}

	return svc.repo.GetOrCreate(ctx, author)
}

func (svc *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	return svc.tx.Do(ctx, func(ctx context.Context) error {
		hasBooks, err := svc.repo.HasBooks(ctx, id)
		if err != nil {
			return fmt.Errorf("check author books: %w", err)
		}
		if hasBooks {
			return app.ErrAuthorHasBooks
		}

		return svc.repo.Delete(ctx, id)
	})
}

func (svc *AuthorService) GetAllAuthors(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Author, int, error) {
	return svc.repo.GetAll(ctx, pagination)
}

func (svc *AuthorService) GetAuthorByID(ctx context.Context, id int) (*entities.Author, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *AuthorService) UpdateAuthor(ctx context.Context, authorCommand *command.UpdateAuthorCommand) (*entities.Author, error) {
	author, err := svc.repo.GetByID(ctx, authorCommand.ID)
	if err != nil {
		return nil, fmt.Errorf("update author: %w", err)
	}

	authorChanged := entities.AuthorChanges{
		Name: authorCommand.Name,
	}

	author, err = author.ApplyAuthor(authorChanged)
	if err != nil {
		return nil, fmt.Errorf("update author: %w", err)
	}

	return svc.repo.Update(ctx, author)
}
