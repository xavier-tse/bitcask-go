package index

import (
	"bytes"

	"github.com/google/btree"
	"github.com/xavier-tse/bitcask-go/data"
)

type Indexer interface {
	// Put 向索引中存储 key 对应的数据位置信息
	Put(key []byte, pos *data.LogRecordPos) bool

	// Get 根据 key 取出对应的索引位置信息
	Get(key []byte) *data.LogRecordPos

	// Delete 根据 key 删除对应的索引位置信息
	Delete(key []byte) bool

	// Size 索引中的数据量
	Size() int

	// Iterator 索引迭代器
	Iterator(reverse bool) Iterator
}

type IndexType = int8

const (
	// Btree 索引
	Btree IndexType = iota + 1
)

// NewIndexer 根据 indexType 初始化索引
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

// Iterator 通用索引迭代器
type Iterator interface {
	// Rewind 重新回到迭代器起点，即第一个数据
	Rewind()

	// Seek 根据传入的 key 查找第一个大于(或小于)等于的目标 key，根据这个 key 开始遍历
	Seek(key []byte)

	// Next 跳转到下一个 key
	Next()

	// Valid 是否已经遍历完所有 key，用于退出
	Valid() bool

	// Key 当前位置的 key 数据
	Key() []byte

	// Value 当前位置的 value 数据
	Value() *data.LogRecordPos

	// Close 关闭迭代器，释放资源
	Close()
}
