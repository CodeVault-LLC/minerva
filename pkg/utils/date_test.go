package utils

import "testing"

func TestCurrentDate(t *testing.T) {
	date := CurrentDate()
	if date.IsZero() {
		t.Error("Expected current date, got zero date")
	}

	if date.Year() < 2023 {
		t.Error("Expected year to be greater than 2023")
	}
}

func TestCompareDates(t *testing.T) {
	date1 := CurrentDate()
	date2 := CurrentDate()

	if !CompareDates(date1, date2) {
		t.Error("Expected dates to be equal")
	}

	date2 = date2.AddDate(0, 0, 1)

	if CompareDates(date1, date2) {
		t.Error("Expected dates to be different")
	}
}

func TestGet24HoursAgo(t *testing.T) {
	date := Get24HoursAgo()
	if date.IsZero() {
		t.Error("Expected date, got zero date")
	}

	if date.Year() < 2023 {
		t.Error("Expected year to be greater than 2023")
	}
}
