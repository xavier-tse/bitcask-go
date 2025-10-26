package data

import "github.com/xavier-tse/bitcask-go/fio"

const DataFileNameSuffix = ".data"

type DataFile struct {
	FileId    uint32
	WriteOff  int64
	ioManager fio.IOManager
}

func OpenDataFile(dirPath string, fileId uint32) (*DataFile, error) {
	return nil, nil
}

func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	return nil, 0, nil
}

func (df *DataFile) Write(buf []byte) error {
	return nil
}

func (df *DataFile) Sync() error {
	return nil
}
