package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    StrArray
		wantErr bool
	}{
		{
			name:  "single string",
			input: `"test string"`,
			want:  StrArray{"test string"},
		},
		{
			name:  "string array",
			input: `["test1","test2"]`,
			want:  StrArray{"test1", "test2"},
		},
		{
			name:  "null",
			input: `null`,
			want:  StrArray{},
		},
		{
			name:  "empty array",
			input: `[]`,
			want:  StrArray{},
		},
		{
			name:    "invalid object",
			input:   `{}`,
			wantErr: true,
		},
		{
			name:    "non string array",
			input:   `[1,2,3]`,
			wantErr: true,
		},
		{
			name:    "mixed types",
			input:   `["test",1]`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   `- [1,2,3]`,
			wantErr: true,
		},
		{
			name:    "raw invalid string",
			input:   `hello`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa StrArray
			err := sa.UnmarshalJSON([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, sa)
		})
	}
}

func TestStrArray_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   StrArray
		want    string
		wantErr bool
	}{
		{
			name:  "nil slice",
			input: nil,
			want:  "null",
		},
		{
			name:  "empty slice",
			input: StrArray{},
			want:  "[]",
		},
		{
			name:  "single value",
			input: StrArray{"a"},
			want:  `["a"]`,
		},
		{
			name:  "multiple values",
			input: StrArray{"a", "b", "c"},
			want:  `["a","b","c"]`,
		},
		{
			name:  "with empty strings",
			input: StrArray{"a", "", "c"},
			want:  `["a","","c"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(data))

			if tt.input == nil {
				return
			}

			var roundtrip StrArray
			err = json.Unmarshal(data, &roundtrip)
			require.NoError(t, err)
			assert.Equal(t, tt.input, roundtrip)
		})
	}
}

func TestStrArray_Trimmed(t *testing.T) {
	input := StrArray{" a ", "b", "  c  "}
	expected := StrArray{"a", "b", "c"}

	result := input.Trimmed()
	assert.Equal(t, expected, result)
}

func TestStrArray_NonEmpty(t *testing.T) {
	input := StrArray{"a", "", "b", "", "c"}
	expected := StrArray{"a", "b", "c"}

	result := input.NonEmpty()
	assert.Equal(t, expected, result)
}

func TestStrArray_Unique(t *testing.T) {
	input := StrArray{"a", "b", "a", "c", "b"}
	expected := StrArray{"a", "b", "c"}

	result := input.Unique()
	assert.Equal(t, expected, result)
}
