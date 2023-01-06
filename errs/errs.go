package errs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Err struct {
	httpcode int
	wrapped  error
	message  string
}

func (err *Err) Error() string {
	if err.wrapped == nil {
		return "the errro is trivial "
	}
	return err.wrapped.Error()
}

func (err *Err) Message() string {
	return err.message
}
func (err *Err) Code() int {
	return err.httpcode
}
func (err *Err) Wrapped() error {
	return err.wrapped
}

func newErrf(err error, message string, httpcode int) error {
	return &Err{
		wrapped:  err,
		message:  message,
		httpcode: httpcode,
	}
}

// Contains revisa el mensaje
func ContainsMessage(err error, message string) bool {
	var myerr *Err
	if errors.As(err, &myerr) {
		return strings.Contains(myerr.message, message)
	}
	return false
}

// NewWithMessage if err isn't errs.Err, bad request is by default
func NewWithMessage(err error, message string) error {
	var myerr *Err
	if errors.As(err, &myerr) {
		return &Err{
			httpcode: myerr.httpcode,
			wrapped:  myerr.wrapped,
			message:  message,
		}
	}
	return BadReqf(err, message)
}
func BadReqf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusBadRequest)
}
func Notfoundf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusForbidden)
}
func Internalf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusInternalServerError)
}
func Unsupportedf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusUnsupportedMediaType)
}
func Unauthorizedf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusUnauthorized)
}
func Forbibbenf(err error, format string, a ...interface{}) error {
	return newErrf(err, fmt.Sprintf(format, a...), http.StatusForbidden)
}
func Bad(format string, a ...interface{}) error {
	return newErrf(nil, fmt.Sprintf(format, a...), http.StatusBadRequest)
}
func NF(format string, a ...interface{}) error {
	return newErrf(nil, fmt.Sprintf(format, a...), http.StatusNotFound)
}
func Internal(format string, a ...interface{}) error {
	return newErrf(nil, fmt.Sprintf(format, a...), http.StatusInternalServerError)
}
