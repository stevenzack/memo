package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type Option struct {
	gorm.Model
	QuestionID uint
	Question   Question
	Text       string
	Video      sql.NullString
	Audio      sql.NullString
	Images     sql.NullString //split by comma
}
