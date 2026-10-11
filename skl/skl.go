// Package skl implements a skip list.
package skl

import (
	"math/rand/v2"
	"sync/atomic"

	"github.com/kevinmingtarja/goleveldb/internal/assert"
)

const (
	// https://github.com/google/leveldb/blob/7ee830d02b623e8ffe0b95d59a74db1e58da04c5/db/skiplist.h#L100
	maxHeight = 12
)

type node struct {
	key   []byte
	value []byte
	next  [maxHeight]atomic.Pointer[node]
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

// Return the first node whose key is greater than or equal to key,
// or nil if no such node exists.
//
// If prev is non-nil, fill prev[level] with pointer to
// previous node at "level" for every level in [0..max_height-1].
func (s *Skiplist) findGreaterOrEqual(key []byte, prev []*node) *node {
	level := s.getHeight() - 1
	x := s.head
	for {
		next := x.next[level].Load()
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
			}
			// Go down next level
			level--
		}
	}
}

// Return the latest node with a key < key,
// or head if no such node exists.
func (s *Skiplist) findLessThan(key []byte) *node {
	level := s.getHeight() - 1
	x := s.head
	for {
		assert.True(x == s.head || s.compare(x.key, key) < 0)
		next := x.next[level].Load()
		if next == nil || s.compare(key, next.key) <= 0 {
			if level == 0 {
				return x
			}
			// Go down next level
			level--
		} else {
			x = next
		}
	}
}

// Return the last node in the list,
// or head if list is empty.
func (s *Skiplist) findLast() *node {
	level := s.getHeight() - 1
	x := s.head
	for {
		next := x.next[level].Load()
		if next == nil {
			if level == 0 {
				return x
			}
			// Go down next level
			level--
		} else {
			x = next
		}
	}
}

func (s *Skiplist) randomHeight() int {
	const kBranching = 4
	height := 1
	// Increase height with probability 1 in kBranching
	for height < maxHeight && rand.IntN(kBranching) == 0 {
		height++
	}
	return height
}

func (s *Skiplist) Put(key, value []byte) {
	var prev [maxHeight]*node
	x := s.findGreaterOrEqual(key, prev[:])

	// We don't allow duplicate insertions.
	assert.True(x == nil || s.compare(key, x.key) != 0)

	height := s.randomHeight()
	if height > int(s.getHeight()) {
		for i := int(s.getHeight()); i < height; i++ {
			prev[i] = s.head
		}
		s.height.Store(int32(height))
	}

	x = newNode(key, value)
	for i := range height {
		x.next[i].Store(prev[i].next[i].Load())
		prev[i].next[i].Store(x)
	}
}

func (s *Skiplist) NewIterator() *Iterator {
	return &Iterator{
		list: s,
	}
}
