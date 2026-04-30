# kcheck

Validador ligero para structs en Go basado en tags.

## Instalación

```bash
go get github.com/tu-usuario/kcheck
```

## Uso básico

```go
package main

import (
    "fmt"
    "time"

    "github.com/tu-usuario/kcheck"
)

type Address struct {
    City string `chk:"required min=3"`
}

type User struct {
    Name      string    `chk:"required min=2 max=50"`
    Email     string    `chk:"required email"`
    Age       int       `chk:"gte=18 lte=120"`
    Code      *string   `chk:"required upper len=6"`
    Status    string    `chk:"required oneof=active,inactive,pending"`
    Website   *string   `chk:"url"`
    CreatedAt time.Time `chk:"required utc"`
    Address   Address
}

func main() {
    code := "ABC123"
    url := "https://example.com"

    user := User{
        Name:      "Kevin",
        Email:     "kevin@example.com",
        Age:       30,
        Code:      &code,
        Status:    "active",
        Website:   &url,
        CreatedAt: time.Now().UTC(),
        Address:   Address{City: "Lima"},
    }

    err := kcheck.Valid(user)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("valid")
}
```

## Skip campos

```go
err := kcheck.Valid(user, "Email")
```

## Select campos

```go
err := kcheck.ValidSelect(user, "Name")
```

## Tags disponibles

### Requerido
- required
- nonil

### Longitud / tamaño
- len=n
- min=n
- max=n

### Comparadores
- gt=n
- gte=n
- lt=n
- lte=n

### Strings
- alpha
- alphanum
- num
- decimal
- lower
- upper

### Formato
- email
- uuid
- url
- ip
- ipv4
- ipv6

### Strings avanzados
- oneof=a,b,c
- prefix=x
- suffix=x
- contains=x

### Fechas
- date → 2006-01-02
- time → 15:04:05
- datetime → 2006-01-02 15:04:05
- utc → 2026-04-30T15:04:05Z

## Tipos soportados

- string
- *string
- int, uint, float
- bool
- time.Time
- punteros
- structs anidados

## Ejemplo de error

Name: el campo es requerido; Email: el campo no es un correo válido

## Custom validator

```go
v := kcheck.New()

v.Register("startsx", func(f kcheck.Field) error {
    if s, ok := f.Value.(string); ok {
        if !strings.HasPrefix(s, "x") {
            return fmt.Errorf("debe empezar con x")
        }
    }
    return nil
})

type DTO struct {
    Code string `chk:"startsx"`
}

err := v.Struct(DTO{Code: "abc"})
```