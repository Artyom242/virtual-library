package httpapi

import "github.com/gin-gonic/gin"

func NewRouter(bookHandler *BookHandler, authorHandler *AuthorHandler) *gin.Engine {
	router := gin.Default()

	books := router.Group("/books")
	{
		books.GET("", bookHandler.GetAll)
		books.GET("/:id", bookHandler.GetByID)
		books.POST("", bookHandler.Create)
		books.POST("/with-author", bookHandler.CreateBookWithNewAuthor)
		books.DELETE("/:id", bookHandler.DeleteBook)
		books.PATCH("/:id", bookHandler.Update)
	}

	authors := router.Group("/authors")
	{
		authors.GET("", authorHandler.GetAll)
		authors.GET("/:id", authorHandler.GetByID)
		authors.POST("", authorHandler.Create)
		authors.DELETE("/:id", authorHandler.Delete)
		authors.PATCH("/:id", authorHandler.Update)
	}

	return router
}
