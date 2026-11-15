package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"virtual-library/catalog-service/internal/catalog"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *catalog.BookService
}

func NewBookHandler(service *catalog.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetByID(c *gin.Context) {
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Book Id",
		})
		return
	}

	book, err := h.service.GetByID(id)
	if errors.Is(err, catalog.ErrBookNotFound) {
		c.JSON(404, gin.H{
			"message": "Book not found",
		})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"data": book,
	})
	return
}
func (h *BookHandler) GetAll(c *gin.Context) {
	books, err := h.service.GetAll()

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": books,
	})
	return
}
func (h *BookHandler) CreateBook(c *gin.Context) {
	var request CreateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Book Data",
		})
		return
	}

	book := catalog.Book{
		Title:  request.Title,
		Author: request.Author,
		Year:   request.Year,
	}

	book, err := h.service.CreateBook(book)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "Book created",
		"data":    book,
	})
	return
}
func (h *BookHandler) DeleteBook(c *gin.Context) {
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Book Id",
		})
		return
	}

	err = h.service.Delete(id)
	if errors.Is(err, catalog.ErrBookNotFound) {
		c.JSON(404, gin.H{
			"message": "Book not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	}

	c.Status(204)
	return
}
