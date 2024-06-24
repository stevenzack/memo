package db

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	QuestionID uint
	Question   Question
	IsCorrect  bool
	Text       string
}
