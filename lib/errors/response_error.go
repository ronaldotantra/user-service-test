package errors

import (
	"net/http"
	"strings"
)

type baseError struct {
	code    int
	message string
}

func newBaseError(code int, msg string) baseError {
	return baseError{
		code:    code,
		message: msg,
	}
}

func (err baseError) Error() string {
	return strings.ToLower(err.message)
}

func NewNotFoundError(message string) baseError {
	return newBaseError(http.StatusNotFound, message)
}

func NewForbiddenError(message string) baseError {
	return newBaseError(http.StatusForbidden, message)
}

func NewBadRequestError(message string) baseError {
	return newBaseError(http.StatusBadRequest, message)
}

func NewConflictError(message string) baseError {
	return newBaseError(http.StatusConflict, message)
}
