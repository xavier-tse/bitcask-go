package index

import (
	"bytes"

	"github.com/google/btree"
	"github.com/xavier-tse/bitcask-go/data"
)

type Indexer interface {
	Put(key []byte, pos *data.LogRecordPos) bool
	Get(key []byte) *data.LogRecordPos
	Delete(key []byte) bool
}

type IndexType = int8

const (
	Btree IndexType = iota + 1
)

func NewIndexer(typ IndexType) Indexer {
	switch typ {
	case Btree:
		return NewBTree()
	default:
		panic("unsupported index type")
	}
}

type Item struct {
	key []byte
	pos *data.LogRecordPos
}

func (it *Item) Less(bi btree.Item) bool {
	return bytes.Compare(it.key, bi.(*Item).key) == -1
}
