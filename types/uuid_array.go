package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (sa *UUIDArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	switch obj := jsonObj.(type) {
	case string:
		uuidValue, err := uuid.Parse(obj)
		if err != nil {
			return err
		}
		*sa = UUIDArray([]uuid.UUID{uuidValue})
		return nil
	case nil:
		*sa = UUIDArray([]uuid.UUID{})
		return nil
	case []interface{}:
		s := make([]uuid.UUID, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrUnsupportedType
			}
			uuidValue, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			s = append(s, uuidValue)
		}
		*sa = UUIDArray(s)
		return nil
	}
	return ErrUnsupportedType
}
