package bitcask_go

import "os"

type Options struct {
	// 数据目录
	DirPath string

	// 数据文件大小
	DataFileSize int64

	// 每次写入数据是否持久化
	SyncWrites bool

	// 索引类型
	IndexType IndexType
}

// IteratorOptions 迭代器配置项
type IteratorOptions struct {
	// 遍历前序为指定值的 key，默认为空
	Prefix []byte

	// 是否反向遍历
	Reverse bool
}

type IndexType = int8

const (
	BTree IndexType = iota + 1
)

var DefaultOptions = Options{
	DirPath:      os.TempDir(),
	DataFileSize: 256 * 1024 * 1024, // 256MB
	SyncWrites:   false,
	IndexType:    BTree,
}

var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}
