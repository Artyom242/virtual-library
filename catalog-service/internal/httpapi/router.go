package httpapi

import "github.com/gin-gonic/gin"

func NewRouter(handler *BookHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/books/:id", handler.GetByID)
	router.GET("/books", handler.GetAll)
	router.POST("/books", handler.CreateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	return router
}
