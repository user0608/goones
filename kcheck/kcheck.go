package kcheck

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

const DefaultTagName = "chk"

var ErrInvalidInput = errors.New("kcheck: invalid input")

type ValidatorFunc func(Field) error

type Field struct {
	Name      string
	Path      string
	Tag       string
	Param     string
	Value     any
	Kind      reflect.Kind
	IsNil     bool
	IsPointer bool
}

type Validator struct {
	mu    sync.RWMutex
	tag   string
	funcs map[string]ValidatorFunc
}

type mode int

const (
	modeSkip mode = iota
	modeSelect
)

type options struct {
	mode   mode
	fields map[string]struct{}
}

func New() *Validator {
	v := &Validator{
		tag:   DefaultTagName,
		funcs: make(map[string]ValidatorFunc),
	}

	v.RegisterDefaults()
	return v
}

func (v *Validator) Register(name string, fn ValidatorFunc) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.funcs[name] = fn
}

func (v *Validator) get(name string) (ValidatorFunc, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	fn, ok := v.funcs[name]
	return fn, ok
}

func (v *Validator) Struct(input any) error {
	return v.structWithOptions(input, options{
		mode:   modeSkip,
		fields: map[string]struct{}{},
	})
}

func (v *Validator) StructSkip(input any, skips ...string) error {
	return v.structWithOptions(input, options{
		mode:   modeSkip,
		fields: toSet(skips),
	})
}

func (v *Validator) StructSelect(input any, selected ...string) error {
	return v.structWithOptions(input, options{
		mode:   modeSelect,
		fields: toSet(selected),
	})
}

func (v *Validator) structWithOptions(input any, opts options) error {
	if input == nil {
		return ErrInvalidInput
	}

	rv := reflect.ValueOf(input)

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ErrInvalidInput
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	var errs Errors
	v.validateStruct(rv, "", opts, &errs)

	return errs.Err()
}

func (v *Validator) validateStruct(rv reflect.Value, parentPath string, opts options, errs *Errors) {
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		if sf.PkgPath != "" {
			continue
		}

		fieldName := sf.Name

		path := fieldName
		if parentPath != "" {
			path = parentPath + "." + fieldName
		}

		ignored := shouldIgnore(fieldName, path, opts)

		if ignored && !shouldDiveForSelectedPath(path, opts) {
			continue
		}

		tag := sf.Tag.Get(v.tag)
		if tag == "-" {
			continue
		}

		if shouldDive(fv) {
			v.validateStruct(indirectValue(fv), path, opts, errs)
		}

		if ignored {
			continue
		}

		if tag == "" {
			continue
		}

		field := buildField(path, fieldName, fv)

		for _, rule := range parseRules(tag) {
			field.Tag = rule.Name
			field.Param = rule.Param

			fn, ok := v.get(rule.Name)
			if !ok {
				errs.Add(path, fmt.Sprintf("validador [%s] no registrado", rule.Name))
				continue
			}

			if err := fn(field); err != nil {
				errs.Add(path, err.Error())
			}
		}
	}
}

func shouldIgnore(fieldName string, path string, opts options) bool {
	switch opts.mode {
	case modeSkip:
		return inSet(opts.fields, fieldName) || inSet(opts.fields, path)

	case modeSelect:
		return !inSet(opts.fields, fieldName) && !inSet(opts.fields, path)

	default:
		return false
	}
}

func shouldDiveForSelectedPath(path string, opts options) bool {
	if opts.mode != modeSelect {
		return false
	}

	prefix := path + "."

	for selected := range opts.fields {
		if strings.HasPrefix(selected, prefix) {
			return true
		}
	}

	return false
}

func toSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		set[value] = struct{}{}
	}

	return set
}

func inSet(set map[string]struct{}, value string) bool {
	_, ok := set[value]
	return ok
}

func shouldDive(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return false
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	_, isTime := v.Interface().(time.Time)
	return !isTime
}

func indirectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}

	return v
}

func buildField(path string, name string, v reflect.Value) Field {
	field := Field{
		Name: name,
		Path: path,
	}

	if !v.IsValid() {
		field.IsNil = true
		return field
	}

	field.IsPointer = v.Kind() == reflect.Pointer

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			field.IsNil = true
			field.Kind = v.Kind()
			return field
		}

		v = v.Elem()
	}

	field.Kind = v.Kind()

	if v.CanInterface() {
		field.Value = v.Interface()
	}

	return field
}

type rule struct {
	Name  string
	Param string
}

func parseRules(tag string) []rule {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	parts := strings.Fields(tag)
	rules := make([]rule, 0, len(parts))

	for _, part := range parts {
		name, param, hasParam := strings.Cut(part, "=")
		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		r := rule{Name: name}

		if hasParam {
			r.Param = strings.TrimSpace(param)
		}

		rules = append(rules, r)
	}

	return rules
}

var defaultValidator = New()

func Register(name string, fn ValidatorFunc) {
	defaultValidator.Register(name, fn)
}

func Struct(input any) error {
	return defaultValidator.Struct(input)
}

func Valid(i any, skips ...string) error {
	return defaultValidator.StructSkip(i, skips...)
}

func ValidSelect(i any, selected ...string) error {
	return defaultValidator.StructSelect(i, selected...)
}
