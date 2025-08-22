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

func TestMarshalJSON(t *testing.T) {
	tm := time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC)
	dt := types.NewDateTimeOnly(tm)

	data, err := json.Marshal(dt)
	assert.NoError(t, err)
	assert.Equal(t, `"`+tm.Format(time.DateTime)+`"`, string(data))
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{`"2025-08-22 15:30:00"`, time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC)},
		{`null`, time.Time{}},
	}

	for _, tt := range tests {
		var dt types.DateTimeOnly
		err := json.Unmarshal([]byte(tt.input), &dt)
		assert.NoError(t, err)
		assert.True(t, dt.Time.Equal(tt.expected), "input: %s", tt.input)
	}
}

func TestUnmarshalParam(t *testing.T) {
	param := `"2025-08-22 15:30:00"`
	var dt types.DateTimeOnly
	err := dt.UnmarshalParam(param)
	assert.NoError(t, err)
	assert.Equal(t, "2025-08-22 15:30:00", dt.String())
}

func TestGobEncodeDecode(t *testing.T) {
	original := types.NewDateTimeOnly(time.Now().Truncate(time.Second))

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(original)
	assert.NoError(t, err)

	var decoded types.DateTimeOnly
	err = gob.NewDecoder(&buf).Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, original.Time, decoded.Time)
}

func TestValue(t *testing.T) {
	now := types.NewDateTimeOnly(time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC))
	val, err := now.Value()
	assert.NoError(t, err)

	strVal, ok := val.(string)
	assert.True(t, ok)

	parsed, err := time.Parse(time.DateTime, strVal)
	assert.NoError(t, err)
	assert.Equal(t, now.Time, parsed)
}

func TestScan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "from time.Time",
			input:    time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC),
		},
		{
			name:     "from *time.Time",
			input:    ptrTime(time.Date(2025, 8, 22, 12, 0, 0, 0, time.UTC)),
			expected: time.Date(2025, 8, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "from DateTimeOnly",
			input:    types.NewDateTimeOnly(time.Date(2025, 8, 22, 10, 0, 0, 0, time.UTC)),
			expected: time.Date(2025, 8, 22, 10, 0, 0, 0, time.UTC),
		},
		{
			name:     "from *DateTimeOnly",
			input:    ptrDTO(types.NewDateTimeOnly(time.Date(2025, 8, 22, 8, 15, 0, 0, time.UTC))),
			expected: time.Date(2025, 8, 22, 8, 15, 0, 0, time.UTC),
		},
		{
			name:     "from string (standard layout)",
			input:    "2025-08-22 15:30:00",
			expected: time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC),
		},
		{
			name:     "from string (fallback layout)",
			input:    "2025-08-22T15:30:00Z",
			expected: time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC),
		},
		{
			name:     "from []byte",
			input:    []byte("2025-08-22 15:30:00"),
			expected: time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC),
		},
		{
			name:     "from empty string",
			input:    "",
			expected: time.Time{},
		},
		{
			name:     "from nil",
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
			var dt types.DateTimeOnly
			err := dt.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err, "esperaba error pero fue nil")
			} else {
				assert.NoError(t, err, "no esperaba error pero fue: %v", err)
				assert.Equal(t, tt.expected, dt.Time, "tiempo incorrecto")
			}
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrDTO(d types.DateTimeOnly) *types.DateTimeOnly {
	return &d
}

func TestString(t *testing.T) {
	tm := time.Date(2025, 8, 22, 15, 30, 0, 0, time.UTC)
	dt := types.NewDateTimeOnly(tm)
	assert.Equal(t, "2025-08-22 15:30:00", dt.String())
}

func TestDateTimeOnly_ToTimeInLocation(t *testing.T) {
	t.Run("with valid location", func(t *testing.T) {
		loc, err := time.LoadLocation("America/Lima")
		assert.NoError(t, err)

		dt := types.NewDateTimeOnly(time.Date(2025, 8, 22, 15, 45, 30, 0, time.UTC))

		result := dt.ToTimeInLocation(loc)

		expected := time.Date(2025, 8, 22, 15, 45, 30, 0, loc)

		assert.Equal(t, expected.Format(time.RFC3339Nano), result.Format(time.RFC3339Nano))
	})

	t.Run("with nil location (fallback to UTC)", func(t *testing.T) {
		dt := types.NewDateTimeOnly(time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC))

		result := dt.ToTimeInLocation(nil)

		expected := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

		assert.Equal(t, expected.Format(time.RFC3339Nano), result.Format(time.RFC3339Nano))
	})
}
