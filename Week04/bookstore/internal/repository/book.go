package repository

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/model"
)

type BookRepository interface {
	GetAll() []model.Book
	GetByID(id uint) model.Book
	Save(book model.Book) model.Book
	Delete(book model.Book)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (b *bookRepository) GetAll() []model.Book {
	var books []model.Book
	b.db.Find(&books)
	return books
}

func (b *bookRepository) GetByID(id uint) model.Book {
	var book model.Book
	b.db.First(&book, id)
	return book
}

func (b *bookRepository) Save(book model.Book) model.Book {
	log.Println(book)
	b.db.Save(&book)
	return book
}

func (b *bookRepository) Delete(book model.Book) {
	fmt.Println(book)
	b.db.Delete(&book)
}
