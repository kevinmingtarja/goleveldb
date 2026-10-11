package skl

import (
	"bytes"
	"testing"
)

func TestPut(t *testing.T) {
	tests := []struct {
		name    string
		compare func(a, b []byte) int
		want    []string
	}{
		{
			name:    "ascending",
			compare: bytes.Compare,
			want:    []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name: "descending",
			compare: func(a, b []byte) int {
				return -bytes.Compare(a, b)
			},
			want: []string{"f", "e", "d", "c", "b", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSkipList(tt.compare)
			for _, key := range []string{"d", "b", "f", "a", "e", "c"} {
				s.Put([]byte(key), []byte("value-"+key))
			}

			iter := s.NewIterator()
			i := 0
			for iter.SeekToFirst(); iter.Valid(); iter.Next() {
				if got := string(iter.Key()); got != tt.want[i] {
					t.Fatalf("key = %q, want %q", got, tt.want[i])
				}
				if got, want := string(iter.Value()), "value-"+tt.want[i]; got != want {
					t.Fatalf("value for %q = %q, want %q", tt.want[i], got, want)
				}
				i++
			}
			if i != len(tt.want) {
				t.Fatalf("length mismatch: %d and %d", i, len(tt.want))
			}
		})
	}
}
