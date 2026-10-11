package skl

import (
	"bytes"
	"testing"
)

func TestContains(t *testing.T) {
	descending := func(a, b []byte) int {
		return -bytes.Compare(a, b)
	}
	keys := []string{"d", "b", "f", "a", "e", "c"}

	tests := []struct {
		name    string
		compare func(a, b []byte) int
		keys    []string
		target  string
		want    bool
	}{
		{"ascending/empty", bytes.Compare, nil, "a", false},
		{"descending/empty", descending, nil, "a", false},
		{"ascending/first", bytes.Compare, keys, "a", true},
		{"ascending/middle", bytes.Compare, keys, "c", true},
		{"ascending/last", bytes.Compare, keys, "f", true},
		{"ascending/between", bytes.Compare, keys, "bb", false},
		{"ascending/before_first", bytes.Compare, keys, "0", false},
		{"ascending/past_end", bytes.Compare, keys, "z", false},
		{"descending/first", descending, keys, "f", true},
		{"descending/middle", descending, keys, "c", true},
		{"descending/last", descending, keys, "a", true},
		{"descending/between", descending, keys, "bb", false},
		{"descending/before_first", descending, keys, "z", false},
		{"descending/past_end", descending, keys, "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestSkiplist(tt.compare, tt.keys)
			if got := s.Contains([]byte(tt.target)); got != tt.want {
				t.Fatalf("Contains(%q) = %t, want %t", tt.target, got, tt.want)
			}
		})
	}
}

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
