package types_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user0608/goones/types"
)

func TestOnlyDateMarshalJSON(t *testing.T) {
	d := types.NewDateOnly(time.Date(2025, 8, 22, 14, 0, 0, 0, time.UTC))
	data, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-08-22"`, string(data))
}

func TestOnlyDateUnmarshalJSON(t *testing.T) {
	var d types.DateOnly
	err := json.Unmarshal([]byte(`"2025-08-22"`), &d)
	assert.NoError(t, err)
	assert.Equal(t, "2025-08-22", d.String())
}

func TestOnlyDateUnmarshalJSON_Null(t *testing.T) {
	var d types.DateOnly
	err := json.Unmarshal([]byte(`null`), &d)
	assert.NoError(t, err)
	assert.True(t, d.Time.IsZero())
}

func TestOnlyDateUnmarshalParam(t *testing.T) {
	var d types.DateOnly
	err := d.UnmarshalParam(`"2025-08-22"`)
	assert.NoError(t, err)
	assert.Equal(t, "2025-08-22", d.String())
}

func TestOnlyDateGobEncodeDecode(t *testing.T) {
	d := types.NewDateOnly(time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC))
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(d)
	assert.NoError(t, err)

	var decoded types.DateOnly
	err = gob.NewDecoder(&buf).Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, d.Time, decoded.Time)
}

func TestOnlyDateValue(t *testing.T) {
	d := types.NewDateOnly(time.Date(2025, 8, 22, 14, 0, 0, 0, time.UTC))
	val, err := d.Value()
	assert.NoError(t, err)

	strVal, ok := val.(string)
	assert.True(t, ok)
	assert.Equal(t, "2025-08-22", strVal)
}

func TestDateOnly_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "time.Time",
			input:    time.Date(2025, 8, 22, 14, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "*time.Time",
			input:    ptrTime(time.Date(2025, 8, 21, 23, 59, 59, 0, time.UTC)),
			expected: time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "DateOnly",
			input:    types.NewDateOnly(time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)),
			expected: time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "*DateOnly",
			input:    ptrDO(types.NewDateOnly(time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC))),
			expected: time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "string with DateOnly layout",
			input:    "2025-08-18",
			expected: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "string with DateTime fallback",
			input:    "2025-08-17T00:00:00Z",
			expected: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "[]byte with DateOnly layout",
			input:    []byte("2025-08-16"),
			expected: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "empty string",
			input:    "",
			expected: time.Time{},
		},
		{
			name:     "nil",
			input:    nil,
			expected: time.Time{},
		},
		{
			name:    "invalid string",
			input:   "invalid-date",
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   12345,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var do types.DateOnly
			err := do.Scan(tt.input)

			if tt.wantErr {
				assert.Error(t, err, "esperaba error pero fue nil")
			} else {
				assert.NoError(t, err, "no esperaba error pero fue: %v", err)
				assert.Equal(t, tt.expected, do.Time, "valor incorrecto")
			}
		})
	}
}

func ptrDO(d types.DateOnly) *types.DateOnly {
	return &d
}

func TestOnlyDateOnlyDateString(t *testing.T) {
	d := types.NewDateOnly(time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, "2025-08-22", d.String())
}

func TestDateOnly_ToTimeInLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/Lima")
	assert.NoError(t, err)

	input := time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC)
	dateOnly := types.NewDateOnly(input)

	result := dateOnly.ToTimeInLocation(loc)

	expected := time.Date(2025, 8, 22, 0, 0, 0, 0, loc)

	assert.Equal(t, expected.Format(time.RFC3339), result.Format(time.RFC3339))
}
