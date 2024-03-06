package utils

import "time"

func IsValidDateFormat(dateString string) bool {
	layout := "02/01/2006"
	_, err := time.Parse(layout, dateString)
	return err == nil
}
