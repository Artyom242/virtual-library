package domain

import "strings"

type Author struct {
	id   int
	name string
}

func NewAuthor(name string) (Author, error) {
	if strings.TrimSpace(name) == "" {
		return Author{}, ErrInvalidAuthorData
	}

	return Author{name: name}, nil
}

func (a Author) ID() int {
	return a.id
}

func (a Author) Name() string {
	return a.name
}

func HydrateAuthor(id int, name string) Author {
	return Author{
		id:   id,
		name: name,
	}
}

type AuthorPatch struct {
	Name *string
}
