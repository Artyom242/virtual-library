package httpapi

import (
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/httpapi/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	bookService interfaces.BookService,
	authorService interfaces.AuthorService,
	categoryService interfaces.CategoryService,
) *gin.Engine {
	router := gin.Default()
	mainHandler := handler.NewCombinedHandler(bookService, authorService, categoryService)
	generated.RegisterHandlers(router, mainHandler)

	return router
}
