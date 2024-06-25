package db

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	BookID uint `gorm:"index"`
	Book   Book
	Text   string
}
