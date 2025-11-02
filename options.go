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
