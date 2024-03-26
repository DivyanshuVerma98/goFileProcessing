package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
)

func IsValidDateFormat(dateString string) error {
	layout := "02/01/2006"
	_, err := time.Parse(layout, dateString)
	return err
}

func ValidateTransactionType(transaction_type string) error {
	var val_list = []string{constants.Primary, constants.Adjustment, constants.Endorsement}
	transaction_type = strings.TrimSpace(transaction_type)
	transaction_type = strings.ToUpper(transaction_type)
	if len(transaction_type) == 0 {
		return fmt.Errorf("transaction_type can't be empty")
	}
	for _, val := range val_list {
		if transaction_type == val {
			return nil
		}
	}
	return fmt.Errorf(" Invalid transaction_type. Choose from " + strings.Join(val_list, ", "))
}

func ValidateHeaders(headers []string, csv_to_model_map map[string]string) error {
	field_list := []string{}
	for val := range csv_to_model_map {
		field_list = append(field_list, val)
	}
	for _, val := range headers {
		_, exists := csv_to_model_map[val]
		if !exists {
			return fmt.Errorf(" Unsupported head: %s", val)
		}
		for i, fval := range field_list {
			if fval == val {
				field_list = append(field_list[:i], field_list[i+1:]...)
				break
			}
		}
	}
	if len(field_list) != 0 {
		return fmt.Errorf(" Missing headers: %s", strings.Join(field_list, ", "))
	}
	return nil
}
