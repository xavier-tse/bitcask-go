package bitcask_go

type Options struct {
	DirPath      string
	DataFileSize int64
	SyncWrites   bool
	IndexType    IndexType
}

type IndexType = int8

const (
	BTree IndexType = iota + 1
)
