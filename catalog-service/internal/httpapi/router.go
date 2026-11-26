package httpapi

import (
	"virtual-library/catalog-service/internal/app/interface"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/httpapi/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	bookService _interface.BookService,
	authorService _interface.AuthorService,
	categoryService _interface.CategoryService,
) *gin.Engine {
	router := gin.Default()
	mainHandler := handler.NewCombinedHandler(bookService, authorService, categoryService)
	generated.RegisterHandlers(router, mainHandler)

	return router
}
