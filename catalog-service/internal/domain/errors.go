package domain

import "errors"

// books
var (
	ErrBookNotFound           = errors.New("book not found")
	ErrInvalidBookID          = errors.New("invalid book id")
	ErrInvalidBookData        = errors.New("invalid book data")
	ErrEmptyBookPatch         = errors.New("empty book patch")
	ErrMapBookToDomain        = errors.New("failed to map book to domain")
	ErrInvalidBookDescription = errors.New("empty book description")
)

// authors
var (
	ErrAuthorNotFound      = errors.New("author not found")
	ErrInvalidAuthorID     = errors.New("invalid author id")
	ErrInvalidAuthorData   = errors.New("invalid author data")
	ErrEmptyAuthorPatch    = errors.New("empty author patch")
	ErrAuthorHasBooks      = errors.New("author has books")
	ErrAuthorAlreadyExists = errors.New("author already exists")
	ErrMapAuthorToDomain   = errors.New("failed to map author to domain")
)

// other
var (
	ErrValidationError   = errors.New("validation error")
	ErrInternalServer    = errors.New("internal server error")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)
