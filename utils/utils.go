package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MapToStruct(data map[string]string, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf(" Invalid input. Expecting a non-nil pointer to struct")
	}

	for fieldName, fieldValue := range data {
		field := val.Elem().FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf(" Invalid field name: %s", fieldName)
		}
		fieldValue = strings.TrimSpace(fieldValue)
		if strings.ToUpper(fieldValue) == "NA" {
			fieldValue = ""
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if len(fieldValue) > 0 {
				num, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return fmt.Errorf(" Failed to parse int field %s: %s", fieldName, err)
				}
				field.SetInt(num)
			}
		case reflect.Float32, reflect.Float64:
			if len(fieldValue) > 0 {
				num, err := strconv.ParseFloat(fieldValue, 64)
				if err != nil {
					return fmt.Errorf(" Failed to parse float field %s: %s", fieldName, err)
				}
				field.SetFloat(num)
			}
		case reflect.Bool:
			if len(fieldValue) > 0 {
				boolVal, err := strconv.ParseBool(fieldValue)
				if err != nil {
					return fmt.Errorf(" Failed to parse bool field %s: %s", fieldName, err)
				}
				field.SetBool(boolVal)
			}
		case reflect.Struct:
			if len(fieldValue) > 0 {
				formats := []string{"02/01/2006", "02-01-2006"}
				if field.Type() == reflect.TypeOf(time.Time{}) {
					var date time.Time
					var err error
					for _, format := range formats {
						date, err = time.Parse(format, fieldValue)
						if err == nil {
							break
						}
					}
					field.Set(reflect.ValueOf(date))
				} else if field.Type() == reflect.TypeOf(sql.NullTime{}) {
					// Handle sql.NullTime type
					if strings.TrimSpace(fieldValue) != "" {
						var date time.Time
						var err error
						for _, format := range formats {
							date, err = time.Parse(format, fieldValue)
							if err == nil {
								break
							}
						}
						nullTime := sql.NullTime{Time: date, Valid: true}
						field.Set(reflect.ValueOf(nullTime))
					} else {
						nullTime := sql.NullTime{Valid: false}
						field.Set(reflect.ValueOf(nullTime))
					}
				} else {
					return fmt.Errorf(" Unsupported struct field type for %s", fieldName)
				}
			}
		default:
			return fmt.Errorf(" Unsupported field type for %s", fieldName)
		}
	}
	return nil
}
