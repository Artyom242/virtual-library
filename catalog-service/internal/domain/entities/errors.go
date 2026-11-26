package entities

import "errors"

// books
var (
	ErrBookNotFound           = errors.New("book not found")
	ErrInvalidBookID          = errors.New("invalid book id")
	ErrInvalidBookData        = errors.New("invalid book data")
	ErrInvalidBookDescription = errors.New("invalid book description")
	ErrFailed                 = errors.New("invalid book description")
)

// authors
var (
	ErrAuthorNotFound    = errors.New("author not found")
	ErrInvalidAuthorID   = errors.New("invalid author id")
	ErrInvalidAuthorData = errors.New("invalid author data")
)

// categories
var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrInvalidCategoryID     = errors.New("invalid category id")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotUnique     = errors.New("categories are not unique")
)

var ErrValidationError = errors.New("validation error")
