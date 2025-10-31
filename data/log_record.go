package data

import "encoding/binary"

type LogRecordType byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// crc type keySize valueSize
//
//	4 + 1 +   5   +   5 = 15
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecord 的头部信息
type logRecordHeader struct {
	crc        uint32
	recordType LogRecordType
	keySize    uint32
	valueSize  uint32
}

type LogRecordPos struct {
	Fid    uint32
	Offset int64
}

func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	return nil, 0
}

// decodeLogRecord 对字节数组中的 Header 信息解码
func decodeLogRecord(buf []byte) (*logRecordHeader, int64) {
	return nil, 0
}

func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	return 0
}
