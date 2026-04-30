package kcheck

import "strings"

type FieldError struct {
	Field   string
	Message string
}

type Errors struct {
	Items []FieldError
}

func (e *Errors) Add(field string, message string) {
	e.Items = append(e.Items, FieldError{
		Field:   field,
		Message: message,
	})
}

func (e Errors) Error() string {
	if len(e.Items) == 0 {
		return ""
	}

	var b strings.Builder

	for i, item := range e.Items {
		if i > 0 {
			b.WriteString("; ")
		}

		b.WriteString(item.Field)
		b.WriteString(": ")
		b.WriteString(item.Message)
	}

	return b.String()
}

func (e Errors) Err() error {
	if len(e.Items) == 0 {
		return nil
	}

	return e
}
