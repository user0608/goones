package types_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user0608/goones/types"
)

func TestNewDateOnlyFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{`"2025-09-05"`, "2025-09-05", false},
		{`2025-09-05`, "2025-09-05", false},
		{`"2000-01-01"`, "2000-01-01", false},
		{`2000-01-01`, "2000-01-01", false},
		{`"1999-12-31"`, "1999-12-31", false},
		{`1999-12-31`, "1999-12-31", false},

		// errores
		{`"2025-13-01"`, "", true},
		{`2025-13-01`, "", true},
		{`"2025-00-01"`, "", true},
		{`"2025-01-32"`, "", true},
		{`"invalid"`, "", true},
		{`invalid`, "", true},
		{`""`, time.Time{}.Format(time.DateOnly), false}, //zero
		{``, time.Time{}.Format(time.DateOnly), false},   //zero
	}

	for _, tt := range tests {
		dt, err := types.NewDateOnlyFromString(tt.input)

		if tt.wantErr {
			if err == nil {
				t.Errorf("input %s: expected error, got none (dt.String()=%s)", tt.input, dt.String())
			}
			continue
		}

		if err != nil {
			t.Errorf("input %s: unexpected error: %v (dt.String()=%s)", tt.input, err, dt.String())
			continue
		}

		if got := dt.String(); got != tt.expected {
			t.Errorf("input %s: expected %s, got %s", tt.input, tt.expected, got)
		}
	}
}

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

func TestDateOnly_ToUTCDayRange(t *testing.T) {
	loc, err := time.LoadLocation("America/Lima")
	assert.NoError(t, err)

	input := time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC)
	dateOnly := types.NewDateOnly(input)

	start, end := dateOnly.ToUTCDayRange(loc)

	expectedStart := time.Date(2025, 8, 22, 0, 0, 0, 0, loc).UTC()
	expectedEnd := time.Date(2025, 8, 22, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()

	assert.Equal(t, expectedStart.Format(time.RFC3339Nano), start.Format(time.RFC3339Nano))
	assert.Equal(t, expectedEnd.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano))
}

func TestDateOnly_StartAndEndOfDayUTC(t *testing.T) {
	loc, err := time.LoadLocation("America/Lima")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	tests := []struct {
		name          string
		date          string
		expectedStart string
		expectedEnd   string
	}{
		{
			name:          "start and end of day in Lima timezone",
			date:          "2025-09-05",
			expectedStart: "2025-09-05T05:00:00Z",           // 00:00:00 -5h
			expectedEnd:   "2025-09-06T04:59:59.999999999Z", // 23:59:59.999999999 -5h
		},
		{
			name:          "start and end of day - UTC",
			date:          "2000-01-01",
			expectedStart: "2000-01-01T05:00:00Z",           // 00:00:00 -5h
			expectedEnd:   "2000-01-02T04:59:59.999999999Z", // 23:59:59.999999999 -5h
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := types.NewDateOnlyFromString(fmt.Sprintf(`"%s"`, tt.date))
			if err != nil {
				t.Fatalf("error parsing date: %v", err)
			}

			start := dt.StartOfDayUTC(loc)
			end := dt.EndOfDayUTC(loc)

			if got := start.Format(time.RFC3339Nano); got != tt.expectedStart {
				t.Errorf("StartOfDayUTC: expected %s, got %s", tt.expectedStart, got)
			}
			if got := end.Format(time.RFC3339Nano); got != tt.expectedEnd {
				t.Errorf("EndOfDayUTC: expected %s, got %s", tt.expectedEnd, got)
			}
		})
	}
}

func TestDateOnly_StartAndEndOfDayUTC_UTCZone(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name          string
		date          string
		expectedStart string
		expectedEnd   string
	}{
		{
			name:          "UTC start and end of day",
			date:          "2025-09-05",
			expectedStart: "2025-09-05T00:00:00Z",
			expectedEnd:   "2025-09-05T23:59:59.999999999Z",
		},
		{
			name:          "another UTC case",
			date:          "2000-01-01",
			expectedStart: "2000-01-01T00:00:00Z",
			expectedEnd:   "2000-01-01T23:59:59.999999999Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := types.NewDateOnlyFromString(fmt.Sprintf(`"%s"`, tt.date))
			if err != nil {
				t.Fatalf("error parsing date: %v", err)
			}

			start := dt.StartOfDayUTC(loc)
			end := dt.EndOfDayUTC(loc)

			if got := start.Format(time.RFC3339Nano); got != tt.expectedStart {
				t.Errorf("StartOfDayUTC: expected %s, got %s", tt.expectedStart, got)
			}
			if got := end.Format(time.RFC3339Nano); got != tt.expectedEnd {
				t.Errorf("EndOfDayUTC: expected %s, got %s", tt.expectedEnd, got)
			}
		})
	}
}

func TestDateOnly_BuildUTCDayRange(t *testing.T) {
	loc, err := time.LoadLocation("America/Lima")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	type args struct {
		date  string
		start string
		end   string
	}

	tests := []struct {
		name          string
		args          args
		expectedStart string
		expectedEnd   string
	}{
		{
			name: "rango normal mismo día",
			args: args{
				date:  "2025-09-05",
				start: "08:00:00",
				end:   "17:00:00",
			},
			expectedStart: "2025-09-05T13:00:00Z", // UTC = local + 5h
			expectedEnd:   "2025-09-05T22:00:00Z",
		},
		{
			name: "end menor que start, cruza medianoche",
			args: args{
				date:  "2025-09-05",
				start: "22:00:00",
				end:   "02:00:00",
			},
			expectedStart: "2025-09-06T03:00:00Z", // 22:00 -5h
			expectedEnd:   "2025-09-06T07:00:00Z", // +1 día
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := types.NewDateOnlyFromString(fmt.Sprintf(`"%s"`, tt.args.date))
			if err != nil {
				t.Fatalf("error parsing date: %v", err)
			}

			var start, end types.JustTime
			if err := start.UnmarshalParam(tt.args.start); err != nil {
				t.Fatalf("error parsing start: %v", err)
			}
			if err := end.UnmarshalParam(tt.args.end); err != nil {
				t.Fatalf("error parsing end: %v", err)
			}

			left, right := dt.BuildUTCDayRange(loc, start, end)

			if got := left.Format(time.RFC3339); got != tt.expectedStart {
				t.Errorf("start: expected %s, got %s", tt.expectedStart, got)
			}
			if got := right.Format(time.RFC3339); got != tt.expectedEnd {
				t.Errorf("end: expected %s, got %s", tt.expectedEnd, got)
			}
		})
	}
}

func TestDateOnly_BuildUTCDayRange_UTC(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name          string
		date          string
		start         string
		end           string
		expectedStart string
		expectedEnd   string
	}{
		{
			name:          "normal range (same day)",
			date:          "2025-09-05",
			start:         "06:00:00",
			end:           "10:00:00",
			expectedStart: "2025-09-05T06:00:00Z",
			expectedEnd:   "2025-09-05T10:00:00Z",
		},
		{
			name:          "crossing midnight (end before start)",
			date:          "2025-09-05",
			start:         "22:00:00",
			end:           "02:00:00",
			expectedStart: "2025-09-05T22:00:00Z",
			expectedEnd:   "2025-09-06T02:00:00Z", // next day
		},
		{
			name:          "exactly same time (treated as full day)",
			date:          "2025-09-05",
			start:         "00:00:00",
			end:           "00:00:00",
			expectedStart: "2025-09-05T00:00:00Z",
			expectedEnd:   "2025-09-06T00:00:00Z", // +24h
		},
		{
			name:          "same start and end time (treated as full day)",
			date:          "2025-09-05",
			start:         "14:15:16",
			end:           "14:15:16",
			expectedStart: "2025-09-05T14:15:16Z",
			expectedEnd:   "2025-09-06T14:15:16Z", // +24h
		},
		{
			name:          "with seconds and nanos",
			date:          "2025-09-05",
			start:         "01:02:03.123456789",
			end:           "04:05:06.987654321",
			expectedStart: "2025-09-05T01:02:03Z",
			expectedEnd:   "2025-09-05T04:05:06Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := types.NewDateOnlyFromString(fmt.Sprintf(`"%s"`, tt.date))
			if err != nil {
				t.Fatalf("error parsing date: %v", err)
			}

			var start, end types.JustTime
			if err := start.UnmarshalParam(tt.start); err != nil {
				t.Fatalf("error parsing start time: %v", err)
			}
			if err := end.UnmarshalParam(tt.end); err != nil {
				t.Fatalf("error parsing end time: %v", err)
			}

			left, right := dt.BuildUTCDayRange(loc, start, end)

			gotStart := left.Format(time.RFC3339)
			gotEnd := right.Format(time.RFC3339)

			if gotStart != tt.expectedStart {
				t.Errorf("start mismatch: expected %s, got %s", tt.expectedStart, gotStart)
			}
			if gotEnd != tt.expectedEnd {
				t.Errorf("end mismatch: expected %s, got %s", tt.expectedEnd, gotEnd)
			}
		})
	}
}
