package model

import "time"

// Book is the structure for the table files
type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	Isbn      string    `json:"isbn" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
