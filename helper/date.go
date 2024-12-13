package helper

import "time"

func DateTime(dateString string) time.Time {
	layoutFormat := "2006-01-02 15:04:05"
	date, _ := time.Parse(layoutFormat, dateString)
	return date
}

func Date(dateString string) time.Time {
	layoutFormat := "2006-01-02"
	date, _ := time.Parse(layoutFormat, dateString)
	return date
}
