package db

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name string
	Desc string

	Round uint8
}
