package kcheck

import (
	"log"
	"regexp"
	"strings"
)

// SplitCamelCase recibe un texto camelcase y lo separa por espacios
func SplitCamelCase(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1} ${2}")
	}
	return s
}

// StandardSpace elimina los espacios innecesarios entre palabras
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
		log.Println("ERROR:", err)
		return false, "", ""
	}
	if !ok {
		return false, "", ""
	}
	values := strings.Split(s, "=")
	return true, values[0], values[1]
}
