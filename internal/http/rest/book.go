package rest

import (
	"encoding/json"
	errCheck "errors"
	"net/http"
	"strconv"

	"github.com/LibenHailu/sample-books/internal/constant/model"
	"github.com/LibenHailu/sample-books/internal/module/book"
	"github.com/LibenHailu/sample-books/pkg/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookHandler contains a function of handlers for the domain Book
type BookHandler interface {
	GetBook(c *gin.Context)
	GetBooks(c *gin.Context)
	InsertBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

// bookHandler defines all the things neccessary for book handling
type bookHandler struct {
	bookUsecase book.Usecase
}

// BookInit initializes a book handler for the domain book
func BookInit(useCase book.Usecase) BookHandler {
	return &bookHandler{
		bookUsecase: useCase,
	}
}

// GetBook gets a book
// GET /v1/books/:id
func (bh bookHandler) GetBook(c *gin.Context) {

	ID := c.Param("id")
	uintID, err := strconv.ParseUint(ID, 10, 32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Invalid id param"))
		return
	}
	book, err := bh.bookUsecase.GetBook(uint(uintID))
	if err != nil {
		if errCheck.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, errors.ErrorMap("Book not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Unexpected server error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
	return

}

// GetBooks gets books
// GET /v1/books
func (bh bookHandler) GetBooks(c *gin.Context) {
	books, err := bh.bookUsecase.GetBooks()

	if err != nil {

		// c.JSON(http.StatusInternalServerError, errors.ErrorMap("Unexpected server error"))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
	return
}

// InsertBook creates a book
// POST /v1/books
func (bh bookHandler) InsertBook(c *gin.Context) {

	var insertBook model.Book
	if err := c.ShouldBind(&insertBook); err != nil {

		// checks for validation
		valError := errors.ValidationError(err)
		if len(valError) != 0 {
			c.JSON(http.StatusBadRequest, valError)
			return
		}

		c.JSON(http.StatusBadRequest, errors.ErrorMap("Invalid request"))
		return
	}
	book, err := bh.bookUsecase.InsertBook(&insertBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Unexpected server error"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"book": book})
	return
}

// UpdateBook updates a book
// PATCH /v1/books/:id
func (bh bookHandler) UpdateBook(c *gin.Context) {
	ID := c.Param("id")
	uintID, err := strconv.ParseUint(ID, 10, 32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Invalid id param"))
		return
	}

	var book model.Book
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorMap("Bad request"))
		return

	}

	err = json.Unmarshal(jsonData, &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorMap("Bad request"))
		return
	}

	updatedBook, err := bh.bookUsecase.UpdateBook(uint(uintID), &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Unexpected server error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": updatedBook})
	return
}

// DeleteBook deletes a book
// DELETE /v1/books/:id
func (bh bookHandler) DeleteBook(c *gin.Context) {
	ID := c.Param("id")
	uintID, err := strconv.ParseUint(ID, 10, 32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Invalid id param"))
		return
	}

	err = bh.bookUsecase.DeleteBook(uint(uintID))

	if err != nil {
		if errCheck.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, errors.ErrorMap("Book not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, errors.ErrorMap("Unexpected server error"))
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
	return
}
