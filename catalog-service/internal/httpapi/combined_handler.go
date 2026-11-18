package httpapi

import (
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/service"
)

type CombinedHandler struct {
	*BookHandler
	*AuthorHandler
}

var _ generated.ServerInterface = (*CombinedHandler)(nil)

func NewCombinedHandler(bookService *service.BookService, authorService *service.AuthorService) *CombinedHandler {
	return &CombinedHandler{
		BookHandler:   NewBookHandler(bookService),
		AuthorHandler: NewAuthorHandler(authorService),
	}
}
