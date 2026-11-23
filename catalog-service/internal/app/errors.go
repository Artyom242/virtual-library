package app

import "errors"

var (
	ErrEmptyBookPatch      = errors.New("empty book patch")
	ErrEmptyAuthorPatch    = errors.New("empty author patch")
	ErrAuthorHasBooks      = errors.New("author has books")
	ErrAuthorAlreadyExists = errors.New("author already exists")
)
