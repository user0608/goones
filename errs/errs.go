package errs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Err struct {
	httpCode int
	wrapped  error
	message  string
}

func (err *Err) Error() string {
	if err.wrapped == nil {
		return fmt.Sprintf("error: %s", err.message)
	}
	return fmt.Sprintf("error: %s; wrapped: %s", err.message, err.wrapped.Error())
}

func (err *Err) Message() string {
	return err.message
}

func (err *Err) Code() int {
	return err.httpCode
}

func (err *Err) Wrapped() error {
	return err.wrapped
}

func newError(err error, message string, httpCode int) error {
	return &Err{
		wrapped:  err,
		message:  message,
		httpCode: httpCode,
	}
}

// ContainsMessage checks if the message exists in the wrapped error message
func ContainsMessage(err error, message string) bool {
	var customErr *Err
	if errors.As(err, &customErr) {
		return strings.Contains(customErr.message, message)
	}
	return false
}

// NewWithMessage creates a new error with a custom message, defaulting to BadRequest if the error is not of type Err
func NewWithMessage(err error, message string) error {
	var customErr *Err
	if errors.As(err, &customErr) {
		return &Err{
			httpCode: customErr.httpCode,
			wrapped:  customErr.wrapped,
			message:  message,
		}
	}
	return newError(err, message, http.StatusBadRequest)
}

func BadRequestError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusBadRequest)
}

func NotFoundError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusNotFound)
}

func InternalError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusInternalServerError)
}

func UnsupportedMediaTypeError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusUnsupportedMediaType)
}

func UnauthorizedError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusUnauthorized)
}

func ForbiddenError(err error, format string, args ...interface{}) error {
	return newError(err, fmt.Sprintf(format, args...), http.StatusForbidden)
}

func BadRequestf(format string, args ...interface{}) error {
	return newError(nil, fmt.Sprintf(format, args...), http.StatusBadRequest)
}

func NotFoundf(format string, args ...interface{}) error {
	return newError(nil, fmt.Sprintf(format, args...), http.StatusNotFound)
}

func InternalErrorf(format string, args ...interface{}) error {
	return newError(nil, fmt.Sprintf(format, args...), http.StatusInternalServerError)
}

// Version with direct messages (no formatting)
func BadRequestDirect(message string) error {
	return newError(nil, message, http.StatusBadRequest)
}

func NotFoundDirect(message string) error {
	return newError(nil, message, http.StatusNotFound)
}

func InternalErrorDirect(message string) error {
	return newError(nil, message, http.StatusInternalServerError)
}

func UnauthorizedDirect(message string) error {
	return newError(nil, message, http.StatusUnauthorized)
}

func ForbiddenDirect(message string) error {
	return newError(nil, message, http.StatusForbidden)
}

func UnsupportedMediaTypeDirect(message string) error {
	return newError(nil, message, http.StatusUnsupportedMediaType)
}

func WrapError(err error, message string, httpCode int) error {
	if err == nil {
		return nil
	}
	return newError(err, message, httpCode)
}

func IsBadRequest(err error) bool {
	var customErr *Err
	if errors.As(err, &customErr) {
		return customErr.Code() == http.StatusBadRequest
	}
	return false
}

func IsInternalError(err error) bool {
	var customErr *Err
	if errors.As(err, &customErr) {
		return customErr.Code() == http.StatusInternalServerError
	}
	return false
}

func IsErr(err error) bool {
	var customErr *Err
	return errors.As(err, &customErr)
}

func ToSummary(err error) string {
	if err == nil {
		return "No error"
	}
	var customErr *Err
	if errors.As(err, &customErr) {
		return fmt.Sprintf("Error %d: %s", customErr.Code(), customErr.Message())
	}
	return err.Error()
}
