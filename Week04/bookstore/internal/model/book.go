package model

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	ISBN  string
	Price float32
}
