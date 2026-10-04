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

			n := s.head.next[0].Load()
			for _, key := range tt.want {
				if n == nil {
					t.Fatalf("list ended before key %q", key)
				}
				if got := string(n.key); got != key {
					t.Fatalf("key = %q, want %q", got, key)
				}
				if got, want := string(n.value), "value-"+key; got != want {
					t.Fatalf("value for %q = %q, want %q", key, got, want)
				}
				if found := s.findGreaterOrEqual([]byte(key), nil); found != n {
					t.Fatalf("search for %q did not find its inserted node", key)
				}
				n = n.next[0].Load()
			}
			if n != nil {
				t.Fatalf("unexpected extra node with key %q", n.key)
			}
		})
	}
}
