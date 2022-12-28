package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	_ = v.RegisterValidation("phone", validatePhone)
	return &Validator{validator: v}
}

func (v *Validator) Validate(s interface{}) (bool, []string) {
	isValid := true
	var result []string

	errs := v.validator.Struct(s)
	if errs != nil {
		isValid = false
		for _, v := range errs.(validator.ValidationErrors) {
			result = append(result, message(v))
		}
	}

	return isValid, result
}

func message(fieldError validator.FieldError) string {
	description := ""
	tag := fieldError.Tag()
	switch tag {
	case "phone":
		description = "must be phone number in the format (+7|7|8)?[0-9]{10}"
	default:
		description = fmt.Sprintf("must be %s", tag)
	}

	return fmt.Sprintf("Field '%s' %s", fieldError.Field(), description)
}
