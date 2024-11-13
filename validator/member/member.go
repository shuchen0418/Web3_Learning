package member

import (
	"github.com/go-playground/validator/v10"
)

func NameValid(fl validator.FieldLevel) bool {
	if s, ok := fl.Field().Interface().(string); ok {
		if s == "admin" {
			return false
		}
	}
	return true
}
