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

type PGCode string

const (
	Pg_InvalidLengthError       = PGCode("22001")
	Pg_DuplicateRecordError     = PGCode("23505")
	Pg_InvalidFormatError       = PGCode("23514")
	Pg_DependentRecordsError    = PGCode("23503")
	Pg_DataIntegrityError       = PGCode("23000")
	Pg_OperationFailedError     = PGCode("25000")
	Pg_InternalProblemError     = PGCode("26000")
	Pg_UnauthorizedAccessError  = PGCode("28000")
	Pg_TransactionError         = PGCode("2D000")
	Pg_NonexistentResourceError = PGCode("42P01")
	Pg_InvalidFieldValueError   = PGCode("22P02")
	Pg_InvalidJSONValueError    = PGCode("22032")
	Pg_NonNullableFieldsError   = PGCode("23502")
)

var pgErrcodes = map[PGCode]details{
	Pg_InvalidLengthError:       {"Verifique que los campos tengan la longitud correcta de caracteres.", http.StatusBadRequest, false},
	Pg_DuplicateRecordError:     {"El registro ya existe en la base de datos del sistema.", http.StatusBadRequest, false},
	Pg_InvalidFormatError:       {"Uno de los campos no tiene el formato correcto. Consulte con el administrador del sistema.", http.StatusBadRequest, false},
	Pg_DependentRecordsError:    {"Se encontraron otros registros dependientes. No podemos realizar ninguna acción mientras estas relaciones existan.", http.StatusBadRequest, false},
	Pg_DataIntegrityError:       {"Operación restringida debido a un problema de integridad en los datos. Consulte la documentación.", http.StatusBadRequest, false},
	Pg_OperationFailedError:     {"No se pudieron completar las operaciones. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	Pg_InternalProblemError:     {"Hubo un problema interno. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	Pg_UnauthorizedAccessError:  {"Acceso restringido. No podemos realizar la operación.", http.StatusUnauthorized, true},
	Pg_TransactionError:         {"Hubo un problema al realizar la transacción. Por favor, informe la incidencia al equipo técnico.", http.StatusInternalServerError, true},
	Pg_NonexistentResourceError: {"El registro o recurso al que intenta acceder no existe.", http.StatusBadRequest, false},
	Pg_InvalidFieldValueError:   {"El formato o representación de uno de los valores de campo no cumple con los requerimientos.", http.StatusBadRequest, false},
	Pg_InvalidJSONValueError:    {"El valor asignado a uno de los campos de tipo JSON no cumple con los requerimientos.", http.StatusBadRequest, false},
	Pg_NonNullableFieldsError:   {"Hay campos que no deberían ser nulos. Consulte la documentación o al administrador del sistema.", http.StatusBadRequest, false},
}

func AddPgErrs(pgerrcode PGCode, message string, httpcode int, loggable bool) {
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
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "record not found") {
		return newError(err, ErrRecordNotFound, http.StatusBadRequest)
	}

	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.Code == "23503" && strings.Contains(pgerr.Message, "insert or update") {
			return newError(err, message23503, http.StatusBadRequest)
		}
		mutex.RLock()
		defer mutex.RUnlock()
		state, ok := pgErrcodes[PGCode(pgerr.Code)]
		if ok {
			if state.loggable || devmode == "1" {
				return newError(err, state.message, state.httcode)
			}
			return newError(nil, state.message, state.httcode)
		}
	}
	return newError(err, ErrDatabase, http.StatusInternalServerError)
}

func IsPgErrCode(err error, code PGCode) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		return PGCode(pgerr.Code) == code
	}
	return false
}
