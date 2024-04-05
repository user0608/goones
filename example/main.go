package main

import (
	"fmt"
	"os"

	"github.com/user0608/goones/kcheck"
)

type Persona struct {
	NombrePersona string `json:"nombre_persona" chk:"nonil"`
	Edad          int
	Genero        string `chk:"nonil"`
	Direccion     string
	Fecha         string `json:"fecha_nacimiento" chk:"date"`
}

func main() {
	persona := Persona{
		NombrePersona: "2",
		Edad:          30,
		Genero:        "ds",
		Direccion:     "Calle Principal 123",
		Fecha:         "2006-01-02",
	}
	if err := kcheck.Valid(persona); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(persona)
}
