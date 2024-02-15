package validator

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

func NewValidator() *Validator {
	validator := validator.New()
	validator.RegisterValidation("customPassword", validateCustomPassword)
	validator.RegisterValidation("customPhone", validateCustomPhone)
	return &Validator{
		validator: validator,
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func validateCustomPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasUpperCase := regexp.MustCompile("[A-Z]").MatchString(password)
	hasNumber := regexp.MustCompile("[0-9]").MatchString(password)
	hasSpecial := regexp.MustCompile("[^A-Za-z0-9]").MatchString(password)
	length := len(password) >= 6 && len(password) <= 64

	return hasUpperCase && hasNumber && hasSpecial && length
}

func validateCustomPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	regex := regexp.MustCompile(`^\+62[0-9]{10,13}$`)

	return regex.MatchString(phone)
}
