package answer

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/user0608/goones/errs"
)

type Target interface {
	JSON(code int, i interface{}) error
}
type Response struct {
	Type           string      `json:"type,omitempty"` //error-response, success-response
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	ReturnedItems  *int        `json:"returned_items,omitempty"`
	RequestedItems *int        `json:"requested_items,omitempty"`
	CurrentPage    *int        `json:"current_page,omitempty"`
	NumberPages    *int        `json:"number_pages,omitempty"`
	NumberItems    *int        `json:"items,omitempty"`
}

const success_response = "success"
const error_message = "error-message"

const SUCCESS = "operacion realizada"
const CREATED = "registro guardado"
const DELETED = "registro eliminado"
const UPDATED = "registro actualizado"

func Ok(c Target, payload interface{}) error {
	return c.JSON(http.StatusOK, &Response{
		Type: success_response,
		Data: payload,
	})
}
func Message(c Target, message string) error {
	return c.JSON(http.StatusOK, &Response{Message: message})
}

func OkPage(c Target, p Pager) error {
	var currentpage = p.CurrentPage()
	var numberitems = p.NumberItems()
	var numberpages = p.NumberPages()
	var requestedItems = p.RequestedItems()
	var returntedItems = p.ReturnedItems()
	return c.JSON(http.StatusOK, &Response{
		Type:           success_response,
		Data:           p.Data(),
		NumberPages:    &numberpages,
		ReturnedItems:  &returntedItems,
		RequestedItems: &requestedItems,
		CurrentPage:    &currentpage,
		NumberItems:    &numberitems,
	})
}
func NewOK(payload any) Response {
	return Response{
		Type: success_response,
		Data: payload,
	}
}

func NewSms(message string) Response {
	return Response{
		Type:    success_response,
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
			log.Println(strings.TrimSpace(we.Wrapped().Error()))
			return
		}
	}(err, werr)
	return code, message
}
func Err(c Target, err error) error {
	code, message := unwrap(err)
	return c.JSON(code, &Response{Type: error_message, Message: message})
}
func Error(err error, payload ...any) Response {
	_, message := unwrap(err)
	var data any
	if len(payload) > 0 {
		data = payload[0]
	}
	return Response{Type: error_message, Message: message, Data: data}
}

func JsonErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidRequestBody))
}
func QueryErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidQueryParam))
}
