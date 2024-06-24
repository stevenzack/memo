package db

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	BookID uint
	Book   Book
	Text   string
}
