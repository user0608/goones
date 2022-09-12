package errs

import (
	"fmt"
	"net/http"
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
