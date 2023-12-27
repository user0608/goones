package errs

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
)

type details struct {
	message  string
	httcode  int
	loggable bool
}

var mutex sync.RWMutex

var pgErrcodes = map[string]details{
	"22001": {"Verifique que los campos tengan la longitud correcta de caracteres.", http.StatusBadRequest, false},
	"23505": {"El registro ya existe en la base de datos del sistema.", http.StatusBadRequest, false},
	"23514": {"Uno de los campos no tiene el formato correcto. Consulte con el administrador del sistema.", http.StatusBadRequest, false},
	"23503": {"Se encontraron otros registros dependientes. No podemos realizar ninguna acción mientras estas relaciones existan.", http.StatusBadRequest, false},
	"23000": {"Operación restringida debido a un problema de integridad en los datos. Consulte la documentación.", http.StatusBadRequest, false},
	"25000": {"No se pudieron completar las operaciones. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	"26000": {"Hubo un problema interno. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	"28000": {"Acceso restringido. No podemos realizar la operación.", http.StatusUnauthorized, true},
	"2D000": {"Hubo un problema al realizar la transacción. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	"42P01": {"El registro o recurso al que intenta acceder no existe.", http.StatusBadRequest, false},
	"22P02": {"El formato o representación de uno de los valores de campo no cumple con los requerimientos.", http.StatusBadRequest, false},
	"22032": {"El valor asignado a uno de los campos de tipo JSON no cumple con los requerimientos.", http.StatusBadRequest, false},
	"23502": {"Hay campos que no deberían ser nulos. Consulte la documentación o al administrador del sistema.", http.StatusBadRequest, false},
}

func AddPgErrs(pgerrcode, message string, httpcode int, loggable bool) {
	mutex.Lock()
	defer mutex.Unlock()
	pgErrcodes[pgerrcode] = details{message, httpcode, loggable}
}

const (
	ErrInvalidRequestBody          = "La estructura de información enviada es inválida. Por favor, revise la documentación y vuelva a intentar."
	ErrInvalidQueryParam           = "Los parámetros de consulta son inválidos. Favor de revisar la documentación y volver a intentar."
	ErrAuthorizationHeaderNotFound = "La cabecera con el token de utilización no fue encontrada. La operación fue rechazada."
	ErrInvalidToken                = "El token que está utilizando no es válido o ha caducado. Contáctese con el equipo técnico."
	ErrSigningTokenString          = "El token que está utilizando no es genuino. Contáctese con el equipo técnico."
	ErrDatabase                    = "La operación no se pudo realizar debido a algún problema. Contáctese con el equipo técnico."

	ErrRecordNotFound        = "El registro buscado no fue encontrado."
	ErrCreating              = "No se pudo realizar la operación de registro."
	ErrUpdating              = "No se pudo realizar la operación de actualización."
	ErrUserOrPasswordInvalid = "Usuario o contraseña incorrectos."
	ErrIDNotFound            = "Parámetro ID o identificador no encontrado."
	ErrCodeNotFound          = "Parámetro código no encontrado."
	ErrNameNotFound          = "Parámetro nombre no encontrado."
	ErrNotFound              = "No se pudo encontrar ningún recurso asociado a esta consulta."
	ErrGeneric               = "Hubo un error inesperado. Favor de reportar la incidencia al equipo técnico."
	ErrInternal              = ErrGeneric
)
const message23503 = "No se puede realizar la operación debido a asociaciones incompatibles. Asegúrese de que los valores relacionados existan antes de intentar el registro."

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
		mutex.RLock()
		defer mutex.RUnlock()
		state, ok := pgErrcodes[pgerr.Code]
		if ok {
			if state.loggable || devmode == "1" {
				return newErrf(err, state.message, state.httcode)
			}
			return newErrf(nil, state.message, state.httcode)
		}
	}
	return newErrf(err, ErrDatabase, http.StatusInternalServerError)
}
