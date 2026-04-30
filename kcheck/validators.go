package kcheck

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (v *Validator) RegisterDefaults() {
	v.Register("required", required)
	v.Register("nonil", required)

	v.Register("len", length)
	v.Register("min", min)
	v.Register("max", max)

	v.Register("email", email)
	v.Register("uuid", uuidV4)
	v.Register("url", urlValue)

	v.Register("ip", ip)
	v.Register("ipv4", ipv4)
	v.Register("ipv6", ipv6)

	v.Register("alpha", alpha)
	v.Register("alphanum", alphanum)
	v.Register("num", numericString)
	v.Register("decimal", decimalString)

	v.Register("lower", lower)
	v.Register("upper", upper)
	v.Register("oneof", oneOf)

	v.Register("prefix", prefix)
	v.Register("suffix", suffix)
	v.Register("contains", contains)

	v.Register("date", dateValue)
	v.Register("time", timeValue)
	v.Register("datetime", dateTimeValue)
	v.Register("utc", utcValue)

	v.Register("gt", greaterThan)
	v.Register("gte", greaterThanOrEqual)
	v.Register("lt", lessThan)
	v.Register("lte", lessThanOrEqual)
}

func required(f Field) error {
	if f.IsNil {
		return fmt.Errorf("el campo [%s] es requerido", f.Path)
	}

	switch v := f.Value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("el campo [%s] es requerido", f.Path)
		}
	case []string:
		if len(v) == 0 {
			return fmt.Errorf("el campo [%s] es requerido", f.Path)
		}
	}

	return nil
}

func length(f Field) error {
	want, err := strconv.Atoi(f.Param)
	if err != nil {
		return fmt.Errorf("parámetro inválido para len en [%s]", f.Path)
	}

	switch v := f.Value.(type) {
	case string:
		if utf8.RuneCountInString(v) != want {
			return fmt.Errorf("el campo [%s] debe tener longitud [%d]", f.Path, want)
		}
	default:
		return fmt.Errorf("len solo soporta string en [%s]", f.Path)
	}

	return nil
}

func min(f Field) error {
	minVal, err := strconv.ParseFloat(f.Param, 64)
	if err != nil {
		return fmt.Errorf("parámetro inválido para min en [%s]", f.Path)
	}

	if s, ok := f.Value.(string); ok {
		if utf8.RuneCountInString(s) < int(minVal) {
			return fmt.Errorf("el campo [%s] debe tener mínimo [%d] caracteres", f.Path, int(minVal))
		}
		return nil
	}

	num, ok := asFloat(f.Value)
	if !ok {
		return fmt.Errorf("min solo soporta string o número en [%s]", f.Path)
	}

	if num < minVal {
		return fmt.Errorf("el campo [%s] debe ser mayor o igual a [%s]", f.Path, f.Param)
	}

	return nil
}

func max(f Field) error {
	maxVal, err := strconv.ParseFloat(f.Param, 64)
	if err != nil {
		return fmt.Errorf("parámetro inválido para max en [%s]", f.Path)
	}

	if s, ok := f.Value.(string); ok {
		if utf8.RuneCountInString(s) > int(maxVal) {
			return fmt.Errorf("el campo [%s] debe tener máximo [%d] caracteres", f.Path, int(maxVal))
		}
		return nil
	}

	num, ok := asFloat(f.Value)
	if !ok {
		return fmt.Errorf("max solo soporta string o número en [%s]", f.Path)
	}

	if num > maxVal {
		return fmt.Errorf("el campo [%s] debe ser menor o igual a [%s]", f.Path, f.Param)
	}

	return nil
}

func email(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("email solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] no es un correo válido", f.Path)
	}

	return nil
}

func uuidV4(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("uuid solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] debe ser un UUID v4 válido", f.Path)
	}

	return nil
}

func urlValue(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("url solo soporta string en [%s]", f.Path)
	}

	u, err := url.ParseRequestURI(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("el campo [%s] debe ser una URL válida", f.Path)
	}

	return nil
}

func ip(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("ip solo soporta string en [%s]", f.Path)
	}

	if net.ParseIP(s) == nil {
		return fmt.Errorf("el campo [%s] debe ser una IP válida", f.Path)
	}

	return nil
}

func ipv4(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("ipv4 solo soporta string en [%s]", f.Path)
	}

	parsed := net.ParseIP(s)
	if parsed == nil || parsed.To4() == nil {
		return fmt.Errorf("el campo [%s] debe ser una IPv4 válida", f.Path)
	}

	return nil
}

func ipv6(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("ipv6 solo soporta string en [%s]", f.Path)
	}

	parsed := net.ParseIP(s)
	if parsed == nil || parsed.To4() != nil {
		return fmt.Errorf("el campo [%s] debe ser una IPv6 válida", f.Path)
	}

	return nil
}

func alpha(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("alpha solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[A-Za-zÁÉÍÓÚáéíóúÑñÜü]+$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] solo debe contener letras", f.Path)
	}

	return nil
}

func alphanum(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("alphanum solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[A-Za-z0-9ÁÉÍÓÚáéíóúÑñÜü]+$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] solo debe contener letras y números", f.Path)
	}

	return nil
}

func numericString(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("num solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[0-9]+$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] solo debe contener números", f.Path)
	}

	return nil
}

func decimalString(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("decimal solo soporta string en [%s]", f.Path)
	}

	rx := regexp.MustCompile(`^[0-9]+\.[0-9]+$`)
	if !rx.MatchString(s) {
		return fmt.Errorf("el campo [%s] debe ser decimal", f.Path)
	}

	return nil
}

func lower(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("lower solo soporta string en [%s]", f.Path)
	}

	if s != strings.ToLower(s) {
		return fmt.Errorf("el campo [%s] debe estar en minúsculas", f.Path)
	}

	return nil
}

func upper(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("upper solo soporta string en [%s]", f.Path)
	}

	if s != strings.ToUpper(s) {
		return fmt.Errorf("el campo [%s] debe estar en mayúsculas", f.Path)
	}

	return nil
}

func oneOf(f Field) error {
	if strings.TrimSpace(f.Param) == "" {
		return fmt.Errorf("oneof requiere valores en [%s]", f.Path)
	}

	current := fmt.Sprint(f.Value)

	for _, allowed := range strings.Split(f.Param, ",") {
		if current == strings.TrimSpace(allowed) {
			return nil
		}
	}

	return fmt.Errorf("el campo [%s] debe ser uno de: [%s]", f.Path, f.Param)
}

func prefix(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("prefix solo soporta string en [%s]", f.Path)
	}

	if !strings.HasPrefix(s, f.Param) {
		return fmt.Errorf("el campo [%s] debe empezar con [%s]", f.Path, f.Param)
	}

	return nil
}

func suffix(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("suffix solo soporta string en [%s]", f.Path)
	}

	if !strings.HasSuffix(s, f.Param) {
		return fmt.Errorf("el campo [%s] debe terminar con [%s]", f.Path, f.Param)
	}

	return nil
}

func contains(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("contains solo soporta string en [%s]", f.Path)
	}

	if !strings.Contains(s, f.Param) {
		return fmt.Errorf("el campo [%s] debe contener [%s]", f.Path, f.Param)
	}

	return nil
}

func dateValue(f Field) error {
	switch v := f.Value.(type) {
	case string:
		if _, err := time.Parse(time.DateOnly, v); err != nil {
			return fmt.Errorf("el campo [%s] debe ser una fecha válida", f.Path)
		}
	case time.Time:
		return nil
	default:
		return fmt.Errorf("date solo soporta string o time.Time en [%s]", f.Path)
	}

	return nil
}

func timeValue(f Field) error {
	s, ok := f.Value.(string)
	if !ok {
		return fmt.Errorf("time solo soporta string en [%s]", f.Path)
	}

	if _, err := time.Parse(time.TimeOnly, s); err != nil {
		return fmt.Errorf("el campo [%s] debe ser una hora válida", f.Path)
	}

	return nil
}

func dateTimeValue(f Field) error {
	switch v := f.Value.(type) {
	case string:
		if _, err := time.Parse(time.DateTime, v); err != nil {
			return fmt.Errorf("el campo [%s] debe ser fecha y hora válida", f.Path)
		}
	case time.Time:
		return nil
	default:
		return fmt.Errorf("datetime solo soporta string o time.Time en [%s]", f.Path)
	}

	return nil
}

func utcValue(f Field) error {
	switch v := f.Value.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return fmt.Errorf("el campo [%s] debe tener formato UTC RFC3339", f.Path)
		}

		if t.Location() != time.UTC {
			return fmt.Errorf("el campo [%s] debe estar en UTC", f.Path)
		}

	case time.Time:
		if v.Location() != time.UTC {
			return fmt.Errorf("el campo [%s] debe estar en UTC", f.Path)
		}

	default:
		return fmt.Errorf("utc solo soporta string o time.Time en [%s]", f.Path)
	}

	return nil
}

func greaterThan(f Field) error {
	return compareNumber(f, func(a, b float64) bool { return a > b }, "mayor que")
}

func greaterThanOrEqual(f Field) error {
	return compareNumber(f, func(a, b float64) bool { return a >= b }, "mayor o igual que")
}

func lessThan(f Field) error {
	return compareNumber(f, func(a, b float64) bool { return a < b }, "menor que")
}

func lessThanOrEqual(f Field) error {
	return compareNumber(f, func(a, b float64) bool { return a <= b }, "menor o igual que")
}

func compareNumber(f Field, cmp func(float64, float64) bool, label string) error {
	value, ok := asFloat(f.Value)
	if !ok {
		return fmt.Errorf("[%s] solo soporta números", f.Tag)
	}

	param, err := strconv.ParseFloat(f.Param, 64)
	if err != nil {
		return fmt.Errorf("parámetro inválido para [%s] en [%s]", f.Tag, f.Path)
	}

	if !cmp(value, param) {
		return fmt.Errorf("el campo [%s] debe ser %s [%s]", f.Path, label, f.Param)
	}

	return nil
}

func asFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}
