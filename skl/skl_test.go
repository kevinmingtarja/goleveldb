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

func TestSeek(t *testing.T) {
	descending := func(a, b []byte) int {
		return -bytes.Compare(a, b)
	}

	tests := []struct {
		name      string
		compare   func(a, b []byte) int
		target    string
		wantKey   string
		wantValid bool
	}{
		{"ascending/a", bytes.Compare, "a", "a", true},
		{"ascending/b", bytes.Compare, "b", "b", true},
		{"ascending/c", bytes.Compare, "c", "c", true},
		{"ascending/d", bytes.Compare, "d", "d", true},
		{"ascending/e", bytes.Compare, "e", "e", true},
		{"ascending/f", bytes.Compare, "f", "f", true},
		{"ascending/between", bytes.Compare, "bb", "c", true},
		{"ascending/past_end", bytes.Compare, "z", "", false},
		{"descending/f", descending, "f", "f", true},
		{"descending/e", descending, "e", "e", true},
		{"descending/d", descending, "d", "d", true},
		{"descending/c", descending, "c", "c", true},
		{"descending/b", descending, "b", "b", true},
		{"descending/a", descending, "a", "a", true},
		{"descending/between", descending, "bb", "b", true},
		{"descending/past_end", descending, "0", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSkipList(tt.compare)
			for _, key := range []string{"d", "b", "f", "a", "e", "c"} {
				s.Put([]byte(key), []byte("value-"+key))
			}

			iter := s.NewIterator()
			iter.Seek([]byte(tt.target))
			if got := iter.Valid(); got != tt.wantValid {
				t.Fatalf("seek to %q: valid = %t, want %t", tt.target, got, tt.wantValid)
			}
			if !tt.wantValid {
				return
			}
			if got := string(iter.Key()); got != tt.wantKey {
				t.Fatalf("seek to %q: key = %q, want %q", tt.target, got, tt.wantKey)
			}
			if got, want := string(iter.Value()), "value-"+tt.wantKey; got != want {
				t.Fatalf("seek to %q: value = %q, want %q", tt.target, got, want)
			}
		})
	}
}
