package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name string
	Desc string

	Round      uint8  `gorm:"not null;default:0;"`
	DailyPlan  uint16 `gorm:"not null;default:0;"`
	TodayDone  uint16 `gorm:"not null;default:0;"`
	LastDoneAt sql.NullTime
}
