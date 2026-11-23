package interfaces

import (
	"context"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/domain/entities"
	"virtual-library/catalog-service/internal/httpapi/pagination"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, command *command.CreateAuthorCommand) (*entities.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	UpdateAuthor(ctx context.Context, command *command.UpdateAuthorCommand) (*entities.Author, error)
	GetAuthorByID(ctx context.Context, id int) (*entities.Author, error)
	GetAllAuthors(ctx context.Context, pagination *pagination.Pagination) ([]*entities.Author, int, error)
}
