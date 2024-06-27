package db

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name string
	Desc string

	Round     uint8
	DailyPlan uint16
	TodayDone uint16
}
