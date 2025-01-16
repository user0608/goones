package types

import (
	"reflect"
	"testing"
)

func TestStrArray_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		sa      *StrArray
		args    args
		want    StrArray
		wantErr bool
	}{
		{
			name: "simple string",
			args: args{
				data: []byte(`"test string"`),
			},
			sa:   &StrArray{},
			want: StrArray([]string{"test string"}),
		},
		{
			name: "array string",
			args: args{
				data: []byte(`["test1","test2"]`),
			},
			sa:   &StrArray{},
			want: StrArray([]string{"test1", "test2"}),
		},
		{
			name: "null array",
			args: args{
				data: []byte(`null`),
			},
			sa:   &StrArray{},
			want: StrArray([]string{}),
		},
		{
			name: "empty array",
			args: args{
				data: []byte(`[]`),
			},
			sa:   &StrArray{},
			want: StrArray([]string{}),
		},
		{
			name: "empty obj",
			args: args{
				data: []byte(`{}`),
			},
			sa:      &StrArray{},
			want:    StrArray([]string{}),
			wantErr: true,
		},
		{
			name: "unsopported type",
			args: args{
				data: []byte(`[1,2,3]`),
			},
			sa:      &StrArray{},
			want:    StrArray([]string{}),
			wantErr: true,
		},
		{
			name: "invalid JSON",
			args: args{
				data: []byte(`- [1,2,3]`),
			},
			sa:      &StrArray{},
			want:    StrArray([]string{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sa.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("StrArray.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(tt.sa, &tt.want) {
				t.Errorf("StrArray.UnmarshalJSON() not equal = %v, want %v", tt.sa, tt.want)
			}
		})
	}
}
