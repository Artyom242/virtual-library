package service

import (
	"context"
	"errors"
	"virtual-library/catalog-service/internal/domain"
)

type AuthorService struct {
	repo domain.AuthorRepository
}

func NewAuthorService(repo domain.AuthorRepository) *AuthorService {
	return &AuthorService{repo: repo}
}

func (svc *AuthorService) Create(ctx context.Context, author domain.Author) (domain.Author, error) {
	return svc.repo.Create(ctx, author)
}
func (svc *AuthorService) Delete(ctx context.Context, id int) error {
	return svc.repo.Delete(ctx, id)
}
func (svc *AuthorService) GetAll(ctx context.Context) ([]domain.Author, error) {
	return svc.repo.GetAll(ctx)
}
func (svc *AuthorService) GetById(ctx context.Context, id int) (*domain.Author, error) {
	return svc.repo.GetByID(ctx, id)
}
func (svc *AuthorService) Update(ctx context.Context, id int, patch domain.AuthorPatch) (domain.Author, error) {
	author, err := svc.repo.GetByID(ctx, id)

	if errors.Is(err, domain.ErrAuthorNotFound) {
		return domain.Author{}, domain.ErrAuthorNotFound
	}
	if err != nil {
		return domain.Author{}, err
	}

	if patch.Name != nil {
		author.Name = *patch.Name
	}

	return svc.repo.Update(ctx, *author)
}
