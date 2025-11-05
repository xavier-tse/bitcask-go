package bitcask_go

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xavier-tse/bitcask-go/utils"
)

func TestDB_NewIterator(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "bitcask-go-iterator-1")
	opt.DirPath = dir
	db, err := Open(opt)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	iterator := db.NewIterator(DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.Equal(t, false, iterator.Valid())
}

func TestDB_Iterator_One_Value(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "bitcask-go-iterator-2")
	opt.DirPath = dir
	db, err := Open(opt)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(utils.GetTestKey(10), utils.GetTestKey(10))
	assert.Nil(t, err)

	iterator := db.NewIterator(DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.Equal(t, true, iterator.Valid())
	assert.Equal(t, utils.GetTestKey(10), iterator.Key())
	val, err := iterator.Value()
	assert.Nil(t, err)
	assert.Equal(t, utils.GetTestKey(10), val)
}

func TestDB_Iterator_Multi_Value(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "bitcask-go-iterator-3")
	opt.DirPath = dir
	db, err := Open(opt)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put([]byte("aabc"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("abca"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("acba"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("cbaa"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("caba"), utils.RandomValue(10))
	assert.Nil(t, err)

	// 正向迭代
	iter1 := db.NewIterator(DefaultIteratorOptions)
	for iter1.Rewind(); iter1.Valid(); iter1.Next() {
		assert.NotNil(t, iter1.Key())
	}
	iter1.Rewind()
	for iter1.Seek([]byte("c")); iter1.Valid(); iter1.Next() {
		assert.NotNil(t, iter1.Key())
	}

	// 反向迭代
	iterOpts1 := DefaultIteratorOptions
	iterOpts1.Reverse = true
	iter2 := db.NewIterator(iterOpts1)
	for iter2.Rewind(); iter2.Valid(); iter2.Next() {
		assert.NotNil(t, iter2.Key())
	}
	iter2.Rewind()
	for iter2.Seek([]byte("c")); iter2.Valid(); iter2.Next() {
		assert.NotNil(t, iter2.Key())
	}

	// 指定了 prefix
	iterOpts2 := DefaultIteratorOptions
	iterOpts2.Prefix = []byte("aa")
	iter3 := db.NewIterator(iterOpts2)
	for iter3.Rewind(); iter3.Valid(); iter3.Next() {
		t.Log(string(iter3.Key()))
	}
}
