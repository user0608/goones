package types

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDArray_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		sa      *UUIDArray
		args    args
		want    UUIDArray
		wantErr bool
	}{
		{
			name: "simple string",
			args: args{
				data: []byte(`"220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7"`),
			},
			sa:   &UUIDArray{},
			want: UUIDArray([]uuid.UUID{uuid.MustParse("220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7")}),
		},
		{
			name: "array string",
			args: args{
				data: []byte(`["ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a","61151b3f-c62f-428f-a5af-e66fed737b9c"]`),
			},
			sa: &UUIDArray{},
			want: UUIDArray([]uuid.UUID{
				uuid.MustParse("ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a"),
				uuid.MustParse("61151b3f-c62f-428f-a5af-e66fed737b9c"),
			}),
		},
		{
			name: "null array",
			args: args{
				data: []byte(`null`),
			},
			sa:   &UUIDArray{},
			want: UUIDArray([]uuid.UUID{}),
		},
		{
			name: "empty array",
			args: args{
				data: []byte(`[]`),
			},
			sa:   &UUIDArray{},
			want: UUIDArray([]uuid.UUID{}),
		},
		{
			name: "empty obj",
			args: args{
				data: []byte(`{}`),
			},
			sa:      &UUIDArray{},
			want:    UUIDArray([]uuid.UUID{}),
			wantErr: true,
		},
		{
			name: "unsopported type",
			args: args{
				data: []byte(`[1,2,3]`),
			},
			sa:      &UUIDArray{},
			want:    UUIDArray([]uuid.UUID{}),
			wantErr: true,
		},
		{
			name: "unsopported type string",
			args: args{
				data: []byte(`hello`),
			},
			sa:      &UUIDArray{},
			want:    UUIDArray([]uuid.UUID{}),
			wantErr: true,
		},
		{
			name: "unsopported type string",
			args: args{
				data: []byte(`["hello","hello","hello"]`),
			},
			sa:      &UUIDArray{},
			want:    UUIDArray([]uuid.UUID{}),
			wantErr: true,
		},
		{
			name: "invalid JSON",
			args: args{
				data: []byte(`- [1,2,3]`),
			},
			sa:      &UUIDArray{},
			want:    UUIDArray([]uuid.UUID{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sa.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UUIDArray.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(tt.sa, &tt.want) {
				t.Errorf("UUIDArray.UnmarshalJSON() not equal = %v, want %v", tt.sa, tt.want)
			}
		})
	}
}
