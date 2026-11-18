package service

import (
	"context"
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
func (svc *AuthorService) GetByID(ctx context.Context, id int) (*domain.Author, error) {
	return svc.repo.GetByID(ctx, id)
}
func (svc *AuthorService) Update(ctx context.Context, id int, patch domain.AuthorPatch) (domain.Author, error) {
	if patch.Name == nil {
		return domain.Author{}, domain.ErrEmptyAuthorPatch
	}

	return svc.repo.Update(ctx, id, patch)
}
