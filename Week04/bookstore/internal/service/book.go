package service

import (
	"log"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/model"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/repository"
)

type BookService struct {
	BookRepository repository.BookRepository
}

func NewBookService(b repository.BookRepository) BookService {
	return BookService{BookRepository: b}
}

func (b *BookService) GetAll() []model.Book {
	return b.BookRepository.GetAll()
}

func (b *BookService) GetByID(id uint) model.Book {
	return b.BookRepository.GetByID(id)
}

func (b *BookService) Save(book model.Book) model.Book {
	log.Println(book)
	return b.BookRepository.Save(book)
}

func (b *BookService) Delete(book model.Book) {
	b.BookRepository.Delete(book)
}
