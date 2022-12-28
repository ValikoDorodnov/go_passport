package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validatePhone(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^(?:\\+7|7|8)?[0-9]{10}$", fl.Field().String())
	return ok
}
