package pointers

import "time"

func Int(v int) *int {
	return &v
}

func Int8(v int8) *int8 {
	return &v
}

func Int16(v int16) *int16 {
	return &v
}

func Int32(v int32) *int32 {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Uint(v uint) *uint {
	return &v
}

func Uint8(v uint8) *uint8 {
	return &v
}

func Uint16(v uint16) *uint16 {
	return &v
}

func Uint32(v uint32) *uint32 {
	return &v
}

func Uint64(v uint64) *uint64 {
	return &v
}

func Float32(v float32) *float32 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

func CreatePointer[T any](v T) *T {
	return &v
}

func IntNilIfZero(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func Int8NilIfZero(v int8) *int8 {
	if v == 0 {
		return nil
	}
	return &v
}

func Int16NilIfZero(v int16) *int16 {
	if v == 0 {
		return nil
	}
	return &v
}

func Int32NilIfZero(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

func Int64NilIfZero(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func UintNilIfZero(v uint) *uint {
	if v == 0 {
		return nil
	}
	return &v
}

func Uint8NilIfZero(v uint8) *uint8 {
	if v == 0 {
		return nil
	}
	return &v
}

func Uint16NilIfZero(v uint16) *uint16 {
	if v == 0 {
		return nil
	}
	return &v
}

func Uint32NilIfZero(v uint32) *uint32 {
	if v == 0 {
		return nil
	}
	return &v
}

func Uint64NilIfZero(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func Float32NilIfZero(v float32) *float32 {
	if v == 0 {
		return nil
	}
	return &v
}

func Float64NilIfZero(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}

func StringNilIfZero(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func BoolNilIfZero(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}

func TimeNilIfZero(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return &v
}

func ZeroValueNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}
