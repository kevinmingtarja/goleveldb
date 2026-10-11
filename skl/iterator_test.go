package skl

import (
	"bytes"
	"testing"
)

func newTestSkiplist(compare func(a, b []byte) int, keys []string) *Skiplist {
	s := NewSkipList(compare)
	for _, key := range keys {
		s.Put([]byte(key), []byte("value-"+key))
	}
	return s
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
			s := newTestSkiplist(tt.compare, []string{"d", "b", "f", "a", "e", "c"})

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

func TestSeekToLast(t *testing.T) {
	descending := func(a, b []byte) int {
		return -bytes.Compare(a, b)
	}
	keys := []string{"d", "b", "f", "a", "e", "c"}

	tests := []struct {
		name      string
		compare   func(a, b []byte) int
		keys      []string
		wantKey   string
		wantValid bool
	}{
		{"ascending/empty", bytes.Compare, nil, "", false},
		{"descending/empty", descending, nil, "", false},
		{"ascending/single", bytes.Compare, []string{"a"}, "a", true},
		{"descending/single", descending, []string{"a"}, "a", true},
		{"ascending/multiple", bytes.Compare, keys, "f", true},
		{"descending/multiple", descending, keys, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestSkiplist(tt.compare, tt.keys)

			iter := s.NewIterator()
			iter.SeekToLast()
			if got := iter.Valid(); got != tt.wantValid {
				t.Fatalf("valid = %t, want %t", got, tt.wantValid)
			}
			if !tt.wantValid {
				return
			}
			if got := string(iter.Key()); got != tt.wantKey {
				t.Fatalf("key = %q, want %q", got, tt.wantKey)
			}
			if got, want := string(iter.Value()), "value-"+tt.wantKey; got != want {
				t.Fatalf("value = %q, want %q", got, want)
			}
		})
	}
}

func TestPrev(t *testing.T) {
	descending := func(a, b []byte) int {
		return -bytes.Compare(a, b)
	}
	keys := []string{"d", "b", "f", "a", "e", "c"}

	tests := []struct {
		name      string
		compare   func(a, b []byte) int
		keys      []string
		target    string
		wantKey   string
		wantValid bool
	}{
		{"ascending/a", bytes.Compare, keys, "a", "", false},
		{"ascending/b", bytes.Compare, keys, "b", "a", true},
		{"ascending/c", bytes.Compare, keys, "c", "b", true},
		{"ascending/d", bytes.Compare, keys, "d", "c", true},
		{"ascending/e", bytes.Compare, keys, "e", "d", true},
		{"ascending/f", bytes.Compare, keys, "f", "e", true},
		{"descending/f", descending, keys, "f", "", false},
		{"descending/e", descending, keys, "e", "f", true},
		{"descending/d", descending, keys, "d", "e", true},
		{"descending/c", descending, keys, "c", "d", true},
		{"descending/b", descending, keys, "b", "c", true},
		{"descending/a", descending, keys, "a", "b", true},
		{"ascending/single", bytes.Compare, []string{"a"}, "a", "", false},
		{"descending/single", descending, []string{"a"}, "a", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestSkiplist(tt.compare, tt.keys)

			iter := s.NewIterator()
			iter.Seek([]byte(tt.target))
			if !iter.Valid() {
				t.Fatalf("seek to %q returned an invalid iterator", tt.target)
			}
			if got := string(iter.Key()); got != tt.target {
				t.Fatalf("seek to %q returned key %q", tt.target, got)
			}

			iter.Prev()
			if got := iter.Valid(); got != tt.wantValid {
				t.Fatalf("prev from %q: valid = %t, want %t", tt.target, got, tt.wantValid)
			}
			if !tt.wantValid {
				return
			}
			if got := string(iter.Key()); got != tt.wantKey {
				t.Fatalf("prev from %q: key = %q, want %q", tt.target, got, tt.wantKey)
			}
			if got, want := string(iter.Value()), "value-"+tt.wantKey; got != want {
				t.Fatalf("prev from %q: value = %q, want %q", tt.target, got, want)
			}
		})
	}
}
