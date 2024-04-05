### KCheck - Validación de Campos en Go

El paquete `kcheck` proporciona funciones para validar campos de datos en Go. Cada función de validación se invoca utilizando su etiqueta asociada. A continuación, se describen las funciones de validación disponibles junto con sus etiquetas correspondientes:

#### Funciones de Validación Disponibles

1. **uuidv4Func**
   - Etiqueta: `uuid`
   - Validación: Valida si un campo es un identificador UUIDv4 válido.

2. **noNilFunc**
   - Etiqueta: `nonil`
   - Validación: Valida si un campo no está vacío o contiene solo espacios en blanco.

3. **noSpacesStartAndEnd**
   - Etiqueta: `nosp`
   - Validación: Verifica que un campo de texto no comience ni termine con espacios.

4. **sword**
   - Etiqueta: `sword`
   - Validación: Valida si un campo solo contiene caracteres alfanuméricos y guiones bajos.

5. **sTextFunc**
   - Etiqueta: `stxt`
   - Validación: Valida un campo de texto según ciertos criterios, incluyendo restricciones sobre espacios y caracteres específicos.

6. **emailFunc**
   - Etiqueta: `email`
   - Validación: Valida si un campo contiene una dirección de correo electrónico válida.

7. **numberFunc**
   - Etiqueta: `num`
   - Validación: Valida si un campo contiene únicamente caracteres numéricos.

8. **decimalFunc**
   - Etiqueta: `decimal`
   - Validación: Valida si un campo contiene un número decimal válido.

9. **lenghtFunc**
   - Etiqueta: `len`
   - Validación: Valida si un valor tiene una longitud específica.

10. **maxLenghtFunc**
    - Etiqueta: `max`
    - Validación: Valida si un valor no excede una longitud máxima especificada.

11. **minLenghtFunc**
    - Etiqueta: `min`
    - Validación: Valida si un valor tiene al menos una longitud mínima especificada.

12. **regularExpression**
    - Etiqueta: `rgx`
    - Validación: Valida si un valor cumple con una expresión regular especificada.

13. **ipv4Func**
    - Etiqueta: `ip`
    - Validación: Valida si un valor es una dirección IPv4 válida.

14. **dateFunc**
    - Etiqueta: `date`
    - Validación: Valida si el valor de un campo es una fecha válida en el formato "2006-01-02".

15. **timeFunc**
    - Etiqueta: `time`
    - Validación: Valida si el valor de un campo es una hora válida en el formato "15:04:05".

16. **dateTimeFunc**
    - Etiqueta: `datetime`
    - Validación: Valida si el valor de un campo es una fecha y hora válidas en el formato "2006-01-02 15:04:05".

#### Uso

Cada función de validación se invoca utilizando su etiqueta correspondiente como se muestra en el siguiente ejemplo:

```go
    type Persona struct {
        NombrePersona string `json:"nombre_persona" chk:"nonil"`
        Edad          int
        Genero        string `chk:"nonil"`
        Direccion     string
    }

    func main() {
        persona := Persona{
            NombrePersona: "",
            Edad:          30,
            Genero:        "",
            Direccion:     "Calle Principal 123",
        }
        if err := kcheck.Valid(persona); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        fmt.Println(persona)
    }
```