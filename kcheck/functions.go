package kcheck

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ValidFunc func(atom Atom, args string) error
type MapFuncs map[string]ValidFunc

func uuidv4Func(atom Atom, _ string) error {
	const rgx = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	valid, err := regexp.MatchString(rgx, atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.uuidv4Func: %v\n", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	return fmt.Errorf("el campo [%s] debe ser un identificador uuid v4", atom.Name)
}
func numberFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+$", atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.numberFunc: %v\n", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	const message = "todos los caracteres del campo [%s] deben ser numéricos,el valor [%s] es invalido"
	return fmt.Errorf(message, atom.Value, atom.Name)
}
func decimalFunc(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[0-9]+.[0-9]+$", atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.Number: %v\n", err)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	message := "el campo [%s] es decimal, el valor [%s] no es es valido"
	return fmt.Errorf(message, atom.Value, atom.Name)
}
func sword(atom Atom, _ string) error {
	valid, err := regexp.MatchString("^[A-Za-z0-9_]*$", atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.sword: %v\n", err)
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
		log.Printf("ERROR: kcheck.calLens: %v\n", err)
		return 0, 0, ErrorKCHECK
	}
	valueLenght := len(value)
	return lenght, valueLenght, nil
}
func noNilFunc(atom Atom, _ string) error {
	lenght := len(atom.Value)
	if strings.TrimSpace(atom.Value) == "" {
		message := "el campo [%s] no puede quedar vacío"
		if lenght != 0 {
			message = "el campo [%s] no puede solo contener espacios en blanco"
		}
		return fmt.Errorf(message, atom.Name)
	}
	return nil
}
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
			const message = "el campo [%s] no puede contener ninguno de estos caracteres %s"
			return fmt.Errorf(message, atom.Name, denied)

		}
	}
	return nil
}
func emailFunc(atom Atom, _ string) error {
	match, err := regexp.MatchString(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`, atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.emailFunc: %v\n", err)
		return ErrorKCHECK
	}
	if !match {
		message := "el campo [%s] es de tipo correo, [%s] no es un correo válido"
		return fmt.Errorf(message, atom.Name, atom.Value)
	}
	return nil
}

func lenghtFunc(atom Atom, args string) error {
	vLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: kcheck.lenghtFunc field:`%s` args:`%s`: %v\n", atom.Name, args, err)
		return err
	}
	if valueLenght == vLen {
		return nil
	}
	message := "el número de caracteres del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, vLen, atom.Value, valueLenght)
}
func maxLenghtFunc(atom Atom, args string) error {
	maxLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: kcheck.maxLenghtFunc field:`%s` args:`%s`: %v\n", atom.Name, args, err)
		return err
	}
	if valueLenght <= maxLen {
		return nil
	}
	message := "el número de caracteres maximo del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, maxLen, atom.Value, valueLenght)
}
func minLenghtFunc(atom Atom, args string) error {
	minLen, valueLenght, err := calLens(atom.Value, args)
	if err != nil {
		log.Printf("ERROR: kcheck.minLenghtFunc field:[%s] args:[%s]: %v\n", atom.Name, args, err)
		return err
	}
	if valueLenght >= minLen {
		return nil
	}
	message := "el número de caracteres minimo del campo [%s] debe ser [%d], [%s] tiene [%d] caracteres"
	return fmt.Errorf(message, atom.Name, minLen, atom.Value, valueLenght)
}
func regularExpression(atom Atom, args string) error {
	valid, err := regexp.MatchString(args, atom.Value)
	if err != nil {
		log.Printf("ERROR: kcheck.regularExpression: [%v], en el campo [%s] con expresión [%s]\n", err, atom.Name, args)
		return ErrorKCHECK
	}
	if valid {
		return nil
	}
	message := "el valor [%s] en el campo [%s] es inválido, consulte con el administrador para más información"
	return fmt.Errorf(message, atom.Value, atom.Name)
}
