package kcheck

import (
	"strings"
	"testing"
	"time"
)

func strPtr(v string) *string {
	return &v
}

type testAddress struct {
	City string `chk:"required min=3"`
}

type testUser struct {
	Name      string      `chk:"required min=2 max=20"`
	Email     string      `chk:"required email"`
	Age       int         `chk:"gte=18 lte=120"`
	Code      *string     `chk:"required upper len=6"`
	Status    string      `chk:"required oneof=active,inactive,pending"`
	Website   *string     `chk:"url"`
	CreatedAt time.Time   `chk:"required datetime"`
	Address   testAddress `chk:"-"`
}

func TestValidOK(t *testing.T) {
	user := testUser{
		Name:      "Kevin",
		Email:     "kevin@example.com",
		Age:       30,
		Code:      strPtr("ABC123"),
		Status:    "active",
		Website:   strPtr("https://example.com"),
		CreatedAt: time.Now(),
		Address:   testAddress{City: "Lima"},
	}

	if err := Valid(user); err != nil {
		t.Fatalf("expected valid user, got error: %v", err)
	}
}

func TestValidWithPointerOK(t *testing.T) {
	user := &testUser{
		Name:      "Kevin",
		Email:     "kevin@example.com",
		Age:       30,
		Code:      strPtr("ABC123"),
		Status:    "active",
		Website:   strPtr("https://example.com"),
		CreatedAt: time.Now(),
		Address:   testAddress{City: "Lima"},
	}

	if err := Valid(user); err != nil {
		t.Fatalf("expected valid user pointer, got error: %v", err)
	}
}

func TestValidNilInput(t *testing.T) {
	if err := Valid(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestValidNilPointerInput(t *testing.T) {
	var user *testUser

	if err := Valid(user); err == nil {
		t.Fatal("expected error for nil pointer input")
	}
}

func TestValidInvalidNonStructInput(t *testing.T) {
	if err := Valid("invalid"); err == nil {
		t.Fatal("expected error for non-struct input")
	}
}

func TestRequiredString(t *testing.T) {
	type dto struct {
		Name string `chk:"required"`
	}

	err := Valid(dto{Name: ""})
	if err == nil {
		t.Fatal("expected required error")
	}

	if !strings.Contains(err.Error(), "Name") {
		t.Fatalf("expected error to contain field name, got: %v", err)
	}
}

func TestRequiredPointerString(t *testing.T) {
	type dto struct {
		Name *string `chk:"required"`
	}

	if err := Valid(dto{Name: strPtr("Kevin")}); err != nil {
		t.Fatalf("expected valid pointer string, got: %v", err)
	}

	if err := Valid(dto{Name: nil}); err == nil {
		t.Fatal("expected error for nil pointer string")
	}
}

func TestSkipField(t *testing.T) {
	user := testUser{
		Name:      "Kevin",
		Email:     "bad-email",
		Age:       30,
		Code:      strPtr("ABC123"),
		Status:    "active",
		Website:   strPtr("https://example.com"),
		CreatedAt: time.Now(),
		Address:   testAddress{City: "Lima"},
	}

	if err := Valid(user, "Email"); err != nil {
		t.Fatalf("expected Email to be skipped, got error: %v", err)
	}
}

func TestSelectField(t *testing.T) {
	user := testUser{
		Name:      "Kevin",
		Email:     "bad-email",
		Age:       10,
		Code:      strPtr("ABC123"),
		Status:    "invalid",
		Website:   strPtr("bad-url"),
		CreatedAt: time.Now(),
		Address:   testAddress{City: "Lima"},
	}

	if err := ValidSelect(user, "Name"); err != nil {
		t.Fatalf("expected only Name to be validated, got error: %v", err)
	}
}

func TestSelectInvalidField(t *testing.T) {
	user := testUser{
		Name:      "Kevin",
		Email:     "bad-email",
		Age:       10,
		Code:      strPtr("ABC123"),
		Status:    "invalid",
		Website:   strPtr("bad-url"),
		CreatedAt: time.Now(),
		Address:   testAddress{City: "Lima"},
	}

	err := ValidSelect(user, "Email")
	if err == nil {
		t.Fatal("expected Email validation error")
	}

	if !strings.Contains(err.Error(), "Email") {
		t.Fatalf("expected error to contain Email, got: %v", err)
	}
}

func TestNestedStructValidation(t *testing.T) {
	type userWithAddress struct {
		Name    string `chk:"required"`
		Address testAddress
	}

	err := Valid(userWithAddress{
		Name:    "Kevin",
		Address: testAddress{City: "Li"},
	})

	if err == nil {
		t.Fatal("expected nested Address.City validation error")
	}

	if !strings.Contains(err.Error(), "Address.City") {
		t.Fatalf("expected error to contain Address.City, got: %v", err)
	}
}

func TestSkipNestedFieldByPath(t *testing.T) {
	type userWithAddress struct {
		Name    string `chk:"required"`
		Address testAddress
	}

	err := Valid(userWithAddress{
		Name:    "Kevin",
		Address: testAddress{City: "Li"},
	}, "Address.City")

	if err != nil {
		t.Fatalf("expected Address.City to be skipped, got: %v", err)
	}
}

func TestSelectNestedFieldByPath(t *testing.T) {
	type userWithAddress struct {
		Name    string `chk:"required"`
		Address testAddress
	}

	err := ValidSelect(userWithAddress{
		Name:    "",
		Address: testAddress{City: "Li"},
	}, "Address.City")

	if err == nil {
		t.Fatal("expected Address.City validation error")
	}

	if strings.Contains(err.Error(), "Name") {
		t.Fatalf("expected Name to be ignored, got: %v", err)
	}
}

func TestLengthMinMaxString(t *testing.T) {
	type dto struct {
		Code string `chk:"len=4"`
		Name string `chk:"min=2 max=5"`
	}

	if err := Valid(dto{Code: "ABCD", Name: "Kevin"}); err != nil {
		t.Fatalf("expected valid lengths, got: %v", err)
	}

	if err := Valid(dto{Code: "ABC", Name: "K"}); err == nil {
		t.Fatal("expected length/min errors")
	}
}

func TestMinMaxNumbers(t *testing.T) {
	type dto struct {
		Age int     `chk:"min=18 max=65"`
		Pay float64 `chk:"gte=1000 lte=5000"`
	}

	if err := Valid(dto{Age: 30, Pay: 2000}); err != nil {
		t.Fatalf("expected valid numbers, got: %v", err)
	}

	if err := Valid(dto{Age: 10, Pay: 7000}); err == nil {
		t.Fatal("expected numeric range errors")
	}
}

func TestComparators(t *testing.T) {
	type dto struct {
		A int `chk:"gt=10"`
		B int `chk:"gte=10"`
		C int `chk:"lt=10"`
		D int `chk:"lte=10"`
	}

	if err := Valid(dto{A: 11, B: 10, C: 9, D: 10}); err != nil {
		t.Fatalf("expected valid comparators, got: %v", err)
	}

	if err := Valid(dto{A: 10, B: 9, C: 10, D: 11}); err == nil {
		t.Fatal("expected comparator errors")
	}
}

func TestEmailUUIDURLIP(t *testing.T) {
	type dto struct {
		Email string `chk:"email"`
		ID    string `chk:"uuid"`
		URL   string `chk:"url"`
		IP    string `chk:"ip"`
		IPv4  string `chk:"ipv4"`
		IPv6  string `chk:"ipv6"`
	}

	valid := dto{
		Email: "test@example.com",
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		URL:   "https://example.com",
		IP:    "192.168.1.1",
		IPv4:  "192.168.1.1",
		IPv6:  "2001:db8::1",
	}

	if err := Valid(valid); err != nil {
		t.Fatalf("expected valid network fields, got: %v", err)
	}

	invalid := dto{
		Email: "bad",
		ID:    "bad",
		URL:   "bad",
		IP:    "bad",
		IPv4:  "2001:db8::1",
		IPv6:  "192.168.1.1",
	}

	if err := Valid(invalid); err == nil {
		t.Fatal("expected network validation errors")
	}
}

func TestStringFormats(t *testing.T) {
	type dto struct {
		Alpha    string `chk:"alpha"`
		AlphaNum string `chk:"alphanum"`
		Num      string `chk:"num"`
		Decimal  string `chk:"decimal"`
		Lower    string `chk:"lower"`
		Upper    string `chk:"upper"`
	}

	valid := dto{
		Alpha:    "José",
		AlphaNum: "José123",
		Num:      "12345",
		Decimal:  "123.45",
		Lower:    "hello",
		Upper:    "HELLO",
	}

	if err := Valid(valid); err != nil {
		t.Fatalf("expected valid string formats, got: %v", err)
	}

	invalid := dto{
		Alpha:    "José123",
		AlphaNum: "José-123",
		Num:      "123a",
		Decimal:  "123,45",
		Lower:    "Hello",
		Upper:    "Hello",
	}

	if err := Valid(invalid); err == nil {
		t.Fatal("expected string format errors")
	}
}

func TestOneOfPrefixSuffixContains(t *testing.T) {
	type dto struct {
		Status string `chk:"oneof=active,inactive,pending"`
		Code   string `chk:"prefix=USR- suffix=-PE contains=123"`
	}

	if err := Valid(dto{Status: "active", Code: "USR-123-PE"}); err != nil {
		t.Fatalf("expected valid string rules, got: %v", err)
	}

	if err := Valid(dto{Status: "deleted", Code: "ADM-999-AR"}); err == nil {
		t.Fatal("expected oneof/prefix/suffix/contains errors")
	}
}

func TestDateTimeValidation(t *testing.T) {
	type dto struct {
		Date     string    `chk:"date"`
		Time     string    `chk:"time"`
		DateTime string    `chk:"datetime"`
		Created  time.Time `chk:"datetime"`
	}

	valid := dto{
		Date:     "2026-04-30",
		Time:     "15:04:05",
		DateTime: "2026-04-30 15:04:05",
		Created:  time.Now(),
	}

	if err := Valid(valid); err != nil {
		t.Fatalf("expected valid date/time, got: %v", err)
	}

	invalid := dto{
		Date:     "30-04-2026",
		Time:     "15:04",
		DateTime: "2026-04-30T15:04:05",
		Created:  time.Now(),
	}

	if err := Valid(invalid); err == nil {
		t.Fatal("expected date/time errors")
	}
}

func TestCustomValidator(t *testing.T) {
	v := New()

	v.Register("startsx", func(f Field) error {
		value, ok := f.Value.(string)
		if !ok {
			return nil
		}

		if !strings.HasPrefix(value, "x") {
			return ErrInvalidInput
		}

		return nil
	})

	type dto struct {
		Code string `chk:"startsx"`
	}

	if err := v.Struct(dto{Code: "x123"}); err != nil {
		t.Fatalf("expected custom validator valid, got: %v", err)
	}

	if err := v.Struct(dto{Code: "a123"}); err == nil {
		t.Fatal("expected custom validator error")
	}
}

func TestUnknownValidator(t *testing.T) {
	type dto struct {
		Name string `chk:"unknown"`
	}

	err := Valid(dto{Name: "Kevin"})
	if err == nil {
		t.Fatal("expected unknown validator error")
	}

	if !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected error to contain validator name, got: %v", err)
	}
}

func TestUnexportedFieldIgnored(t *testing.T) {
	type dto struct {
		Name   string `chk:"required"`
		secret string `chk:"required"`
	}

	if err := Valid(dto{Name: "Kevin"}); err != nil {
		t.Fatalf("expected unexported field to be ignored, got: %v", err)
	}
}

func TestTagDashIgnored(t *testing.T) {
	type dto struct {
		Name string `chk:"-"`
	}

	if err := Valid(dto{Name: ""}); err != nil {
		t.Fatalf("expected chk dash to ignore field, got: %v", err)
	}
}
