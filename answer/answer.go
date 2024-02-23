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

const SUCCESS = "Operación completada exitosamente."
const CREATED = "Registro guardado con éxito."
const DELETED = "Registro eliminado correctamente."
const UPDATED = "Registro actualizado con éxito."

func Ok(c Target, payload interface{}) error {
	return c.JSON(http.StatusOK, &Response{
		Type: success_response,
		Data: payload,
	})
}

func Message(c Target, message string) error {
	return c.JSON(http.StatusOK, &Response{Message: message})
}

func Success(c Target) error { return c.JSON(http.StatusOK, &Response{Message: SUCCESS}) }

func Created(c Target) error { return c.JSON(http.StatusCreated, &Response{Message: CREATED}) }

func Updated(c Target) error { return c.JSON(http.StatusOK, &Response{Message: UPDATED}) }

func Deleted(c Target) error { return c.JSON(http.StatusOK, &Response{Message: DELETED}) }

func unwrap(err error) (code int, message string) {
	var werr *errs.Err
	code = http.StatusInternalServerError
	message = "Ocurrió un problema. Se produjo un error inesperado."
	if errors.As(err, &werr) {
		code = werr.Code()
		message = werr.Message()
	}
	if werr == nil && err != nil {
		var errSMS = strings.TrimSpace(err.Error())
		if strings.HasPrefix(":", errSMS) {
			code = http.StatusBadRequest
			message = strings.TrimLeft(errSMS, ":")
			return
		}
	}
	go func(err error, we *errs.Err) {
		if we == nil && err != nil {
			log.Println("ERR:", err.Error())
			return
		}
		if we.Wrapped() != nil {
			log.Println("ERR:", strings.TrimSpace(we.Wrapped().Error()))
			return
		}
	}(err, werr)
	return code, message
}

func Err(c Target, err error) error {
	code, message := unwrap(err)
	return c.JSON(code, &Response{Type: error_message, Message: message})
}

func JsonErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidRequestBody))
}

func QueryErr(c Target) error {
	return Err(c, errs.Bad(errs.ErrInvalidQueryParam))
}

func Auto(c Target, err error) error {
	if err != nil {
		return Err(c, err)
	}
	return Success(c)
}
