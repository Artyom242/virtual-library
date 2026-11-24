package entities

import (
	"strings"
)

type Author struct {
	id   int
	name string
}

func (a *Author) ID() int {
	return a.id
}
func (a *Author) Name() string {
	return a.name
}
func (a *Author) validate() error {
	if strings.TrimSpace(a.name) == "" {
		return ErrValidationError
	}

	return nil
}

func NewAuthor(name string) (*Author, error) {
	author := Author{
		name: name,
	}

	if err := author.validate(); err != nil {
		return &Author{}, err
	}

	return &author, nil
}
func HydrateAuthor(id int, name string) *Author {
	return &Author{
		id:   id,
		name: name,
	}
}

type AuthorChanges struct {
	Name *string
}

func (a *Author) ApplyAuthor(author AuthorChanges) (*Author, error) {
	updated := *a

	if author.Name != nil {
		updated.name = *author.Name
	}

	if err := updated.validate(); err != nil {
		return nil, err
	}

	return &updated, nil
}
