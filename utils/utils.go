package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func MapToStruct(data map[string]string, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf(" Invalid input. Expecting a non-nil pointer to struct")
	}

	typ := val.Elem().Type()
	fmt.Println("Type", typ)
	for fieldName, fieldValue := range data {
		field := val.Elem().FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf(" Invalid field name: %s", fieldName)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(fieldValue, 10, 64)
			if err != nil {
				return fmt.Errorf(" Failed to parse int field %s: %s", fieldName, err)
			}
			field.SetInt(num)
		case reflect.Float32, reflect.Float64:
			num, err := strconv.ParseFloat(fieldValue, 64)
			if err != nil {
				return fmt.Errorf(" Failed to parse float field %s: %s", fieldName, err)
			}
			field.SetFloat(num)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return fmt.Errorf(" Failed to parse bool field %s: %s", fieldName, err)
			}
			field.SetBool(boolVal)
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				date, err := time.Parse("02/01/2006", fieldValue)
				if err != nil {
					return fmt.Errorf(" Failed to parse date field %s: %s", fieldName, err)
				}
				field.Set(reflect.ValueOf(date))
			} else {
				return fmt.Errorf(" Unsupported struct field type for %s", fieldName)
			}
		default:
			return fmt.Errorf(" Unsupported field type for %s", fieldName)
		}
	}
	return nil
}
