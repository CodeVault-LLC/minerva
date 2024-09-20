package utils

import "time"

func CurrentDate() time.Time {
	return time.Now()
}

func CompareDates(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day()
}

func Get24HoursAgo() time.Time {
	return time.Now().AddDate(0, 0, -1)
}
