package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type JustTime time.Duration

func NewJustTimeFromString(value string) (JustTime, error) {
	var jt JustTime
	if err := jt.UnmarshalJSON([]byte(value)); err != nil {
		return jt, err
	}
	return jt, nil
}

func NewJustTime(t time.Time) JustTime {
	return JustTime(
		time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute +
			time.Duration(t.Second())*time.Second +
			time.Duration(t.Nanosecond())*time.Nanosecond,
	)
}

func (jt JustTime) String() string {
	return jt.Format()
}

func (jt JustTime) ToTime() time.Time {
	now := time.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return base.Add(time.Duration(jt))
}

func (jt JustTime) ToTimeInLocation(loc *time.Location) time.Time {
	if loc == nil {
		slog.Warn("JustTime: nil location provided, defaulting to UTC")
		loc = time.UTC
	}
	now := time.Now().In(loc)
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return base.Add(time.Duration(jt))
}

func (jt JustTime) Format() string {
	if jt.nanoseconds() > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%09d", jt.hours(), jt.minutes(), jt.seconds(), jt.nanoseconds())
	}
	return fmt.Sprintf("%02d:%02d:%02d", jt.hours(), jt.minutes(), jt.seconds())
}

func (jt JustTime) MarshalJSON() ([]byte, error) {
	if jt == 0 {
		return []byte("null"), nil
	}
	return []byte(`"` + jt.Format() + `"`), nil
}

func (jt *JustTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		*jt = 0
		return nil
	}
	return jt.UnmarshalParam(value)
}

func (jt *JustTime) UnmarshalParam(value string) error {
	var h, m, s, n int
	var parsed bool

	switch {
	case strings.Count(value, ":") == 2 && strings.Contains(value, "."):
		_, err := fmt.Sscanf(value, "%d:%d:%d.%d", &h, &m, &s, &n)
		if err == nil {
			parsed = true
		}
	case strings.Count(value, ":") == 2:
		_, err := fmt.Sscanf(value, "%d:%d:%d", &h, &m, &s)
		if err == nil {
			parsed = true
		}
	case strings.Count(value, ":") == 1:
		_, err := fmt.Sscanf(value, "%d:%d", &h, &m)
		if err == nil {
			s = 0
			parsed = true
		}
	case strings.Count(value, ":") == 0:
		_, err := fmt.Sscanf(value, "%d", &h)
		if err == nil {
			m, s = 0, 0
			parsed = true
		}
	}

	if !parsed {
		return fmt.Errorf("JustTime: formato inválido `%s`", value)
	}

	if h < 0 || h > 23 || m < 0 || m > 59 || s < 0 || s > 59 || n < 0 || n > 999999999 {
		return fmt.Errorf("JustTime: fuera de rango `%s` (h=%d m=%d s=%d n=%d)", value, h, m, s, n)
	}

	*jt = JustTime(
		time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(s)*time.Second +
			time.Duration(n)*time.Nanosecond,
	)
	return nil
}

func (jt JustTime) Value() (driver.Value, error) {
	return jt.Format(), nil
}

func (jt *JustTime) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*jt = 0
	case time.Time:
		*jt = NewJustTime(v)
	case *time.Time:
		if v != nil {
			*jt = NewJustTime(*v)
		} else {
			*jt = 0
		}
	case JustTime:
		*jt = v
	case *JustTime:
		if v != nil {
			*jt = *v
		} else {
			*jt = 0
		}
	case string, []byte:
		str := cast.ToString(v)
		if str == "" {
			*jt = 0
		} else {
			err := jt.UnmarshalParam(str)
			if err != nil {
				return fmt.Errorf("JustTime: error al parsear `%s`: %w", str, err)
			}
		}
	default:
		return fmt.Errorf("JustTime: tipo no soportado (%T): %v", v, v)
	}
	return nil
}

func (jt JustTime) GobEncode() ([]byte, error) {
	return json.Marshal(jt)
}

func (jt *JustTime) GobDecode(b []byte) error {
	return json.Unmarshal(b, jt)
}

// helpers
func (jt JustTime) hours() int {
	return int(time.Duration(jt).Truncate(time.Hour).Hours())
}

func (jt JustTime) minutes() int {
	return int((time.Duration(jt) % time.Hour).Truncate(time.Minute).Minutes())
}

func (jt JustTime) seconds() int {
	return int((time.Duration(jt) % time.Minute).Truncate(time.Second).Seconds())
}

func (jt JustTime) nanoseconds() int {
	return int((time.Duration(jt) % time.Second).Nanoseconds())
}
