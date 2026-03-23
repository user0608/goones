package types

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

var ErrInvalidUUIDArrayInput = errors.New("invalid input for UUIDArray")

var _ json.Unmarshaler = (*UUIDArray)(nil)
var _ json.Marshaler = (*UUIDArray)(nil)

func (sa *UUIDArray) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*sa = []uuid.UUID{}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		res := make([]uuid.UUID, len(arr))
		for i, v := range arr {
			u, err := uuid.Parse(v)
			if err != nil {
				return err
			}
			res[i] = u
		}
		*sa = res
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		u, err := uuid.Parse(single)
		if err != nil {
			return err
		}
		*sa = []uuid.UUID{u}
		return nil
	}

	return ErrInvalidUUIDArrayInput
}

func (sa UUIDArray) MarshalJSON() ([]byte, error) {
	if sa == nil {
		return []byte("null"), nil
	}

	out := make([]string, len(sa))
	for i, v := range sa {
		out[i] = v.String()
	}

	return json.Marshal(out)
}

func (sa UUIDArray) Unique() UUIDArray {
	seen := make(map[uuid.UUID]struct{}, len(sa))
	out := make([]uuid.UUID, 0, len(sa))
	for _, v := range sa {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
