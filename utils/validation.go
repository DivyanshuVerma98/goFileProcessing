package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
)

func IsValidDateFormat(dateString string) bool {
	layout := "02/01/2006"
	_, err := time.Parse(layout, dateString)
	return err == nil
}

func ValidateTransactionType(transaction_type string) (bool, string) {
	var val_list = []string{constants.Primary, constants.Adjustment, constants.Endorsement}
	transaction_type = strings.TrimSpace(transaction_type)
	transaction_type = strings.ToUpper(transaction_type)
	if len(transaction_type) == 0 {
		return false, "transaction_type can't be empty"
	}
	for _, val := range val_list {
		if transaction_type == val {
			return true, "Success"
		}
	}
	return false, "Invalid transaction_type. Choose from " + strings.Join(val_list, ", ")
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

func ValidateMotorBatchData(batch_data *structs.BatchData) {
	for key, val := range batch_data.PolicyDetails.DataMap {

		for _, fval := range constants.MotorMandatoryFields {
			if len(val[fval]) == 0 {
				batch_data.ErrorDetails.MessageMap[key] = fval + " can't be empty."
				continue
			}
		}
		transaction_type := val[constants.TransactionType]
		is_valid, msg := ValidateTransactionType(transaction_type)
		if !is_valid {
			batch_data.ErrorDetails.MessageMap[key] = msg
		}

		insurer_name := val[constants.InsuredName]
		insurer_name = strings.TrimSpace(insurer_name)
		insurer_name = strings.ToLower(insurer_name)

		product := val[constants.Product]
		product = strings.TrimSpace(product)
		product = strings.ToLower(product)

		insurer_mandatory_fields := constants.MotorInsurerJsonData[insurer_name][product]
		for _, fval := range insurer_mandatory_fields {
			if len(val[fval]) == 0 {
				batch_data.ErrorDetails.MessageMap[key] = fval + " can't be empty."
				continue
			}
		}

	}
}
