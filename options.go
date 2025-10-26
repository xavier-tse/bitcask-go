package bitcask_go

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
