package utils

import (
	"fmt"
	"net/http"
)

type (
	// HTTPError contains error information for errors occurring in a HTTP call
	HttpError struct {
		Message    string
		StatusCode int
	}
)

func (err *HttpError) Error() string {
	return fmt.Sprintf("%s [HTTP %d]", err.Message, err.StatusCode)
}

// BadRequestErrorMessage returns a HTTP 400 error with the given error message
func BadRequestErrorMessage(msg string) *HttpError {
	return &HttpError{
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

// BadRequestError returns a HTTP 400 error with the given error
func BadRequestError(err error) *HttpError {
	return BadRequestErrorMessage(err.Error())
}

// NotFoundError returns a HTTP 404 error with the given error message
func NotFoundErrorMessage(msg string) *HttpError {
	return &HttpError{
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

// NotFoundError returns a HTTP 404 error with the given error
func NotFoundError(err error) *HttpError {
	return NotFoundErrorMessage(err.Error())
}

// UnauthorizedErrorMessage returns a HTTP 401 error with the given error
func UnauthorizedError(err error) *HttpError {
	return UnauthorizedErrorMessage(err.Error())
}

// UnauthorizedErrorMessage returns a HTTP 401 error with the given error message
func UnauthorizedErrorMessage(msg string) *HttpError {
	return &HttpError{
		Message:    msg,
		StatusCode: http.StatusUnauthorized,
	}
}

// InternalServerError returns a HTTP 500 error with the given error message
func InternalServerErrorMessage(msg string) *HttpError {
	return &HttpError{
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
	}
}

// InternalServerError returns a HTTP 500 error with the given error
func InternalServerError(err error) *HttpError {
	return InternalServerErrorMessage(err.Error())
}
