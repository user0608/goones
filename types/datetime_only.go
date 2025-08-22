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

type DateTimeOnly struct{ time.Time }

func NewDateTimeOnly(t time.Time) DateTimeOnly {
	str := t.Format(time.DateTime)
	val, _ := time.Parse(time.DateTime, str)
	return DateTimeOnly{Time: val}
}

func (jt DateTimeOnly) ToTimeInLocation(loc *time.Location) time.Time {
	if loc == nil {
		slog.Warn("DateTimeOnly: nil location provided, defaulting to UTC")
		loc = time.UTC
	}
	str := jt.Format(time.DateTime)
	t, err := time.ParseInLocation(time.DateTime, str, loc)
	if err != nil {
		slog.Warn(
			"DateTimeOnly: failed to convert to time.Time in location",
			"value", str,
			"location", loc.String(),
			"error", err.Error(),
		)
		return time.Time{}
	}
	return t
}

var _ json.Marshaler = DateTimeOnly{}

func (do DateTimeOnly) MarshalJSON() ([]byte, error) {
	if do.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + do.Format(time.DateTime) + `"`), nil
}

var _ json.Unmarshaler = (*DateTimeOnly)(nil)

func (od *DateTimeOnly) UnmarshalJSON(data []byte) (err error) {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}
	od.Time, err = time.Parse(time.DateTime, value)
	return err
}

func (do *DateTimeOnly) UnmarshalParam(value string) error {
	return do.UnmarshalJSON([]byte(value))
}

var _ gob.GobEncoder = DateTimeOnly{}

func (do DateTimeOnly) GobEncode() ([]byte, error) {
	return json.Marshal(do)
}

var _ gob.GobDecoder = (*DateTimeOnly)(nil)

func (do *DateTimeOnly) GobDecode(b []byte) error {
	return json.Unmarshal(b, do)
}

var _ driver.Valuer = DateTimeOnly{}

func (do DateTimeOnly) Value() (driver.Value, error) {
	return do.Time.Format(time.DateTime), nil
}

var _ sql.Scanner = (*DateTimeOnly)(nil)

func (do *DateTimeOnly) Scan(value any) error {
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

	case DateTimeOnly:
		*do = v

	case *DateTimeOnly:
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
			t, err := time.Parse(time.DateTime, str)
			if err != nil {
				t, err = cast.ToTimeE(v)
				if err != nil {
					return fmt.Errorf("DateTimeOnly: error al parsear `%s`: %w", str, err)
				}
			}
			do.Time = t
		}

	default:
		return fmt.Errorf("DateTimeOnly: tipo no soportado (%T): %v", v, v)
	}

	*do = NewDateTimeOnly(do.Time)
	return nil
}

var _ fmt.Stringer = DateTimeOnly{}

func (do DateTimeOnly) String() string {
	return do.Format(time.DateTime)
}
