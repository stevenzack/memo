package util

import "time"

func YesterdayAgo(t time.Time) []any {
	return []any{time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location()), time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())}
}

func ThreeDaysAgo(t time.Time) []any {
	return []any{
		time.Date(t.Year(), t.Month(), t.Day()-3, 0, 0, 0, 0, t.Location()),
		time.Date(t.Year(), t.Month(), t.Day()-2, 0, 0, 0, 0, t.Location()),
	}
}

func SevenDaysAgo(t time.Time) []any {
	return []any{
		time.Date(t.Year(), t.Month(), t.Day()-7, 0, 0, 0, 0, t.Location()),
		time.Date(t.Year(), t.Month(), t.Day()-6, 0, 0, 0, 0, t.Location()),
	}
}

func OneMonthAgo(t time.Time) []any {
	return []any{
		time.Date(t.Year(), t.Month(), t.Day()-30, 0, 0, 0, 0, t.Location()),
		time.Date(t.Year(), t.Month(), t.Day()-29, 0, 0, 0, 0, t.Location()),
	}
}
