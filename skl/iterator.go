package skl

import "github.com/kevinmingtarja/goleveldb/internal/assert"

type Iterator struct {
	list *Skiplist
	nd   *node
}

func (it *Iterator) Seek(target []byte) {
	it.nd = it.list.findGreaterOrEqual(target, nil)
}

func (it *Iterator) SeekToFirst() {
	it.nd = it.list.head.next[0].Load()
}

func (it *Iterator) Next() {
	assert.True(it.Valid())
	it.nd = it.nd.next[0].Load()
}

func (it *Iterator) Valid() bool {
	return it.nd != nil
}

func (it *Iterator) Key() []byte {
	assert.True(it.Valid())
	return it.nd.key
}

func (it *Iterator) Value() []byte {
	assert.True(it.Valid())
	return it.nd.value
}
