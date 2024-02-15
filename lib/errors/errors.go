package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidatorError(err error) map[string]interface{} {
	e := make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		var errMessage error
		switch v.Tag() {
		case "required":
			errMessage = fmt.Errorf("Field '%s' must be filled", v.Field())
		case "customPassword":
			errMessage = fmt.Errorf("Field '%s'must be minimum 6 characters and maximum 64 characters, containing at least 1 capital characters AND 1 number AND 1 special (nonalpha-numeric) characters.", v.Field())
		case "customPhone":
			errMessage = fmt.Errorf("Field '%s' must start with the Indonesia country code “+62” and must be at minimum 10 characters and maximum 13 characters.", v.Field())
		case "min":
			errMessage = fmt.Errorf("Field '%s' must greater than %s", v.Field(), v.Param())
		case "max":
			errMessage = fmt.Errorf("Field '%s' must less than %s", v.Field(), v.Param())
		default:
			errMessage = fmt.Errorf("Field '%s': '%v' must satisfy '%s' '%v' criteria", v.Field(), v.Value(), v.Tag(), v.Param())
		}
		e[v.Field()] = fmt.Sprintf("%v", errMessage)
	}
	return e
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Something went wrong"
	var data any
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	switch v := err.(type) {
	case *echo.HTTPError:
		code = v.Code
		data = v.Message
	case validator.ValidationErrors:
		code = http.StatusBadRequest
		message = "Validation error"
		data = NewValidatorError(err)
	case baseError:
		code = v.code
		message = v.Error()
	default:
		message = v.Error()
	}
	errResponse := ErrorResponse{
		Code:    fmt.Sprintf("%d", code),
		Message: message,
		Data:    data,
	}

	c.JSON(code, errResponse)
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
