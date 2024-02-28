package kcheck

import (
	"errors"

	"log/slog"
	"reflect"
	"strings"
	"sync"
)

var TAG = "chk"

var ErrorKCHECK = errors.New("error unexpected, kcheck")
var mutex sync.RWMutex
var tagFuncs = MapFuncs{
	"uuid":    uuidv4Func,
	"nonil":   noNilFunc,
	"nosp":    noSpacesStartAndEnd,
	"sword":   sword,
	"stxt":    sTextFunc,
	"email":   emailFunc,
	"num":     numberFunc,
	"decimal": decimalFunc,
	"len":     lenghtFunc,
	"max":     maxLenghtFunc,
	"min":     minLenghtFunc,
	"rgx":     regularExpression,
}

// OmitFields lista de campos que no se tomaran en cuanta al realizar la verificación
type Fields []string

/*
AddFunc permite registrar una nueva función personalizada, la cual será asociada
al tagKey indicado, si le takKey ya existe, este será remplazado, por ejemplo si
se usa el tagKey `num` este remplaza al existente. La función ValidFunc recibe
como primer parámetro un objeto con los datos del campo a verificar, incluye el
nombre y el valor, y como segundo parámetro recibe el valor después del `=` del
tagKey, por ejemplo el tag es `len` y este recibe un valor `len=10` el 10 es enviado
como segundo parámetro en formato string.
Nota: importante que el registro de nuevas funciona no se haga en tiempo de ejecución,
esto podría generar problemas de acceso por parte de las gorutines
*/
type TagParamExtractor interface {
	GetTagValue(fieldName string) (value string, ok bool)
}
type paramExtractor struct {
	rValue reflect.Value
	rType  reflect.Type
}

func (ex *paramExtractor) GetTagValue(fieldName string) (value string, ok bool) {
	rsf, found := ex.rType.FieldByName(fieldName)
	if !found {
		return
	}
	if rsf.Type.Kind() == reflect.String {
		value = rsf.Tag.Get(TAG)
		ok = true
	}
	return
}

func AddFunc(tagKey string, f ValidFunc) {
	mutex.Lock()
	tagFuncs[tagKey] = f
	mutex.Lock()
}
func getFunc(tagKey string) (ValidFunc, bool) {
	mutex.RLock()
	f, ok := tagFuncs[tagKey]
	mutex.RUnlock()
	return f, ok
}

func (of *Fields) isContain(field string) bool {
	for _, v := range *of {
		if v == field {
			return true
		}
	}
	return false
}

func Valid(i interface{}, skips ...string) error {
	return valid(i, skips, true)
}
func ValidSelect(i interface{}, selecteds ...string) error {
	return valid(i, selecteds, false)
}
func reflectValueAndType(i interface{}) (*reflect.Value, *reflect.Type, error) {
	var rValue reflect.Value
	rType := reflect.TypeOf(i)
	if rType == nil {
		slog.Warn("kcheck: nil value was received")
		return nil, nil, ErrorKCHECK
	}
	switch rType.Kind() {
	case reflect.Struct:
		rValue = reflect.ValueOf(i)
	case reflect.Ptr:
		if rType.Elem().Kind() == reflect.Struct {
			rValue = reflect.ValueOf(i).Elem()
			rType = rType.Elem()
		} else {
			slog.Warn("kcheck: invalid type", "type", rType)
			return nil, nil, ErrorKCHECK
		}
	}
	return &rValue, &rType, nil
}

func valid(i interface{}, filds Fields, isOmit bool) error {
	rValue, rType, err := reflectValueAndType(i)
	if err != nil {
		return err
	}
	for i := 0; i < (*rType).NumField(); i++ {
		rsf := (*rType).Field(i)
		rv := rValue.Field(i)
		if rsf.Type.Kind() == reflect.String {
			tagValues := rsf.Tag.Get(TAG)
			if isOmit {
				if tagValues == "" || filds.isContain(rsf.Name) {
					continue
				}
			} else {
				if !filds.isContain(rsf.Name) {
					continue
				}
			}
			atom := Atom{Name: SplitCamelCase(rsf.Name), Value: rv.String()}
			if err := ValidTarget(tagValues, atom); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidTarget(tags string, atom Atom) error {
	tags = StandardSpace(tags)
	keys := strings.Split(tags, " ")
	for _, key := range keys {
		if f, ok := getFunc(key); ok {
			if err := f(atom, ""); err != nil {
				return err
			}
		} else {
			valid, fkey, keyValues := SplitKeyValue(key)
			if valid {
				if ff, okk := getFunc(fkey); okk {
					if err := ff(atom, keyValues); err != nil {
						return err
					}
				} else {
					slog.Warn("kcheck: tag value invalid", "tag", key, "field", atom.Name)
					return ErrorKCHECK
				}
			} else {
				slog.Warn("kcheck: tag value invalid", "tag", key, "field", atom.Name)
				return ErrorKCHECK
			}
		}
	}
	return nil
}

func BuildTagParamExtractor(i interface{}) (TagParamExtractor, error) {
	rValue, rType, err := reflectValueAndType(i)
	if err != nil {
		return nil, err
	}
	return &paramExtractor{rValue: *rValue, rType: *rType}, nil
}
