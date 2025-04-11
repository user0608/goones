package pointers_test

import (
	"testing"
	"time"

	"github.com/user0608/goones/tools/pointers"
)

func TestPrimitives(t *testing.T) {
	if v := pointers.Int(5); *v != 5 {
		t.Errorf("Int(5) = %d, want 5", *v)
	}
	if v := pointers.Int8(5); *v != 5 {
		t.Errorf("Int8(5) = %d, want 5", *v)
	}
	if v := pointers.Int16(5); *v != 5 {
		t.Errorf("Int16(5) = %d, want 5", *v)
	}
	if v := pointers.Int32(5); *v != 5 {
		t.Errorf("Int32(5) = %d, want 5", *v)
	}
	if v := pointers.Int64(5); *v != 5 {
		t.Errorf("Int64(5) = %d, want 5", *v)
	}
	if v := pointers.Uint(5); *v != 5 {
		t.Errorf("Uint(5) = %d, want 5", *v)
	}
	if v := pointers.Uint8(5); *v != 5 {
		t.Errorf("Uint8(5) = %d, want 5", *v)
	}
	if v := pointers.Uint16(5); *v != 5 {
		t.Errorf("Uint16(5) = %d, want 5", *v)
	}
	if v := pointers.Uint32(5); *v != 5 {
		t.Errorf("Uint32(5) = %d, want 5", *v)
	}
	if v := pointers.Uint64(5); *v != 5 {
		t.Errorf("Uint64(5) = %d, want 5", *v)
	}
	if v := pointers.Float32(5.5); *v != 5.5 {
		t.Errorf("Float32(5.5) = %f, want 5.5", *v)
	}
	if v := pointers.Float64(5.5); *v != 5.5 {
		t.Errorf("Float64(5.5) = %f, want 5.5", *v)
	}
	if v := pointers.String("hello"); *v != "hello" {
		t.Errorf(`String("hello") = %s, want "hello"`, *v)
	}
	if v := pointers.Bool(true); *v != true {
		t.Errorf("Bool(true) = %v, want true", *v)
	}
}

func TestNilIfZero(t *testing.T) {
	if v := pointers.IntNilIfZero(0); v != nil {
		t.Errorf("IntNilIfZero(0) should be nil")
	}
	if v := pointers.IntNilIfZero(1); *v != 1 {
		t.Errorf("IntNilIfZero(1) = %d, want 1", *v)
	}
	if v := pointers.StringNilIfZero(""); v != nil {
		t.Errorf(`StringNilIfZero("") should be nil`)
	}
	if v := pointers.StringNilIfZero("test"); *v != "test" {
		t.Errorf(`StringNilIfZero("test") = %s, want "test"`, *v)
	}
	if v := pointers.BoolNilIfZero(false); v != nil {
		t.Errorf("BoolNilIfZero(false) should be nil")
	}
	if v := pointers.BoolNilIfZero(true); *v != true {
		t.Errorf("BoolNilIfZero(true) = %v, want true", *v)
	}
	if v := pointers.Float64NilIfZero(0); v != nil {
		t.Errorf("Float64NilIfZero(0) should be nil")
	}
	if v := pointers.Float64NilIfZero(3.14); *v != 3.14 {
		t.Errorf("Float64NilIfZero(3.14) = %f, want 3.14", *v)
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	if v := pointers.Time(now); !v.Equal(now) {
		t.Errorf("Time(now) = %v, want %v", *v, now)
	}
}

func TestTimeNilIfZero(t *testing.T) {
	var zero time.Time
	if v := pointers.TimeNilIfZero(zero); v != nil {
		t.Errorf("TimeNilIfZero(zero) should be nil")
	}
	now := time.Now()
	if v := pointers.TimeNilIfZero(now); !v.Equal(now) {
		t.Errorf("TimeNilIfZero(now) = %v, want %v", *v, now)
	}
}

func TestCreatePointer(t *testing.T) {
	i := 10
	if p := pointers.CreatePointer(i); *p != i {
		t.Errorf("CreatePointer(10) = %d, want 10", *p)
	}

	s := "hello"
	if p := pointers.CreatePointer(s); *p != s {
		t.Errorf(`CreatePointer("hello") = %s, want "hello"`, *p)
	}

	now := time.Now()
	if p := pointers.CreatePointer(now); !p.Equal(now) {
		t.Errorf("CreatePointer(time) = %v, want %v", *p, now)
	}
}

func TestZeroValueNil(t *testing.T) {
	if p := pointers.ZeroValueNil(0); p != nil {
		t.Errorf("ZeroValueNil(0) should be nil")
	}
	if p := pointers.ZeroValueNil(1); *p != 1 {
		t.Errorf("ZeroValueNil(1) = %d, want 1", *p)
	}

	if p := pointers.ZeroValueNil(""); p != nil {
		t.Errorf(`ZeroValueNil("") should be nil`)
	}
	if p := pointers.ZeroValueNil("go"); *p != "go" {
		t.Errorf(`ZeroValueNil("go") = %s, want "go"`, *p)
	}

	var zeroTime time.Time
	if p := pointers.ZeroValueNil(zeroTime); p != nil {
		t.Errorf("ZeroValueNil(zero time) should be nil")
	}
	now := time.Now()
	if p := pointers.ZeroValueNil(now); !p.Equal(now) {
		t.Errorf("ZeroValueNil(now) = %v, want %v", *p, now)
	}
}
