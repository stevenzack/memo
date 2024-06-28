package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	BookID uint `gorm:"index"`
	Book   Book
	Text   string
	Video  sql.NullString
	Audio  sql.NullString

	FirstReview sql.NullTime  `gorm:"index"`
	Slayed      sql.NullBool  `gorm:"index"`
	Done        sql.NullBool  `gorm:"index"`
	WrongCount  sql.NullInt16 `gorm:"index"`
}
