package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"virtual-library/catalog-service/internal/catalog"
	"virtual-library/catalog-service/internal/memory"

	"github.com/gin-gonic/gin"
)

var response struct {
	Data    catalog.Book `json:"data"`
	Message string       `json:"message"`
}

func TestBookHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := memory.NewMemoryBookRepository()

	expected, err := repo.Create(catalog.Book{
		Author: "John Doe",
		Title:  "The Big Buck",
		Year:   1942,
	})

	if err != nil {
		t.Fatal(err)
	}

	service := catalog.NewService(repo)
	handler := NewBookHandler(service)

	router := gin.New()
	router.GET("/book/:id", handler.GetByID)
	req, err := http.NewRequest(
		http.MethodGet,
		"/book/"+strconv.Itoa(expected.ID),
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %v, got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Body.Len() == 0 {
		t.Errorf("expected: %v, got: empty", expected.ID)
	}

	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response.Data, expected) {
		t.Errorf("expected: %v, got: %v", expected, response.Data)
	}
}
func TestBookHandler_GetInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := memory.NewMemoryBookRepository()
	service := catalog.NewService(repo)
	handler := NewBookHandler(service)
	router := gin.New()

	router.GET("/book/:id", handler.GetByID)
	req, err := http.NewRequest(
		"GET",
		"/book/abc",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected: %v, got: %v", http.StatusBadRequest, recorder.Code)
	}
}
func TestBookHandler_GetNotExistBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := memory.NewMemoryBookRepository()
	service := catalog.NewService(repo)
	handler := NewBookHandler(service)
	router := gin.New()
	router.GET("/book/:id", handler.GetByID)
	req, err := http.NewRequest("GET", "/book/99", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Message != "Book not found" {
		t.Errorf("expected message: Book not found, got: %v", response.Message)
	}
}
func TestBookHandler_DeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := memory.NewMemoryBookRepository()
	book, err := repo.Create(catalog.Book{
		Author: "John Doe",
		Title:  "The Big Buck",
		Year:   1942,
	})
	if err != nil {
		t.Fatal(err)
	}
	service := catalog.NewService(repo)
	handler := NewBookHandler(service)
	router := gin.New()
	router.DELETE("/book/:id", handler.DeleteBook)
	req, err := http.NewRequest("DELETE", "/book/"+strconv.Itoa(book.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	_, err = repo.GetByID(book.ID)
	if !errors.Is(err, catalog.ErrBookNotFound) {
		t.Errorf("expected: %v, got: %v", catalog.ErrBookNotFound, err)
	}

	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected: %v, got: %v", http.StatusNoContent, recorder.Code)
	}
}
func TestBookHandler_DeleteNotExistBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := memory.NewMemoryBookRepository()
	service := catalog.NewService(repo)
	handler := NewBookHandler(service)
	router := gin.New()
	router.DELETE("/book/:id", handler.DeleteBook)
	req, err := http.NewRequest("DELETE", "/book/99", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	err = json.NewDecoder(recorder.Body).Decode(&response)

	if err != nil {
		t.Fatal(err)
	}

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected: %v, got: %v", http.StatusNoContent, recorder.Code)
	}

	if response.Message != "Book not found" {
		t.Errorf("expected message: Book not found, got: %v", response.Message)
	}
}
