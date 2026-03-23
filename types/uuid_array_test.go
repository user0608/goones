package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    UUIDArray
		wantErr bool
	}{
		{
			name:  "single uuid",
			input: `"220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7"`,
			want:  UUIDArray{uuid.MustParse("220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7")},
		},
		{
			name:  "uuid array",
			input: `["ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a","61151b3f-c62f-428f-a5af-e66fed737b9c"]`,
			want: UUIDArray{
				uuid.MustParse("ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a"),
				uuid.MustParse("61151b3f-c62f-428f-a5af-e66fed737b9c"),
			},
		},
		{
			name:  "null",
			input: `null`,
			want:  UUIDArray{},
		},
		{
			name:  "empty array",
			input: `[]`,
			want:  UUIDArray{},
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
			name:    "invalid raw string",
			input:   `hello`,
			wantErr: true,
		},
		{
			name:    "invalid uuid in array",
			input:   `["hello","hello"]`,
			wantErr: true,
		},
		{
			name:    "malformed json",
			input:   `- [1,2,3]`,
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid uuid",
			input:   `["ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a","invalid"]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa UUIDArray
			err := sa.UnmarshalJSON([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sa) != len(tt.want) {
				t.Fatalf("length mismatch: got %d, want %d", len(sa), len(tt.want))
			}

			for i := range sa {
				if sa[i] != tt.want[i] {
					t.Fatalf("element mismatch at %d: got %v, want %v", i, sa[i], tt.want[i])
				}
			}
		})
	}
}

func TestUUIDArray_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   UUIDArray
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
			input: UUIDArray{},
			want:  "[]",
		},
		{
			name: "single uuid",
			input: UUIDArray{
				uuid.MustParse("220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7"),
			},
			want: `["220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7"]`,
		},
		{
			name: "multiple uuids",
			input: UUIDArray{
				uuid.MustParse("ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a"),
				uuid.MustParse("61151b3f-c62f-428f-a5af-e66fed737b9c"),
			},
			want: `["ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a","61151b3f-c62f-428f-a5af-e66fed737b9c"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(data) != tt.want {
				t.Fatalf("unexpected output: got %s, want %s", data, tt.want)
			}

			var roundTrip UUIDArray
			if err := json.Unmarshal(data, &roundTrip); err != nil {
				t.Fatalf("roundtrip unmarshal failed: %v", err)
			}

			if len(roundTrip) != len(tt.input) {
				t.Fatalf("roundtrip length mismatch: got %d, want %d", len(roundTrip), len(tt.input))
			}

			for i := range roundTrip {
				if roundTrip[i] != tt.input[i] {
					t.Fatalf("roundtrip mismatch at %d: got %v, want %v", i, roundTrip[i], tt.input[i])
				}
			}
		})
	}
}

func TestUUIDArray_Unique(t *testing.T) {
	u1 := uuid.MustParse("220a4e82-0dc0-46bf-98f9-e0e3c06b2bf7")
	u2 := uuid.MustParse("ced02bc8-d1c7-4bcc-8eaa-a6ad58df251a")
	u3 := uuid.MustParse("61151b3f-c62f-428f-a5af-e66fed737b9c")

	tests := []struct {
		name  string
		input UUIDArray
		want  UUIDArray
	}{
		{
			name:  "no duplicates",
			input: UUIDArray{u1, u2, u3},
			want:  UUIDArray{u1, u2, u3},
		},
		{
			name:  "with duplicates",
			input: UUIDArray{u1, u2, u1, u3, u2},
			want:  UUIDArray{u1, u2, u3},
		},
		{
			name:  "all duplicates",
			input: UUIDArray{u1, u1, u1},
			want:  UUIDArray{u1},
		},
		{
			name:  "empty slice",
			input: UUIDArray{},
			want:  UUIDArray{},
		},
		{
			name:  "nil slice",
			input: nil,
			want:  UUIDArray{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Unique()
			assert.Equal(t, tt.want, result)
		})
	}
}
