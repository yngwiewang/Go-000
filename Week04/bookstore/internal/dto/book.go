package dto

import (
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/model"
)

type BookDTO struct {
	ID    uint    `json:"id,string,omitempty"`
	ISBN  string  `json:"isbn"`
	Price float32 `json:"price,string"`
}

func ToBook(bookDTO BookDTO) model.Book {
	return model.Book{
		ISBN:  bookDTO.ISBN,
		Price: bookDTO.Price,
	}
}

func ToBookDTO(book model.Book) BookDTO {
	return BookDTO{
		ID:    book.ID,
		ISBN:  book.ISBN,
		Price: book.Price,
	}
}

func ToBookDTOs(books []model.Book) []BookDTO {
	bookdtos := make([]BookDTO, len(books))
	for i, v := range books {
		bookdtos[i] = ToBookDTO(v)
	}
	return bookdtos
}
