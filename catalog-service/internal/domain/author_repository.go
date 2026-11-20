package domain

import (
	"context"
)

type AuthorRepository interface {
	Create(c context.Context, author Author) (Author, error)
	GetByID(c context.Context, authorID int) (*Author, error)
	GetAll(c context.Context, page int, limit int) ([]Author, int, error)
	Delete(c context.Context, authorID int) error
	Update(c context.Context, id int, author AuthorPatch) (Author, error)
	GetOrCreate(ctx context.Context, author Author) (Author, error)
	HasBooks(c context.Context, authorID int) (bool, error)
}
