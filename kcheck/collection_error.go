package kcheck

import (
	"fmt"
	"strings"
)

type CollectionError struct {
	errors []error
}

func (e *CollectionError) AppendError(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err)
}
func (e *CollectionError) GetErr() error {
	if len(e.errors) == 0 {
		return nil
	}
	return e
}
func (e *CollectionError) Error() string {
	var errorMessages []string
	for _, e := range e.errors {
		errorMessages = append(errorMessages, fmt.Sprint("- ", e.Error()))
	}
	return strings.Join(errorMessages, " \n")
}
