package httpapi

import (
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(bookService *service.BookService, authorService *service.AuthorService) *gin.Engine {
	router := gin.Default()

	handler := NewCombinedHandler(bookService, authorService)
	generated.RegisterHandlers(router, handler)

	return router
}
