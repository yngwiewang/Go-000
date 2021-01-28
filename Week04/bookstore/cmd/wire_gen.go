// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/jinzhu/gorm"
	"github.com/yngwiewang/Go-000/Week04/bookstore/api/v1"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/repository"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/service"
)

// Injectors from wire.go:

func InitBookAPI(db *gorm.DB) v1.BookAPI {
	bookRepository := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookAPI := v1.NewBookAPI(bookService)
	return bookAPI
}