package entities

import (
	"strings"
)

type Category struct {
	id   int
	name string
}

func (c *Category) ID() int {
	return c.id
}
func (c *Category) Name() string {
	return c.name
}
func (c *Category) validate() error {
	if strings.TrimSpace(c.name) == "" {
		return ErrValidationError
	}

	return nil
}

func NewCategory(name string) (*Category, error) {
	category := Category{
		name: name,
	}

	if err := category.validate(); err != nil {
		return &Category{}, err
	}

	return &category, nil
}
func HydrateCategory(id int, name string) *Category {
	return &Category{
		id:   id,
		name: name,
	}
}

type CategoryChanges struct {
	Name *string
}

func (c *Category) ApplyCategory(category CategoryChanges) (*Category, error) {
	updated := *c

	if category.Name != nil {
		updated.name = *category.Name
	}

	if err := updated.validate(); err != nil {
		return nil, err
	}

	return &updated, nil
}
