package db

import "gorm.io/gorm"

type Option struct {
	gorm.Model
	QuestionID uint
	Question   Question
	Text       string
	Video      string
	Audio      string
}
