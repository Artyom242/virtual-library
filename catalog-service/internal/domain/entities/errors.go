package entities

import "errors"

// books
var (
	ErrBookNotFound           = errors.New("book not found")
	ErrInvalidBookID          = errors.New("invalid book id")
	ErrInvalidBookData        = errors.New("invalid book data")
	ErrInvalidBookDescription = errors.New("invalid book description")
)

// authors
var (
	ErrAuthorNotFound    = errors.New("author not found")
	ErrInvalidAuthorID   = errors.New("invalid author id")
	ErrInvalidAuthorData = errors.New("invalid author data")
)

var ErrValidationError = errors.New("validation error")
