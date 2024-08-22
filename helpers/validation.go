package helpers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateStruct(s interface{}) error {
	validate = validator.New()

	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := getJSONTag(s, err.Field())
			tag := err.Tag()
			param := err.Param()

			errorMessages := fmt.Sprintf("Field '%s' is invalid: %s %s", field, tag, param)

			if tag == "required" {
				errorMessages = fmt.Sprintf("Field '%s': is %s", field, tag)
			}

			if tag == "email" {
				errorMessages = fmt.Sprintf("Field '%s': invalid %s", field, tag)
			}

			if tag == "min" {
				errorMessages = fmt.Sprintf("Field '%s': minimal %s characters", field, param)
			}

			return errors.New(errorMessages)
		}
	}

	return nil
}

func getJSONTag(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("json")
}
