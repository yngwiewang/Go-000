package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/dto"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/model"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/service"
)

type BookAPI struct {
	BookService service.BookService
}

func NewBookAPI(b service.BookService) BookAPI {
	return BookAPI{BookService: b}
}

func (b *BookAPI) GetAll(c *gin.Context) {
	books := b.BookService.GetAll()

	c.JSON(http.StatusOK, gin.H{"books": dto.ToBookDTOs(books)})
}

func (b *BookAPI) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book := b.BookService.GetByID(uint(id))

	c.JSON(http.StatusOK, gin.H{"book": dto.ToBookDTO(book)})
}

func (b *BookAPI) Create(c *gin.Context) {
	var bookDTO dto.BookDTO
	err := c.BindJSON(&bookDTO)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	createBook := b.BookService.Save(dto.ToBook(bookDTO))

	c.JSON(http.StatusOK, gin.H{"book": dto.ToBookDTO(createBook)})
}

func (b *BookAPI) Update(c *gin.Context) {
	var bookDTO dto.BookDTO
	err := c.BindJSON(&bookDTO)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	book := b.BookService.GetByID(uint(id))
	log.Println(book)
	if book == (model.Book{}) {
		c.Status(http.StatusBadRequest)
		return
	}

	book.ISBN = bookDTO.ISBN
	book.Price = bookDTO.Price
	log.Println(book)
	b.BookService.Save(book)

	c.Status(http.StatusOK)
}

func (b *BookAPI) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book := b.BookService.GetByID(uint(id))
	if book == (model.Book{}) {
		c.Status(http.StatusBadRequest)
		return
	}
	fmt.Println(book)
	b.BookService.Delete(book)

	c.Status(http.StatusOK)
}
