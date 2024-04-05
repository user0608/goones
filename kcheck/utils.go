package kcheck

import (
	"log/slog"
	"regexp"
	"strings"
)

// StandardSpace elimina los espacios innecesarios entre palabras en una cadena de texto,
// dejando solamente un espacio entre cada palabra.
func StandardSpace(s string) string {
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// SplitKeyValue recibe un cadena de texto `len=89`, separado por un `=`
// esta función devuelve (valido,key,valor)
func SplitKeyValue(s string) (bool, string, string) {
	ok, err := regexp.MatchString("^[a-zA-Z]+=(.)+$", s)
	if err != nil {
		slog.Warn("kcheck.SplitKeyValue", "error", err)
		return false, "", ""
	}
	if !ok {
		return false, "", ""
	}
	values := strings.Split(s, "=")
	return true, values[0], values[1]
}
