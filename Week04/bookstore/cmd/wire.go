package main

import (
	"github.com/jinzhu/gorm"

	"github.com/google/wire"
	v1 "github.com/yngwiewang/Go-000/Week04/bookstore/api/v1"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/repository"
	"github.com/yngwiewang/Go-000/Week04/bookstore/internal/service"
)

func initBookAPI(db *gorm.DB) v1.BookAPI {
	wire.Build(repository.NewBookRepository, service.NewBookService, v1.NewBookAPI)
	return v1.BookAPI{}
}
