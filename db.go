package bitcask_go

import (
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/xavier-tse/bitcask-go/data"
	"github.com/xavier-tse/bitcask-go/index"
)

type DB struct {
	options    Options
	mu         *sync.RWMutex
	fileIds    []int // 文件id, 仅用于加载索引
	activeFile *data.DataFile
	olderFiles map[uint32]*data.DataFile
	index      index.Indexer
}

func Open(options Options) (*DB, error) {
	if err := checkOptions(options); err != nil {
		return nil, err
	}

	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	db := &DB{
		options:    options,
		mu:         new(sync.RWMutex),
		olderFiles: make(map[uint32]*data.DataFile),
		index:      index.NewIndexer(index.IndexType(options.IndexType)),
	}

	if err := db.loadDataFiles(); err != nil {
		return nil, err
	}

	if err := db.loadIndexFromDataFiles(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}

	pos, err := db.appendLogRecord(logRecord)
	if err != nil {
		return err
	}

	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdateFailed
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}

	logRecordPos := db.index.Get(key)
	if logRecordPos == nil {
		return nil, ErrKeyNotFound
	}

	var dataFile *data.DataFile
	if db.activeFile.FileId == logRecordPos.Fid {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFiles[logRecordPos.Fid]
	}
	if dataFile == nil {
		return nil, ErrDataFileNotFound
	}

	logRecord, _, err := dataFile.ReadLogRecord(logRecordPos.Offset)
	if err != nil {
		return nil, err
	}

	if logRecord.Type != data.LogRecordDeleted {
		return nil, ErrDataFileNotFound
	}
	return logRecord.Value, nil
}

func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecordPos, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	encRecord, size := data.EncodeLogRecord(logRecord)
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		db.olderFiles[db.activeFile.FileId] = db.activeFile

		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	if db.options.SyncWrites {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	pos := &data.LogRecordPos{
		Fid:    db.activeFile.FileId,
		Offset: writeOff,
	}
	return pos, nil
}

func (db *DB) setActiveDataFile() error {
	var initialFileId uint32 = 0
	if db.activeFile != nil {
		initialFileId = db.activeFile.FileId + 1
	}

	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileId)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}

func checkOptions(options Options) error {
	if options.DirPath == "" {
		return errors.New("database dir is empty")
	}
	if options.DataFileSize <= 0 {
		return errors.New("database data file size must be positive")
	}
	return nil
}

func (db *DB) loadDataFiles() error {
	dirEntries, err := os.ReadDir(db.options.DirPath)
	if err != nil {
		return err
	}

	var fileIds []int
	for _, entry := range dirEntries {
		if strings.HasSuffix(entry.Name(), data.DataFileNameSuffix) {
			splitNames := strings.Split(entry.Name(), ".")
			fileId, err := strconv.Atoi(splitNames[0])
			if err != nil {
				return ErrDataDirectoryCorrupted
			}
			fileIds = append(fileIds, fileId)
		}
	}
	sort.Ints(fileIds)
	db.fileIds = fileIds

	for i, fileId := range fileIds {
		dataFile, err := data.OpenDataFile(db.options.DirPath, uint32(fileId))
		if err != nil {
			return err
		}

		if i == len(fileIds)-1 {
			db.activeFile = dataFile
		} else {
			db.olderFiles[uint32(fileId)] = dataFile
		}
	}
	return nil
}

func (db *DB) loadIndexFromDataFiles() error {
	if len(db.fileIds) == 0 {
		return nil
	}

	for i, fid := range db.fileIds {
		var fileId = uint32(fid)
		var dataFile *data.DataFile
		if fileId == db.activeFile.FileId {
			dataFile = db.activeFile
		} else {
			dataFile = db.olderFiles[fileId]
		}

		var offset int64 = 0
		for {
			logRecord, size, err := dataFile.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			logRecordPos := &data.LogRecordPos{
				Fid:    fileId,
				Offset: offset,
			}
			if logRecord.Type == data.LogRecordDeleted {
				db.index.Delete(logRecord.Key)
			} else {
				db.index.Put(logRecord.Key, logRecordPos)
			}

			offset += size
		}

		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}
	return nil
}
