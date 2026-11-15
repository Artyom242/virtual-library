package catalog

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *BookService
}

type CreateBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required,min=1"`
}

func NewBookHandler(service *BookService) *BookHandler {
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
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	} else if book == nil {
		c.JSON(404, gin.H{
			"message": "Book not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": book,
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
	request := CreateBookRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Book Data",
		})
		return
	}

	book := Book{
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
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
		return
	}

	c.Status(204)
	return
}
