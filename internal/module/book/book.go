package book

import "github.com/LibenHailu/sample-books/internal/constant/model"

// InsertBook implements business and calls a database operation to persist a new book
func (s *service) InsertBook(book *model.Book) (*model.Book, error) {
	book, err := s.bookPersist.InsertBook(book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBooks implements business and calls a database operation to get books
func (s *service) GetBooks() ([]model.Book, error) {
	books, err := s.bookPersist.GetBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

// GetBook implements business and calls a database operation to get a book
func (s *service) GetBook(ID uint) (*model.Book, error) {
	book, err := s.bookPersist.GetBook(ID)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// UpdateBook implements business and calls a database operation to update a book
func (s *service) UpdateBook(ID uint, book *model.Book) (*model.Book, error) {
	book, err := s.bookPersist.UpdateBook(ID, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// DeleteBook implements buisines rules and calls a database operation to delete a book
func (s *service) DeleteBook(ID uint) error {
	return s.bookPersist.DeleteBook(ID)
}
