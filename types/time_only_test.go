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

func TestJustTime_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{"time.Time", time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC), "12:00:00", false},
		{"*time.Time", ptrTime(time.Date(1, 1, 1, 23, 59, 59, 0, time.UTC)), "23:59:59", false},
		{"JustTime", types.NewJustTime(time.Date(1, 1, 1, 7, 15, 0, 0, time.UTC)), "07:15:00", false},
		{"*JustTime", ptrJT(types.NewJustTime(time.Date(1, 1, 1, 1, 2, 3, 0, time.UTC))), "01:02:03", false},
		{"string", "18:45:59", "18:45:59", false},
		{"[]byte", []byte("22:22:22"), "22:22:22", false},
		{"string con nanos", "12:00:00.000000123", "12:00:00.000000123", false},
		{"nil", nil, "00:00:00", false},
		{"invalid string", "bad-time", "", true},
		{"unsupported type", 123, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jt types.JustTime
			err := jt.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, jt.String())
			}
		})
	}
}

func TestJustTime_String(t *testing.T) {
	jt := types.NewJustTime(time.Date(1, 1, 1, 8, 8, 8, 0, time.UTC))
	assert.Equal(t, "08:08:08", jt.String())

	jtNano := types.NewJustTime(time.Date(1, 1, 1, 1, 1, 1, 999, time.UTC))
	assert.Contains(t, jtNano.String(), "01:01:01.")
}

func TestJustTime_MarshalJSON(t *testing.T) {
	jt := types.NewJustTime(time.Date(1, 1, 1, 14, 0, 0, 0, time.UTC))
	data, err := json.Marshal(jt)
	assert.NoError(t, err)
	assert.Equal(t, `"14:00:00"`, string(data))

	var zero types.JustTime
	zeroData, err := json.Marshal(zero)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(zeroData))
}

func TestJustTime_UnmarshalJSON(t *testing.T) {
	var jt types.JustTime
	err := json.Unmarshal([]byte(`"05:45:30"`), &jt)
	assert.NoError(t, err)
	assert.Equal(t, "05:45:30", jt.String())

	err = json.Unmarshal([]byte(`"null"`), &jt)
	assert.NoError(t, err)
	assert.Equal(t, "00:00:00", jt.String())

	err = json.Unmarshal([]byte(`"bad"`), &jt)
	assert.Error(t, err)
}

func TestJustTime_UnmarshalParam(t *testing.T) {
	var jt types.JustTime
	err := jt.UnmarshalParam("06:00:00")
	assert.NoError(t, err)
	assert.Equal(t, "06:00:00", jt.String())

	err = jt.UnmarshalParam("invalid")
	assert.Error(t, err)
}

func TestJustTime_GobEncodeDecode(t *testing.T) {
	original := types.NewJustTime(time.Date(1, 1, 1, 20, 0, 0, 0, time.UTC))

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(original)
	assert.NoError(t, err)

	var decoded types.JustTime
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}

func ptrJT(t types.JustTime) *types.JustTime {
	return &t
}

func TestJustTime_ToTime(t *testing.T) {
	jt := types.NewJustTime(time.Date(1, 1, 1, 14, 30, 0, 0, time.UTC))
	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 14, 30, 0, 0, time.UTC)

	result := jt.ToTime()

	assert.Equal(t, expected.Year(), result.Year())
	assert.Equal(t, expected.Month(), result.Month())
	assert.Equal(t, expected.Day(), result.Day())
	assert.Equal(t, expected.Hour(), result.Hour())
	assert.Equal(t, expected.Minute(), result.Minute())
	assert.Equal(t, expected.Second(), result.Second())
}

func TestJustTime_ToTimeInLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/Lima")
	assert.NoError(t, err)

	jt := types.NewJustTime(time.Date(1, 1, 1, 9, 15, 0, 0, time.UTC))
	now := time.Now().In(loc)
	expected := time.Date(now.Year(), now.Month(), now.Day(), 9, 15, 0, 0, loc)

	result := jt.ToTimeInLocation(loc)

	assert.Equal(t, expected.Year(), result.Year())
	assert.Equal(t, expected.Month(), result.Month())
	assert.Equal(t, expected.Day(), result.Day())
	assert.Equal(t, expected.Hour(), result.Hour())
	assert.Equal(t, expected.Minute(), result.Minute())
	assert.Equal(t, expected.Second(), result.Second())
	assert.Equal(t, loc.String(), result.Location().String())

	assert.Equal(t, expected.Format(time.RFC3339), result.Format(time.RFC3339))
}
