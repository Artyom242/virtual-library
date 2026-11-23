package handler

import (
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/httpapi/generated"
)

type CombinedHandler struct {
	*BookHandler
	*AuthorHandler
}

var _ generated.ServerInterface = (*CombinedHandler)(nil)

func NewCombinedHandler(bookService interfaces.BookService, authorService interfaces.AuthorService) *CombinedHandler {
	return &CombinedHandler{
		BookHandler:   NewBookHandler(bookService),
		AuthorHandler: NewAuthorHandler(authorService),
	}
}
