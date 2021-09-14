package book

import (
	"github.com/LibenHailu/sample-books/internal/constant/model"
	"github.com/LibenHailu/sample-books/internal/storage/persistence"
)

// Usecase interface contains function of business logic for domian Book
type Usecase interface {
	InsertBook(book *model.Book) (*model.Book, error)
	GetBooks() ([]model.Book, error)
	GetBook(ID uint) (*model.Book, error)
	UpdateBook(ID uint, book *model.Book) (*model.Book, error)
	DeleteBook(ID uint) error
}

//Service defines all neccessary service for the domain Book
type service struct {
	bookPersist persistence.BookPersistence
}

// Initialize creates a new object with UseCase type
func Initialize(bookPersist persistence.BookPersistence) Usecase {
	return &service{
		bookPersist,
	}
}
