package handler

import (
	"virtual-library/catalog-service/internal/app/interface"
	"virtual-library/catalog-service/internal/httpapi/generated"
)

type CombinedHandler struct {
	*BookHandler
	*AuthorHandler
	*CategoryHandler
}

var _ generated.ServerInterface = (*CombinedHandler)(nil)

func NewCombinedHandler(
	bookService _interface.BookService,
	authorService _interface.AuthorService,
	categoryService _interface.CategoryService,
) *CombinedHandler {
	return &CombinedHandler{
		BookHandler:     NewBookHandler(bookService),
		AuthorHandler:   NewAuthorHandler(authorService),
		CategoryHandler: NewCategoryHandler(categoryService),
	}
}
