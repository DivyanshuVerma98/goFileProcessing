package utils

import (
	"strings"
	"time"
)

func IsValidDateFormat(dateString string) bool {
	layout := "02/01/2006"
	_, err := time.Parse(layout, dateString)
	return err == nil
}

func ValidateHeaders(headers []string, csv_to_model_map map[string]interface{}) (bool, string) {
	field_list := []string{}
	for val := range csv_to_model_map {
		field_list = append(field_list, val)
	}
	for _, val := range headers {
		if csv_to_model_map[val] == nil {
			return false, "Unsupported head: " + val
		}
		for i, fval := range field_list {
			if fval == val {
				field_list = append(field_list[:i], field_list[i+1:]...)
				break
			}
		}
	}
	if len(field_list) != 0 {
		return false, "Missing headers: " + strings.Join(field_list, ", ")
	}
	return true, "Success"
}
