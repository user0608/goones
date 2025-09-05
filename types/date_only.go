package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type DateOnly struct{ time.Time }

func NewDateOnlyFromString(value string) (DateOnly, error) {
	var jt DateOnly
	if err := jt.UnmarshalJSON([]byte(value)); err != nil {
		return jt, err
	}
	return jt, nil
}

func NewDateOnly(t time.Time) DateOnly {
	str := t.Format(time.DateOnly)
	val, _ := time.Parse(time.DateOnly, str)
	return DateOnly{Time: val}
}

func (jt DateOnly) ToTimeInLocation(loc *time.Location) time.Time {
	if loc == nil {
		slog.Warn("DateOnly: nil location provided, defaulting to UTC")
		loc = time.UTC
	}

	str := jt.Format(time.DateOnly)
	t, err := time.ParseInLocation(time.DateOnly, str, loc)
	if err != nil {
		slog.Warn(
			"DateOnly: failed to convert to time.Time in location",
			"value", str,
			"location", loc.String(),
			"error", err.Error(),
		)
		return time.Time{}
	}

	return t
}

func (jt DateOnly) ToUTCDayRange(loc *time.Location) (time.Time, time.Time) {
	day := jt.ToTimeInLocation(loc)
	start := time.Date(
		day.Year(),
		day.Month(),
		day.Day(), 0, 0, 0, 0,
		day.Location(),
	).UTC()
	end := time.Date(
		day.Year(),
		day.Month(),
		day.Day(), 23, 59, 59,
		int(time.Second-time.Nanosecond),
		day.Location(),
	).UTC()
	return start, end
}

func (jt DateOnly) StartOfDayUTC(loc *time.Location) time.Time {
	day := jt.ToTimeInLocation(loc)
	start := time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		0, 0, 0, 0,
		day.Location(),
	).UTC()
	return start
}

func (jt DateOnly) EndOfDayUTC(loc *time.Location) time.Time {
	day := jt.ToTimeInLocation(loc)
	end := time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		23, 59, 59,
		int(time.Second-time.Nanosecond),
		day.Location(),
	).UTC()
	return end
}

// BuildUTCDayRange builds a UTC time range from a base date (DateOnly)
// and two times (start and end). If the end time is earlier than the start time,
// it is assumed to belong to the following day.
func (jt DateOnly) BuildUTCDayRange(loc *time.Location, start, end JustTime) (time.Time, time.Time) {
	day := jt.ToTimeInLocation(loc)

	left := day.Add(time.Duration(start))
	right := day.Add(time.Duration(end))

	// if the end time is less than or equal to the start time → add one day
	if !right.After(left) {
		right = right.Add(24 * time.Hour)
	}

	return left.UTC(), right.UTC()
}

var _ json.Marshaler = DateOnly{}

func (do DateOnly) MarshalJSON() ([]byte, error) {
	if do.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + do.Format(time.DateOnly) + `"`), nil
}

var _ json.Unmarshaler = (*DateOnly)(nil)

func (od *DateOnly) UnmarshalJSON(data []byte) (err error) {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}
	od.Time, err = time.Parse(time.DateOnly, value)
	return
}

func (do *DateOnly) UnmarshalParam(value string) (err error) {
	return do.UnmarshalJSON([]byte(value))
}

var _ gob.GobEncoder = DateOnly{}

func (do DateOnly) GobEncode() ([]byte, error) {
	return json.Marshal(do)
}

var _ gob.GobDecoder = (*DateOnly)(nil)

func (do *DateOnly) GobDecode(b []byte) error { return json.Unmarshal(b, do) }

var _ driver.Valuer = DateOnly{}

func (do DateOnly) Value() (driver.Value, error) {
	return do.Time.Format(time.DateOnly), nil
}

var _ sql.Scanner = (*DateOnly)(nil)

func (do *DateOnly) Scan(value any) (err error) {
	switch v := value.(type) {
	case nil:
		do.Time = time.Time{}

	case time.Time:
		do.Time = v

	case *time.Time:
		if v != nil {
			do.Time = *v
		} else {
			do.Time = time.Time{}
		}

	case DateOnly:
		*do = v

	case *DateOnly:
		if v != nil {
			*do = *v
		} else {
			do.Time = time.Time{}
		}

	case string, []byte:
		str := cast.ToString(v)
		if str == "" {
			do.Time = time.Time{}
		} else {
			t, err := time.Parse(time.DateOnly, str)
			if err != nil {
				t, err = cast.ToTimeE(v)
				if err != nil {
					return fmt.Errorf("DateOnly: error al parsear `%s`: %w", str, err)
				}
			}
			do.Time = t
		}

	default:
		return fmt.Errorf("DateOnly: tipo no soportado (%T): %v", v, v)
	}

	*do = NewDateOnly(do.Time)
	return nil
}

func (do DateOnly) String() string { return do.Format(time.DateOnly) }
