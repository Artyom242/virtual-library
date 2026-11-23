package entities

import (
	"fmt"
	"strings"
	"virtual-library/catalog-service/internal/app"
)

type Book struct {
	id          int
	title       string
	year        int
	authorID    int
	description *string
}

func (b *Book) ID() int {
	return b.id
}
func (b *Book) Title() string {
	return b.title
}
func (b *Book) Year() int {
	return b.year
}
func (b *Book) AuthorID() int {
	return b.authorID
}
func (b *Book) Description() *string {
	return b.description
}
func (b *Book) validate() error {
	if strings.Trim(b.title, " ") == "" {
		return fmt.Errorf("%w: title is required", ErrValidationError)
	}
	if b.year <= 1 {
		return fmt.Errorf("%w: year must be greater than 1", ErrValidationError)
	}
	if b.authorID <= 1 {
		return fmt.Errorf("%w: author_id must be greater than 1", ErrValidationError)
	}
	if b.description != nil && strings.Trim(*b.description, " ") == "" {
		return fmt.Errorf("%w: description must not be empty", ErrValidationError)
	}

	return nil
}

func NewBook(title string, year int, authorID int, description *string) (*Book, error) {
	book := Book{
		title:       title,
		year:        year,
		authorID:    authorID,
		description: description,
	}

	err := book.validate()
	if err != nil {
		return nil, err
	}

	return &book, nil
}
func HydrateBook(id int, title string, year int, authorID int, description *string) *Book {
	return &Book{
		id:          id,
		title:       title,
		year:        year,
		authorID:    authorID,
		description: description,
	}
}

type BookChanges struct {
	Title       *string
	Year        *int
	AuthorID    *int
	Description *string
}

func (b *Book) ApplyChanges(changes BookChanges) error {
	updated := *b

	if changes.isEmpty() {
		return app.ErrEmptyBookPatch
	}

	if changes.Title != nil {
		updated.title = *changes.Title
	}
	if changes.Year != nil {
		updated.year = *changes.Year
	}
	if changes.AuthorID != nil {
		updated.authorID = *changes.AuthorID
	}
	if changes.Description != nil {
		updated.description = changes.Description
	}

	if err := updated.validate(); err != nil {
		return err
	}

	*b = updated
	return nil
}

func (bc BookChanges) isEmpty() bool {
	if bc.Title == nil && bc.Year == nil && bc.AuthorID == nil && bc.Description == nil {
		return true
	}
	return false
}
