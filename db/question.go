package db

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	BookID uint `gorm:"index"`
	Book   Book
	Text   string

	Slayed     bool   `gorm:"index"`
	Done       bool   `gorm:"index"`
	WrongCount uint16 `gorm:"index"`
}
