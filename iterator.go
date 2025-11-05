package bitcask_go

import (
	"bytes"

	"github.com/xavier-tse/bitcask-go/index"
)

// Iterator 迭代器
type Iterator struct {
	indexIter index.Iterator // 索引迭代器
	db        *DB
	options   IteratorOptions
}

func (db *DB) NewIterator(opts IteratorOptions) *Iterator {
	indexIter := db.index.Iterator(opts.Reverse)
	return &Iterator{
		db:        db,
		indexIter: indexIter,
		options:   opts,
	}
}

// Rewind 重新回到迭代器起点，即第一个数据
func (it *Iterator) Rewind() {
	it.indexIter.Rewind()
	it.skip2Next()
}

// Seek 根据传入的 key 查找第一个大于(或小于)等于的目标 key，根据这个 key 开始遍历
func (it *Iterator) Seek(key []byte) {
	it.indexIter.Seek(key)
	it.skip2Next()
}

// Next 跳转到下一个 key
func (it *Iterator) Next() {
	it.indexIter.Next()
	it.skip2Next()
}

// Valid 是否已经遍历完所有 key，用于退出
func (it *Iterator) Valid() bool {
	return it.indexIter.Valid()
}

// Key 当前位置的 key 数据
func (it *Iterator) Key() []byte {
	return it.indexIter.Key()
}

// Value 当前位置的 value 数据
func (it *Iterator) Value() ([]byte, error) {
	logRecordPos := it.indexIter.Value()
	it.db.mu.RLock()
	defer it.db.mu.RUnlock()
	return it.db.getValueByPosition(logRecordPos)
}

// Close 关闭迭代器，释放资源
func (it *Iterator) Close() {
	it.indexIter.Close()
}

func (it *Iterator) skip2Next() {
	prefixLen := len(it.options.Prefix)
	if prefixLen == 0 {
		return
	}

	for ; it.indexIter.Valid(); it.indexIter.Next() {
		key := it.indexIter.Key()
		if prefixLen <= len(key) && bytes.Compare(it.options.Prefix, key[:prefixLen]) == 0 {
			break
		}
	}
}
