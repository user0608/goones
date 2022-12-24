package errs

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgconn"
)

type details struct {
	message  string
	httcode  int
	loggable bool
}

var mutex sync.RWMutex
var pgErrcodes = map[string]details{
	"22001": {"verifique que los campos tenga la longitud correcta de caracteres", http.StatusBadRequest, false},
	"23505": {"el registro ya existe en la base de datos del sistema", http.StatusBadRequest, false},
	"23514": {"uno de los campos no tiene el formato correcto, consulte con el administrador del sistema", http.StatusBadRequest, false},
	"23503": {"se encontró otros registros dependientes, no podemos realizar ninguna acción mientras estas relaciones existan", http.StatusBadRequest, false},
	"23000": {"operación restringida, problema de integridad en los datos, consulte documentación", http.StatusBadRequest, false},
	"25000": {"no se pudo completar las operaciones, favor de reporte la incidencia al equipo técnico", http.StatusInternalServerError, true},
	"26000": {"hubo un problema interno, favor de reporte la incidencia al equipo técnico", http.StatusInternalServerError, true},
	"28000": {"Acceso restringido, no podemos realizar la operación", http.StatusUnauthorized, true},
	"2D000": {"hubo un problema al realizar la transacción, favor de reporte la incidencia al equipo técnico", http.StatusInternalServerError, true},
	"42P01": {"el registro o recurso al que intenta acceder, no existe", http.StatusBadRequest, false},
	"22P02": {"el formato o representación de uno de los valores de campo no cumple con los requerimientos", http.StatusBadRequest, false},
	"22032": {"el valor asignado a uno de los campos de tipo JSON, no cumple con los requerimientos", http.StatusBadRequest, false},
	"23502": {"hay campos que no deberían de ser nullables, consulte la documentación o al administrador del sistema", http.StatusBadRequest, false},
}

func AddPgErrs(pgerrcode, message string, httpcode int, loggable bool) {
	mutex.Lock()
	defer mutex.Unlock()
	pgErrcodes[pgerrcode] = details{message, httpcode, loggable}
}

const (
	ErrInvalidRequestBody          = "la estructura de información enviada es inválida, revisar la documentación y volver a intentar"
	ErrInvalidQueryParam           = "los parámetros de consulta son inválidos, favor de revisar la documentación y volver a intentar"
	ErrAuthorizationHeaderNotFound = "la cabecera con el token de utilización no fue encontrado, la operación fue rechazada"
	ErrInvalidToken                = "el token que está utilizando no es válido o está caducado, contáctese con el equipo técnico"
	ErrSigningTokenString          = "el token que está utilizando no es genuino, contáctese con el equipo técnico"
	ErrDatabase                    = "la operacion no se pudo realizar de bebido a algún problema, contáctese con el equipo técnico"

	ErrRecordNotFaund        = "el registro buscando no fue encontrado"
	ErrCreating              = "no se puedo realizar la operacion de registro"
	ErrUpdating              = "no se puedo realizar la operacion de actualización"
	ErrUserOrPasswordInvalid = "usuario o contraseña son incorrectos"
	ErrIDNotFound            = "parametro ID o identificador no encontrado"
	ErrCodeNotFound          = "parametro codigo no fue encontrado"
	ErrNameNotFound          = "parametro name no encontrado"
	ErrNotFound              = "no se pudo encontrar ningún recurso asociado a esta consulta"
	ErrGeneric               = "hubo un error, no esperado, favor de reporte la incidencia al equipo técnico"
	ErrInternal              = ErrGeneric
)

const message23503 = "asociación incompatible, verifique la existencia de las identidades relacionadas al registro"

var devmode = os.Getenv("DEV_MODE")

func Devmode() { devmode = "1" }

// Pgf Encapsula el error retornado por PostgreSQL y prepara los mensajes,
// código de error para la respuesta al cliente HTTP, esta respuesta es
// gestionada por `answer`
func Pgf(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.Code == "23503" && strings.Contains(pgerr.Message, "insert or update") {
			return newErrf(err, message23503, http.StatusBadRequest)
		}
		d, ok := pgErrcodes[pgerr.Code]
		if ok {
			if d.loggable || devmode == "1" {
				return newErrf(err, d.message, d.httcode)
			}
			return newErrf(nil, d.message, d.httcode)
		}
	}
	return newErrf(err, ErrDatabase, http.StatusInternalServerError)
}
