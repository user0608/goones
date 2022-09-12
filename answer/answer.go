package answer

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/user0608/goones/errs"
)

type Target interface {
	JSON(code int, i interface{}) error
}
type Response struct {
	Type       string      `json:"type,omitempty"` //error-response, success-response
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Page       *int        `json:"page,omitempty"`
	TotalPages *int        `json:"tatal_pages,omitempty"`
	TotalItems *int        `json:"total_items,omitempty"`
	NumItems   *int        `json:"num_items"`
}

const success_response = "success"
const success_message = "success-message"
const error_message = "error-message"

func Ok(c Target, payload interface{}) error {
	return c.JSON(http.StatusOK, &Response{
		Type: success_response,
		Data: payload,
	})
}
func payloadLen(payload interface{}) int {
	var vlen = 0
	switch reflect.TypeOf(payload).Kind() {
	case reflect.Slice:
		vlen = reflect.ValueOf(payload).Len()
	}
	return vlen
}
func OkPage(c Target, payload interface{}, page, totalPages, totalItems int) error {
	var numItems = payloadLen(payload)
	return c.JSON(http.StatusOK, &Response{
		Type:       success_response,
		Data:       payload,
		Page:       &page,
		TotalPages: &totalPages,
		TotalItems: &totalItems,
		NumItems:   &numItems,
	})
}
func NewOK(payload interface{}) Response {
	return Response{
		Type: success_response,
		Data: payload,
	}
}

func NewSms(message string) Response {
	return Response{
		Type:    success_message,
		Message: message,
	}
}
func unwrap(err error) (code int, message string) {
	var werr *errs.Err
	code = 400
	message = "algo paso, hubo un error no esperado"
	if errors.As(err, &werr) {
		code = werr.Code()
		message = werr.Message()
	}
	go func(e error, we *errs.Err) {
		if we == nil {
			log.Println(e.Error())
			return
		}
		if we.Wrapped() != nil {
			log.Println(we.Wrapped().Error())
			return
		}
	}(err, werr)
	return code, message
}
func Err(c Target, err error) error {
	code, message := unwrap(err)
	return c.JSON(code, &Response{Type: error_message, Message: message})
}
func Error(err error) Response {
	_, message := unwrap(err)
	return Response{Type: error_message, Message: message}
}

func JsonErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidRequestBody))
}
func QueryErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidQueryParam))
}
