package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidationErrors(req interface{}) map[string]string {
	errors := ValidateStruct(req)
	if len(errors) > 0 {
		errorMap := make(map[string]string)
		for _, err := range errors {
			errorMap[err.Field] = err.Message
		}
		return errorMap
	}
	return nil
}

func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = getJSONFieldName(s, err.Field())
			element.Message = getErrorMessage(err, element.Field)
			errors = append(errors, element)
		}
	}
	return errors
}

func getJSONFieldName(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fieldName); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
				return jsonTag[:commaIndex]
			}
			return jsonTag
		}
	}

	return toSnakeCase(fieldName)
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}
	return string(result)
}

func getErrorMessage(fe validator.FieldError, jsonFieldName string) string {
	fieldDisplayName := strings.ReplaceAll(jsonFieldName, "_", " ")

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("please enter your %s", fieldDisplayName)
	case "email":
		return "please enter a valid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fieldDisplayName, fe.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", fieldDisplayName, fe.Param())
	case "oneof":
		return fmt.Sprintf("please select a valid option for %s", fieldDisplayName)
	case "gte":
		return fmt.Sprintf("%s must be at least %s", fieldDisplayName, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fieldDisplayName, fe.Param())
	case "lte":
		return fmt.Sprintf("%s cannot be more than %s", fieldDisplayName, fe.Param())
	case "alpha":
		return fmt.Sprintf("%s can only contain letters", fieldDisplayName)
	case "alphanum":
		return fmt.Sprintf("%s can only contain letters and numbers", fieldDisplayName)
	case "numeric":
		return fmt.Sprintf("%s must be a number", fieldDisplayName)
	case "url":
		return "Please enter a valid website URL"
	case "uuid":
		return "Invalid ID format"
	case "eqfield":
		return fmt.Sprintf("%s must match with the previous field", fieldDisplayName)
	case "datetime":
		return fmt.Sprintf("invalid %s format, please enter valid date and time", fieldDisplayName)
	case "date":
		return fmt.Sprintf("invalid %s format, please enter valid date", fieldDisplayName)
	case "time":
		return fmt.Sprintf("invalid %s format, please enter valid time", fieldDisplayName)
	case "unique":
		return fmt.Sprintf("%s cannot contain duplicate values", fieldDisplayName)
	default:
		return fmt.Sprintf("please check your %s and try again", fieldDisplayName)
	}
}
