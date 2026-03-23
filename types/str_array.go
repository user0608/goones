package types

import (
	"encoding/json"
	"errors"
	"strings"
)

type StrArray []string

var ErrInvalidStrArrayInput = errors.New("invalid input for StrArray")

var _ json.Unmarshaler = (*StrArray)(nil)
var _ json.Marshaler = (*StrArray)(nil)

func (sa *StrArray) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*sa = StrArray{}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		if arr == nil {
			*sa = StrArray{}
		} else {
			*sa = arr
		}
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*sa = []string{single}
		return nil
	}

	return ErrInvalidStrArrayInput
}

func (sa StrArray) MarshalJSON() ([]byte, error) {
	if sa == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(sa))
}

func (sa StrArray) Trimmed() StrArray {
	out := make([]string, len(sa))
	for i, v := range sa {
		out[i] = strings.TrimSpace(v)
	}
	return out
}

func (sa StrArray) NonEmpty() StrArray {
	out := make([]string, 0, len(sa))
	for _, v := range sa {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func (sa StrArray) Unique() StrArray {
	seen := make(map[string]struct{}, len(sa))
	out := make([]string, 0, len(sa))
	for _, v := range sa {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
