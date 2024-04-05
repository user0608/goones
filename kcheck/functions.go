package kcheck

import (
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ValidFunc func(atom Atom, args string) error
type MapFuncs map[string]ValidFunc

// uuidv4Func valida si un campo es un identificador UUIDv4 válido.
// Retorna un error si el campo no cumple con el formato UUIDv4.
func uuidv4Func(atom Atom, _ string) error {
	const rgx = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	valid, err := regexp.MatchString(rgx, atom.Value)
	if err != nil {
		slog.Warn("kcheck.uuidv4Func", "error", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	return fmt.Errorf("el campo [%s] debe ser un identificador uuid v4", atom.Name)
}

// numberFunc valida si un campo contiene únicamente caracteres numéricos.
// Retorna un error si el campo contiene caracteres que no son números.
func numberFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+$", atom.Value)
	if err != nil {
		slog.Warn("kcheck.numberFunc", "error", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	const message = "todos los caracteres del campo [%s] deben ser numéricos,el valor [%s] es invalido"
	return fmt.Errorf(message, atom.Value, atom.Name)
}

// decimalFunc valida si un campo contiene un número decimal válido.
// Retorna un error si el campo no contiene un número decimal válido.
// Ejemplos de números decimales válidos:
// 1. "123.45"
// 2. "0.123"
// 3. "456.789"
// 4. "987654.321"
// 5. "10.00"
// 6. "3.14159"
func decimalFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+.[0-9]+$", atom.Value)
	if err != nil {
		slog.Warn("kcheck.decimalFunc", "error", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	message := "el campo [%s] es decimal, el valor [%s] no es es valido"
	return fmt.Errorf(message, atom.Value, atom.Name)
}

// sword valida si un campo solo contiene caracteres alfanuméricos y guiones bajos.
// Retorna un error si el campo contiene caracteres no permitidos.
func sword(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[A-Za-z0-9_]*$", atom.Value)
	if err != nil {
		slog.Warn("kcheck.sword", "error", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	message := "en el campo [%s] solo están permitidos caracteres numéricos y alfabéticos"
	return fmt.Errorf(message, atom.Name)
}

// calLens retorna el valor slen convertido en int, el numero de caracteres del value y error en caso exista
// Utilizado por Lenght, MaxLenght, MinLenght
func calLens(value string, slen string) (int, int, error) {
	lenght, err := strconv.Atoi(slen)
	if err != nil {
		slog.Warn("kcheck.calLens", "error", err)
		return 0, 0, ErrorKCHECK
	}
	valueLenght := len(value)
	return lenght, valueLenght, nil
}

// noNilFunc valida si un campo está vacío o contiene solo espacios en blanco.
// Retorna un error si el campo está vacío o contiene solo espacios en blanco.
func noNilFunc(atom Atom, _ string) error {
	lenght := len(atom.Value)
	if strings.TrimSpace(atom.Value) == "" {
		message := "El campo [%s] no puede quedar vacío"
		if lenght != 0 {
			message = "el campo [%s] no puede solo contener espacios en blanco"
		}
		return fmt.Errorf(message, atom.Name)
	}
	return nil
}

// noSpacesStartAndEnd verifica que un campo de texto no comience ni termine con espacios.
// Retorna un error si el campo comienza o termina con espacios.
func noSpacesStartAndEnd(atom Atom, _ string) error {
	matchStartSpace, _ := regexp.MatchString("^( .)", atom.Value)
	if matchStartSpace {
		message := "el campo [%s] no puede contener espacios al inicio"
		return fmt.Errorf(message, atom.Name)
	}
	matchEndSpace, _ := regexp.MatchString("(. )$", atom.Value)
	if matchEndSpace {
		message := "el campo [%s] no puede contener espacios al final"
		return fmt.Errorf(message, atom.Name)
	}
	return nil
}

// sTextFunc valida un campo de texto según ciertos criterios.
// Comprueba que el campo no comience o termine con espacios,
// no contenga palabras separadas por más de dos espacios consecutivos
// y no contenga caracteres específicos.
// Retorna un error si alguna de estas condiciones no se cumple.
func sTextFunc(atom Atom, args string) error {
	denied := "!\"#$%&'()*+,./:;<=>?@[\\]^_}{~|"
	if err := noSpacesStartAndEnd(atom, args); err != nil {
		return err
	}
	match, _ := regexp.MatchString("( ){3}", atom.Value)
	if match {
		const message = "el campo [%s] no puede tener palabras separadas por más de 2 espacios"
		return fmt.Errorf(message, atom.Name)
	}
	for _, c := range atom.Value {
		if strings.ContainsRune(denied, c) {
			const message = "el campo [%s] no puede contener ninguno de estos caracteres [%s]"
			return fmt.Errorf(message, atom.Name, denied)

		}
	}
	return nil
}

// emailFunc valida si un campo contiene una dirección de correo electrónico válida.
// Retorna un error si el campo no contiene una dirección de correo electrónico válida.
func emailFunc(atom Atom, _ string) error {
	match, err := regexp.MatchString(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`, atom.Value)
	if err != nil {
		slog.Warn("kcheck.emailFunc", "error", err)
		return ErrorKCHECK
	}
	if !match {
		message := "el campo [%s] es de tipo correo, [%s] no es un correo válido"
		return fmt.Errorf(message, atom.Name, atom.Value)
	}
	return nil
}

// lenghtFunc valida si un valor tiene una longitud específica.
// Retorna un error si la longitud del valor no coincide con la longitud específica.
func lenghtFunc(atom Atom, args string) error {
	vLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		slog.Warn("kcheck.lenghtFunc", "error", err, "field", atom.Name, "args", args)
		return err
	}
	if valueLenght == vLen {
		return nil
	}
	message := "el número de caracteres del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, vLen, atom.Value, valueLenght)
}

// maxLenghtFunc valida si un valor no excede una longitud máxima especificada.
// Retorna un error si la longitud del valor es mayor que la longitud máxima especificada.
func maxLenghtFunc(atom Atom, args string) error {
	maxLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		slog.Warn("kcheck.maxLenghtFunc", "error", err, "field", atom.Name, "args", args)
		return err
	}
	if valueLenght <= maxLen {
		return nil
	}
	message := "el número de caracteres maximo del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, maxLen, atom.Value, valueLenght)
}

// minLenghtFunc valida si un valor tiene al menos una longitud mínima especificada.
// Retorna un error si la longitud del valor es menor que la longitud mínima especificada.
func minLenghtFunc(atom Atom, args string) error {
	minLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		slog.Warn("kcheck.minLenghtFunc", "error", err, "field", atom.Name, "args", args)
		return err
	}
	if valueLenght >= minLen {
		return nil
	}
	message := "el número de caracteres minimo del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, minLen, atom.Value, valueLenght)
}

// regularExpression valida si un valor cumple con una expresión regular especificada.
// Retorna un error si el valor no coincide con la expresión regular.
func regularExpression(atom Atom, args string) error {
	valid, err := regexp.MatchString(args, atom.Value)
	if err != nil {
		slog.Warn("kcheck.regularExpression", "error", err, "field", atom.Name, "args", args)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	message := "el valor [%s] en el campo [%s] es inválido, consulte con el administrador para más información"
	return fmt.Errorf(message, atom.Value, atom.Name)
}

// ipv4Func valida si un valor es una dirección IPv4 válida.
// Retorna un error si el valor no es una dirección IPv4 válida.
func ipv4Func(atom Atom, _ string) error {
	ip := net.ParseIP(atom.Value)
	if ip == nil || ip.To4() == nil {
		message := "el campo [%s] no es una dirección IPv4 válida"
		return fmt.Errorf(message, atom.Name)
	}
	return nil
}

// dateFunc valida si el valor de un campo es una fecha válida en el formato "2006-01-02".
// Retorna un error si el valor no es una fecha válida.
func dateFunc(atom Atom, _ string) error {
	if err := noNilFunc(atom, ""); err != nil {
		return err
	}
	_, err := time.Parse(time.DateOnly, atom.Value)
	if err != nil {
		message := "el valor [%s] en el campo [%s] no es una fecha válida"
		return fmt.Errorf(message, atom.Value, atom.Name)
	}
	return nil
}

// timeFunc valida si el valor de un campo es una hora válida en el formato "15:04:05".
// Retorna un error si el valor no es una hora válida.
func timeFunc(atom Atom, _ string) error {
	if err := noNilFunc(atom, ""); err != nil {
		return err
	}
	_, err := time.Parse(time.TimeOnly, atom.Value)
	if err != nil {
		message := "el valor [%s] en el campo [%s] no es una hora válida"
		return fmt.Errorf(message, atom.Value, atom.Name)
	}
	return nil
}

// dateTimeFunc valida si el valor de un campo es una fecha y hora válidas en el formato "2006-01-02 15:04:05".
// Retorna un error si el valor no es una fecha y hora válidas.
func dateTimeFunc(atom Atom, _ string) error {
	if err := noNilFunc(atom, ""); err != nil {
		return err
	}
	_, err := time.Parse(time.DateTime, atom.Value)
	if err != nil {
		message := "el valor [%s] en el campo [%s] no es una fecha y hora válidas"
		return fmt.Errorf(message, atom.Value, atom.Name)
	}
	return nil
}
