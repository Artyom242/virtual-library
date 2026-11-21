package domain

import (
	"fmt"
	"strings"
)

type Book struct {
	id          int
	title       string
	year        int
	authorID    int
	description *string
}

func (b Book) ID() int {
	return b.id
}
func (b Book) Title() string {
	return b.title
}
func (b Book) Year() int {
	return b.year
}
func (b Book) AuthorID() int {
	return b.authorID
}
func (b Book) Description() *string {
	return b.description
}

func NewBook(title string, year int, authorID int, description *string) (Book, error) {
	if title == "" {
		return Book{}, fmt.Errorf("%w title must not be empty", ErrValidationError)
	}
	if err := validateOptionalText(description); err != nil {
		return Book{}, fmt.Errorf("%w : description must not be empty", ErrInvalidBookDescription)
	}

	return Book{
		title:       title,
		year:        year,
		authorID:    authorID,
		description: description,
	}, nil
}

func HydrateBook(id int, title string, year int, authorID int, description *string) Book {
	return Book{
		id:          id,
		title:       title,
		year:        year,
		authorID:    authorID,
		description: description,
	}
}

type BookPatch struct {
	Title       *string
	AuthorID    *int
	Year        *int
	Description *string
}

func validateOptionalText(value *string) error {
	if value != nil && strings.TrimSpace(*value) == "" {
		return ErrInvalidBookDescription
	}

	return nil
}

func (p BookPatch) IsEmpty() bool {
	return p.Title == nil &&
		p.AuthorID == nil &&
		p.Year == nil &&
		p.Description == nil
}

func (p BookPatch) Validate() error {
	if p.Description != nil && strings.TrimSpace(*p.Description) == "" {
		return ErrInvalidBookDescription
	}

	return nil
}

type CreateBookWithAuthor struct {
	Title       string
	Year        int
	Description *string
}
