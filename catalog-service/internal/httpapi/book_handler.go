package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/dto"
	"virtual-library/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Id",
		})
		return
	}

	book, err := h.service.GetByID(ctx, id)
	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}
func (h *BookHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	books, err := h.service.GetAll(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})
}
func (h *BookHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Data",
		})
		return
	}

	book := domain.Book{
		Title:    request.Title,
		AuthorID: request.AuthorID,
		Year:     request.Year,
	}

	book, err := h.service.Create(ctx, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created",
		"data":    book,
	})
}
func (h *BookHandler) CreateBookWithNewAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateBookWithAuthorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Data",
		})
		return
	}

	book := domain.Book{
		Title: request.Title,
		Year:  request.Year,
	}

	author := domain.Author{
		Name: request.AuthorName,
	}

	createdBook, err := h.service.CreateBookWithNewAuthor(ctx, book, author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created",
		"data":    createdBook,
	})
}
func (h *BookHandler) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Id",
		})
		return
	}

	err = h.service.Delete(ctx, id)
	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	}

	c.Status(http.StatusNoContent)
}
func (h *BookHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Id",
		})
		return
	}

	var request dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Book Data",
		})
	}

	patchBook := domain.Book{
		ID:       id,
		Title:    *request.Title,
		Year:     *request.Year,
		AuthorID: *request.AuthorID,
	}

	book, err := h.service.Update(ctx, patchBook)
	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated",
		"data":    book,
	})
}
