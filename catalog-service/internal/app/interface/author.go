package _interface

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, command *command.CreateAuthorCommand) (*entity.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	UpdateAuthor(ctx context.Context, command *command.UpdateAuthorCommand) (*entity.Author, error)
	GetAuthorByID(ctx context.Context, id int) (*entity.Author, error)
	GetAllAuthors(ctx context.Context, pagination *pagination.Pagination) ([]*entity.Author, int, error)
}
