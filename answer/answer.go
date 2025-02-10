package answer

import (
	"errors"
	"log/slog"
	"math"
	"net/http"
	"reflect"
	"strings"

	"github.com/user0608/goones/errs"
)

type Target interface {
	JSON(code int, i interface{}) error
}
type Response struct {
	Type    string      `json:"type,omitempty"` //error-response, success-response
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

const success_response = "success"
const error_message = "error-message"

const SUCCESS = "Operación completada exitosamente"
const CREATED = "Registro guardado con éxito"
const DELETED = "Registro eliminado correctamente"
const UPDATED = "Registro actualizado con éxito"

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

func UnwrapErr(err error) (code int, message string) {
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
			slog.Error("internal error", "error", err)
			return
		}
		if we.Wrapped() != nil {
			slog.Error("internal error", "error", we.Wrapped())
			return
		}
	}(err, werr)
	return code, message
}

func Err(c Target, err error) error {
	code, message := UnwrapErr(err)
	return c.JSON(code, &Response{Type: error_message, Message: message})
}

func JsonErr(c Target) error {
	return Err(c, errs.BadRequestDirect(errs.ErrInvalidRequestBody))
}

func QueryErr(c Target) error {
	return Err(c, errs.BadRequestDirect(errs.ErrInvalidQueryParam))
}

func Auto(c Target, err error) error {
	if err != nil {
		return Err(c, err)
	}
	return Success(c)
}

func AutoOK(c Target, data, err error) error {
	if err != nil {
		return Err(c, err)
	}
	return Ok(c, data)
}

type PageResponse struct {
	Response
	// Page: current page
	Page int64 `json:"page"`
	// PerPage: number of items per page
	PerPage int64 `json:"perPage"`
	// TotalPages: total pages
	TotalPages int64 `json:"totalPages"`
	// TotalItems: total items on the data source
	TotalItems int64 `json:"totalItems"`
	// Items: number of items on the current page
	Items int64 `json:"items"`
}

// page: current page
// perPage: number of items per page
// totalItems: total items on the data source
func OKPage(c Target, page int64, perPage int64, totalItems int64, data any) error {
	return c.JSON(http.StatusOK, &PageResponse{
		Response:   Response{Type: success_response, Data: data},
		Page:       page,
		PerPage:    perPage,
		TotalItems: totalItems,
		Items:      TotalItems(data),
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(perPage))),
	})
}

type LimitOffsetResponse struct {
	Response
	// Limit: number of items per page
	Limit int64 `json:"limit"`
	// Offset: number of items to skip
	Offset int64 `json:"offset"`
	// TotalItems: total items on the data source
	TotalItems int64 `json:"totalItems"`
	// Items: number of items on the current page
	Items int64 `json:"items"`
}

func OKLimitOffset(c Target, limit int64, offset int64, totalItems int64, data any) error {
	return c.JSON(http.StatusOK, &LimitOffsetResponse{
		Response:   Response{Type: success_response, Data: data},
		Limit:      limit,
		Offset:     offset,
		TotalItems: totalItems,
		Items:      TotalItems(data),
	})
}

// if the data is an array, return the number of elements
// otherwise, return 1
func TotalItems(data any) int64 {
	typeOf := reflect.TypeOf(data)
	var kind reflect.Kind
	kind = typeOf.Kind()
	if kind == reflect.Pointer {
		kind = typeOf.Elem().Kind()
	}
	if kind == reflect.Slice {
		return int64(reflect.ValueOf(data).Len())
	}
	return 1
}
