// Package skl implements a skip list.
package skl

import (
	"sync/atomic"
)

const (
	// https://github.com/google/leveldb/blob/7ee830d02b623e8ffe0b95d59a74db1e58da04c5/db/skiplist.h#L100
	maxHeight = 12
)

type node struct {
	key   []byte
	value []byte
	next  [maxHeight]*node
}

type Skiplist struct {
	height  atomic.Int32 // Current height. 1 <= height <= maxHeight.
	head    *node
	compare func(a, b []byte) int
}

func newNode(key, value []byte) *node {
	return &node{
		key:   key,
		value: value,
	}
}

func NewSkipList(compare func(a, b []byte) int) *Skiplist {
	head := newNode(nil, nil)
	s := &Skiplist{
		head:    head,
		compare: compare,
	}
	s.height.Store(1)
	return s
}

func (s *Skiplist) getHeight() int32 {
	return s.height.Load()
}

func (s *Skiplist) findGreaterOrEqual(key []byte, prev []*node) *node {
	level := s.getHeight() - 1
	x := s.head
	for {
		next := x.next[level]
		if next != nil && s.compare(key, next.key) > 0 {
			// key is greater, keep searching in the same level
			x = next
		} else {
			// Keep track of each level's last pointer
			if prev != nil {
				prev[level] = x
			}
			if level == 0 {
				return next
			} else {
				// Go down next level
				level--
			}
		}
	}
}
