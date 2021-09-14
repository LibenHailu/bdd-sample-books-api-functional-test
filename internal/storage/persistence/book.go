package persistence

import (
	"log"

	"github.com/LibenHailu/sample-books/internal/constant/model"
	"gorm.io/gorm"
)

// BookPersistence contains a list of functions for database operation for table book
type BookPersistence interface {
	InsertBook(book *model.Book) (*model.Book, error)
	GetBooks() ([]model.Book, error)
	GetBook(ID uint) (*model.Book, error)
	UpdateBook(ID uint, book *model.Book) (*model.Book, error)
	DeleteBook(ID uint) error
}

// bookPersistence defines all the neccsary things for the database opertation on domina book
type bookPersistence struct {
	conn *gorm.DB
}

// BookInit creates a new BookPersistence object
func BookInit(conn *gorm.DB) BookPersistence {
	return &bookPersistence{
		conn,
	}
}

// InsertBook persists book data to the database
func (bp *bookPersistence) InsertBook(book *model.Book) (*model.Book, error) {

	err := bp.conn.Create(book).Error
	if err != nil {
		log.Printf("Error when saving book to db: %v", err)
		return nil, err
	}
	return book, nil
}

// GetBooks returns all Books of our system
func (bp *bookPersistence) GetBooks() ([]model.Book, error) {
	books := []model.Book{}

	err := bp.conn.Find(&books).Error

	if err != nil {
		log.Printf("Error when quering books: %v", err)
		return nil, err
	}
	return books, nil
}

// GetBook returns a Books with an id given
func (bp *bookPersistence) GetBook(ID uint) (*model.Book, error) {
	book := model.Book{
		ID: ID,
	}

	err := bp.conn.First(&book).Error

	if err != nil {
		log.Printf("Error when finding a user to update")
		return nil, err
	}

	return &book, nil
}

// UpdateBook updates book
func (bp *bookPersistence) UpdateBook(ID uint, book *model.Book) (*model.Book, error) {

	updated := model.Book{
		ID: ID,
	}

	err := bp.conn.First(&updated).Error

	if err != nil {
		log.Printf("Error when finding a user to update: %v", err)
		return nil, err
	}

	err = bp.conn.Model(&updated).Updates(book).Error

	if err != nil {
		log.Printf("Error when updating a book: %v", err)
		return nil, err
	}

	return &updated, nil

}

// DeleteBook delete book
func (bp *bookPersistence) DeleteBook(ID uint) error {
	book := model.Book{
		ID: ID,
	}

	err := bp.conn.First(&book).Error

	if err != nil {
		log.Printf("Error when finding a user to update")
		return err
	}

	err = bp.conn.Delete(&model.Book{}, ID).Error

	if err != nil {
		log.Printf("Error when deleting a book: %v", err)
		return err
	}
	return nil
}
