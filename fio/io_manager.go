package fio

const DataFilePerm = 0644

type IOManager interface {
	// Read 从文件给定位置读取对应数据
	Read([]byte, int64) (int, error)

	// Write 把字节数组写入到文件中
	Write([]byte) (int, error)

	// Sync 持久化数据
	Sync() error

	// Close 关闭文件
	Close() error

	// Size 得到文件大小
	Size() (int64, error)
}

func NewIOManager(fileName string) (IOManager, error) {
	return NewFileIOManager(fileName)
}
