package domain

import "errors"

var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrBookNotFound   = errors.New("book not found")
)
