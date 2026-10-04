// Package key defines the internal database key format.
package key

import (
	"bytes"
	"encoding/binary"

	"github.com/kevinmingtarja/goleveldb/internal/assert"
)

// ValueType identifies a value or deletion entry.
type ValueType uint8

const (
	// https://github.com/google/leveldb/blob/7ee830d02b623e8ffe0b95d59a74db1e58da04c5/db/dbformat.h#L54
	TypeDeletion ValueType = 0
	TypeValue    ValueType = 1

	// We reserve eight bits for the type.
	MaxSequenceNumber uint64 = (1 << 56) - 1
)

// Encode returns an internal key containing the user key, sequence number, and value type.
func Encode(key []byte, seq uint64, typ ValueType) []byte {
	assert.True(seq <= MaxSequenceNumber)
	assert.True(typ <= TypeValue)
	// Add 8 more bytes (64 bits) at the end for sequence number + type
	out := make([]byte, len(key)+8)
	copy(out, key)
	binary.BigEndian.PutUint64(out[len(key):], (seq<<8)|uint64(typ))
	return out
}

// Decode extracts the user key, sequence number, and entry type from an internal key.
func Decode(encoded []byte) (userKey []byte, seq uint64, typ ValueType, err error) {
	panic("not implemented")
}

// Compare orders internal keys by ascending user key, descending sequence number, and descending type.
func Compare(a, b []byte) int {
	if cmp := bytes.Compare(a[:len(a)-8], b[:len(b)-8]); cmp != 0 {
		return cmp
	}
	// Reverse the cmp for descending sequence number and type
	return -bytes.Compare(a[len(a)-8:], b[len(b)-8:])
}
