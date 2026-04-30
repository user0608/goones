package kcheck

import (
	"strings"
	"testing"
	"time"
)

func TestValidatorsRequired(t *testing.T) {
	if err := required(Field{Path: "Name", Value: "Kevin"}); err != nil {
		t.Fatalf("expected valid required, got %v", err)
	}

	if err := required(Field{Path: "Name", Value: ""}); err == nil {
		t.Fatal("expected required error")
	}

	if err := required(Field{Path: "Name", IsNil: true}); err == nil {
		t.Fatal("expected nil required error")
	}
}

func TestValidatorsLengthMinMaxString(t *testing.T) {
	if err := length(Field{Path: "Code", Value: "José", Param: "4"}); err != nil {
		t.Fatalf("expected valid length, got %v", err)
	}

	if err := length(Field{Path: "Code", Value: "José", Param: "3"}); err == nil {
		t.Fatal("expected length error")
	}

	if err := min(Field{Path: "Name", Value: "José", Param: "4"}); err != nil {
		t.Fatalf("expected valid min string, got %v", err)
	}

	if err := min(Field{Path: "Name", Value: "José", Param: "5"}); err == nil {
		t.Fatal("expected min string error")
	}

	if err := max(Field{Path: "Name", Value: "José", Param: "4"}); err != nil {
		t.Fatalf("expected valid max string, got %v", err)
	}

	if err := max(Field{Path: "Name", Value: "José", Param: "3"}); err == nil {
		t.Fatal("expected max string error")
	}
}

func TestValidatorsMinMaxNumbers(t *testing.T) {
	if err := min(Field{Path: "Age", Value: 18, Param: "18"}); err != nil {
		t.Fatalf("expected valid min number, got %v", err)
	}

	if err := min(Field{Path: "Age", Value: 17, Param: "18"}); err == nil {
		t.Fatal("expected min number error")
	}

	if err := max(Field{Path: "Age", Value: 120, Param: "120"}); err != nil {
		t.Fatalf("expected valid max number, got %v", err)
	}

	if err := max(Field{Path: "Age", Value: 121, Param: "120"}); err == nil {
		t.Fatal("expected max number error")
	}
}

func TestValidatorsComparators(t *testing.T) {
	if err := greaterThan(Field{Path: "A", Value: 11, Param: "10", Tag: "gt"}); err != nil {
		t.Fatalf("expected gt valid, got %v", err)
	}

	if err := greaterThan(Field{Path: "A", Value: 10, Param: "10", Tag: "gt"}); err == nil {
		t.Fatal("expected gt error")
	}

	if err := greaterThanOrEqual(Field{Path: "B", Value: 10, Param: "10", Tag: "gte"}); err != nil {
		t.Fatalf("expected gte valid, got %v", err)
	}

	if err := greaterThanOrEqual(Field{Path: "B", Value: 9, Param: "10", Tag: "gte"}); err == nil {
		t.Fatal("expected gte error")
	}

	if err := lessThan(Field{Path: "C", Value: 9, Param: "10", Tag: "lt"}); err != nil {
		t.Fatalf("expected lt valid, got %v", err)
	}

	if err := lessThan(Field{Path: "C", Value: 10, Param: "10", Tag: "lt"}); err == nil {
		t.Fatal("expected lt error")
	}

	if err := lessThanOrEqual(Field{Path: "D", Value: 10, Param: "10", Tag: "lte"}); err != nil {
		t.Fatalf("expected lte valid, got %v", err)
	}

	if err := lessThanOrEqual(Field{Path: "D", Value: 11, Param: "10", Tag: "lte"}); err == nil {
		t.Fatal("expected lte error")
	}
}

func TestValidatorsEmailUUIDURL(t *testing.T) {
	if err := email(Field{Path: "Email", Value: "test@example.com"}); err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}

	if err := email(Field{Path: "Email", Value: "bad-email"}); err == nil {
		t.Fatal("expected email error")
	}

	if err := uuidV4(Field{Path: "ID", Value: "550e8400-e29b-41d4-a716-446655440000"}); err != nil {
		t.Fatalf("expected valid uuid, got %v", err)
	}

	if err := uuidV4(Field{Path: "ID", Value: "550e8400-e29b-11d4-a716-446655440000"}); err == nil {
		t.Fatal("expected uuid error")
	}

	if err := urlValue(Field{Path: "URL", Value: "https://example.com/path?q=1"}); err != nil {
		t.Fatalf("expected valid url, got %v", err)
	}

	if err := urlValue(Field{Path: "URL", Value: "example.com"}); err == nil {
		t.Fatal("expected url error")
	}
}

func TestValidatorsIP(t *testing.T) {
	if err := ip(Field{Path: "IP", Value: "192.168.1.1"}); err != nil {
		t.Fatalf("expected valid ip, got %v", err)
	}

	if err := ip(Field{Path: "IP", Value: "2001:db8::1"}); err != nil {
		t.Fatalf("expected valid ipv6 as ip, got %v", err)
	}

	if err := ip(Field{Path: "IP", Value: "bad"}); err == nil {
		t.Fatal("expected ip error")
	}

	if err := ipv4(Field{Path: "IPv4", Value: "192.168.1.1"}); err != nil {
		t.Fatalf("expected valid ipv4, got %v", err)
	}

	if err := ipv4(Field{Path: "IPv4", Value: "2001:db8::1"}); err == nil {
		t.Fatal("expected ipv4 error")
	}

	if err := ipv6(Field{Path: "IPv6", Value: "2001:db8::1"}); err != nil {
		t.Fatalf("expected valid ipv6, got %v", err)
	}

	if err := ipv6(Field{Path: "IPv6", Value: "192.168.1.1"}); err == nil {
		t.Fatal("expected ipv6 error")
	}
}

func TestValidatorsStringFormats(t *testing.T) {
	if err := alpha(Field{Path: "Name", Value: "JoséÑandu"}); err != nil {
		t.Fatalf("expected alpha valid, got %v", err)
	}

	if err := alpha(Field{Path: "Name", Value: "José123"}); err == nil {
		t.Fatal("expected alpha error")
	}

	if err := alphanum(Field{Path: "Code", Value: "José123"}); err != nil {
		t.Fatalf("expected alphanum valid, got %v", err)
	}

	if err := alphanum(Field{Path: "Code", Value: "José-123"}); err == nil {
		t.Fatal("expected alphanum error")
	}

	if err := numericString(Field{Path: "Num", Value: "12345"}); err != nil {
		t.Fatalf("expected num valid, got %v", err)
	}

	if err := numericString(Field{Path: "Num", Value: "123a"}); err == nil {
		t.Fatal("expected num error")
	}

	if err := decimalString(Field{Path: "Decimal", Value: "123.45"}); err != nil {
		t.Fatalf("expected decimal valid, got %v", err)
	}

	if err := decimalString(Field{Path: "Decimal", Value: "123,45"}); err == nil {
		t.Fatal("expected decimal error")
	}

	if err := lower(Field{Path: "Lower", Value: "hello"}); err != nil {
		t.Fatalf("expected lower valid, got %v", err)
	}

	if err := lower(Field{Path: "Lower", Value: "Hello"}); err == nil {
		t.Fatal("expected lower error")
	}

	if err := upper(Field{Path: "Upper", Value: "HELLO"}); err != nil {
		t.Fatalf("expected upper valid, got %v", err)
	}

	if err := upper(Field{Path: "Upper", Value: "Hello"}); err == nil {
		t.Fatal("expected upper error")
	}
}

func TestValidatorsOneOfPrefixSuffixContains(t *testing.T) {
	if err := oneOf(Field{Path: "Status", Value: "active", Param: "active,inactive,pending"}); err != nil {
		t.Fatalf("expected oneof valid, got %v", err)
	}

	if err := oneOf(Field{Path: "Status", Value: "deleted", Param: "active,inactive,pending"}); err == nil {
		t.Fatal("expected oneof error")
	}

	if err := prefix(Field{Path: "Code", Value: "USR-123", Param: "USR-"}); err != nil {
		t.Fatalf("expected prefix valid, got %v", err)
	}

	if err := prefix(Field{Path: "Code", Value: "ADM-123", Param: "USR-"}); err == nil {
		t.Fatal("expected prefix error")
	}

	if err := suffix(Field{Path: "File", Value: "report.pdf", Param: ".pdf"}); err != nil {
		t.Fatalf("expected suffix valid, got %v", err)
	}

	if err := suffix(Field{Path: "File", Value: "report.txt", Param: ".pdf"}); err == nil {
		t.Fatal("expected suffix error")
	}

	if err := contains(Field{Path: "Text", Value: "hello world", Param: "world"}); err != nil {
		t.Fatalf("expected contains valid, got %v", err)
	}

	if err := contains(Field{Path: "Text", Value: "hello world", Param: "admin"}); err == nil {
		t.Fatal("expected contains error")
	}
}

func TestValidatorsDateTimeUTC(t *testing.T) {
	if err := dateValue(Field{Path: "Date", Value: "2026-04-30"}); err != nil {
		t.Fatalf("expected date valid, got %v", err)
	}

	if err := dateValue(Field{Path: "Date", Value: "30-04-2026"}); err == nil {
		t.Fatal("expected date error")
	}

	if err := timeValue(Field{Path: "Time", Value: "15:04:05"}); err != nil {
		t.Fatalf("expected time valid, got %v", err)
	}

	if err := timeValue(Field{Path: "Time", Value: "15:04"}); err == nil {
		t.Fatal("expected time error")
	}

	if err := dateTimeValue(Field{Path: "DateTime", Value: "2026-04-30 15:04:05"}); err != nil {
		t.Fatalf("expected datetime valid, got %v", err)
	}

	if err := dateTimeValue(Field{Path: "DateTime", Value: "2026-04-30T15:04:05Z"}); err == nil {
		t.Fatal("expected datetime error")
	}

	if err := utcValue(Field{Path: "UTC", Value: "2026-04-30T15:04:05Z"}); err != nil {
		t.Fatalf("expected utc string valid, got %v", err)
	}

	if err := utcValue(Field{Path: "UTC", Value: "2026-04-30T15:04:05-05:00"}); err == nil {
		t.Fatal("expected utc offset error")
	}

	if err := utcValue(Field{Path: "UTC", Value: "2026-04-30 15:04:05"}); err == nil {
		t.Fatal("expected utc format error")
	}

	if err := utcValue(Field{Path: "UTC", Value: time.Date(2026, 4, 30, 15, 4, 5, 0, time.UTC)}); err != nil {
		t.Fatalf("expected utc time.Time valid, got %v", err)
	}

	if err := utcValue(Field{Path: "UTC", Value: time.Date(2026, 4, 30, 15, 4, 5, 0, time.FixedZone("PET", -5*60*60))}); err == nil {
		t.Fatal("expected non utc time.Time error")
	}
}

func TestValidatorsUnsupportedTypes(t *testing.T) {
	if err := length(Field{Path: "Age", Value: 10, Param: "2"}); err == nil {
		t.Fatal("expected length unsupported type error")
	}

	if err := email(Field{Path: "Email", Value: 10}); err == nil {
		t.Fatal("expected email unsupported type error")
	}

	if err := uuidV4(Field{Path: "ID", Value: 10}); err == nil {
		t.Fatal("expected uuid unsupported type error")
	}

	if err := urlValue(Field{Path: "URL", Value: 10}); err == nil {
		t.Fatal("expected url unsupported type error")
	}

	if err := ip(Field{Path: "IP", Value: 10}); err == nil {
		t.Fatal("expected ip unsupported type error")
	}

	if err := ipv4(Field{Path: "IPv4", Value: 10}); err == nil {
		t.Fatal("expected ipv4 unsupported type error")
	}

	if err := ipv6(Field{Path: "IPv6", Value: 10}); err == nil {
		t.Fatal("expected ipv6 unsupported type error")
	}

	if err := alpha(Field{Path: "Alpha", Value: 10}); err == nil {
		t.Fatal("expected alpha unsupported type error")
	}

	if err := alphanum(Field{Path: "AlphaNum", Value: 10}); err == nil {
		t.Fatal("expected alphanum unsupported type error")
	}

	if err := numericString(Field{Path: "Num", Value: 10}); err == nil {
		t.Fatal("expected num unsupported type error")
	}

	if err := decimalString(Field{Path: "Decimal", Value: 10}); err == nil {
		t.Fatal("expected decimal unsupported type error")
	}

	if err := lower(Field{Path: "Lower", Value: 10}); err == nil {
		t.Fatal("expected lower unsupported type error")
	}

	if err := upper(Field{Path: "Upper", Value: 10}); err == nil {
		t.Fatal("expected upper unsupported type error")
	}

	if err := prefix(Field{Path: "Prefix", Value: 10, Param: "x"}); err == nil {
		t.Fatal("expected prefix unsupported type error")
	}

	if err := suffix(Field{Path: "Suffix", Value: 10, Param: "x"}); err == nil {
		t.Fatal("expected suffix unsupported type error")
	}

	if err := contains(Field{Path: "Contains", Value: 10, Param: "x"}); err == nil {
		t.Fatal("expected contains unsupported type error")
	}

	if err := dateValue(Field{Path: "Date", Value: 10}); err == nil {
		t.Fatal("expected date unsupported type error")
	}

	if err := timeValue(Field{Path: "Time", Value: 10}); err == nil {
		t.Fatal("expected time unsupported type error")
	}

	if err := dateTimeValue(Field{Path: "DateTime", Value: 10}); err == nil {
		t.Fatal("expected datetime unsupported type error")
	}

	if err := utcValue(Field{Path: "UTC", Value: 10}); err == nil {
		t.Fatal("expected utc unsupported type error")
	}
}

func TestValidatorStructOK(t *testing.T) {
	type Address struct {
		City string `chk:"required min=3"`
	}

	type User struct {
		Name      string    `chk:"required min=2 max=20"`
		Email     string    `chk:"required email"`
		Age       int       `chk:"gte=18 lte=120"`
		Code      *string   `chk:"required upper len=6"`
		Status    string    `chk:"required oneof=active,inactive,pending"`
		Website   *string   `chk:"url"`
		CreatedAt time.Time `chk:"required utc"`
		Address   Address
	}

	user := User{
		Name:      "Kevin",
		Email:     "kevin@example.com",
		Age:       30,
		Code:      strPtr("ABC123"),
		Status:    "active",
		Website:   strPtr("https://example.com"),
		CreatedAt: time.Date(2026, 4, 30, 15, 4, 5, 0, time.UTC),
		Address:   Address{City: "Lima"},
	}

	if err := Valid(user); err != nil {
		t.Fatalf("expected valid user, got %v", err)
	}
}

func TestValidatorStructInvalid(t *testing.T) {
	type User struct {
		Name  string `chk:"required"`
		Email string `chk:"email"`
		Age   int    `chk:"gte=18"`
	}

	err := Valid(User{
		Name:  "",
		Email: "bad",
		Age:   10,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := err.Error()
	for _, want := range []string{"Name", "Email", "Age"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("expected error to contain %s, got %v", want, err)
		}
	}
}

func TestValidatorSkip(t *testing.T) {
	type User struct {
		Name  string `chk:"required"`
		Email string `chk:"email"`
	}

	err := Valid(User{
		Name:  "Kevin",
		Email: "bad",
	}, "Email")
	if err != nil {
		t.Fatalf("expected Email skipped, got %v", err)
	}
}

func TestValidatorSelect(t *testing.T) {
	type User struct {
		Name  string `chk:"required"`
		Email string `chk:"email"`
		Age   int    `chk:"gte=18"`
	}

	err := ValidSelect(User{
		Name:  "Kevin",
		Email: "bad",
		Age:   10,
	}, "Name")
	if err != nil {
		t.Fatalf("expected only Name validated, got %v", err)
	}
}

func TestValidatorNestedSelectAndSkip(t *testing.T) {
	type Address struct {
		City string `chk:"required min=3"`
	}

	type User struct {
		Name    string `chk:"required"`
		Address Address
	}

	user := User{
		Name:    "",
		Address: Address{City: "Li"},
	}

	err := ValidSelect(user, "Address.City")
	if err == nil {
		t.Fatal("expected Address.City error")
	}

	if strings.Contains(err.Error(), "Name") {
		t.Fatalf("expected Name ignored, got %v", err)
	}

	err = Valid(user, "Address.City")
	if err == nil {
		t.Fatal("expected Name error")
	}

	if strings.Contains(err.Error(), "Address.City") {
		t.Fatalf("expected Address.City skipped, got %v", err)
	}
}

func TestValidatorInvalidInput(t *testing.T) {
	if err := Valid(nil); err == nil {
		t.Fatal("expected nil input error")
	}

	var user *struct {
		Name string `chk:"required"`
	}

	if err := Valid(user); err == nil {
		t.Fatal("expected nil pointer input error")
	}

	if err := Valid("bad"); err == nil {
		t.Fatal("expected non struct input error")
	}
}

func TestValidatorUnknownValidator(t *testing.T) {
	type User struct {
		Name string `chk:"unknown"`
	}

	err := Valid(User{Name: "Kevin"})
	if err == nil {
		t.Fatal("expected unknown validator error")
	}

	if !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected unknown in error, got %v", err)
	}
}

func TestValidatorCustomRegister(t *testing.T) {
	v := New()

	v.Register("startsx", func(f Field) error {
		s, ok := f.Value.(string)
		if !ok {
			return nil
		}
		if !strings.HasPrefix(s, "x") {
			return ErrInvalidInput
		}
		return nil
	})

	type DTO struct {
		Code string `chk:"startsx"`
	}

	if err := v.Struct(DTO{Code: "x123"}); err != nil {
		t.Fatalf("expected custom validator valid, got %v", err)
	}

	if err := v.Struct(DTO{Code: "a123"}); err == nil {
		t.Fatal("expected custom validator error")
	}
}
