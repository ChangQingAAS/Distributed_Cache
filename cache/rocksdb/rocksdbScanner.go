package rocksdb

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/7.0.3/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/7.0.3 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"distribute_cache/cache/scanner"
	"unsafe"
)

type rocksdbScanner struct {
	iterator    *C.rocksdb_iterator_t
	initialized bool
}

func (s *rocksdbScanner) Close() {
	C.rocksdb_iter_destroy(s.iterator)
}

func (s *rocksdbScanner) Scan() bool {
	if !s.initialized {
		C.rocksdb_iter_seek_to_first(s.iterator)
		s.initialized = true
	} else {
		C.rocksdb_iter_next(s.iterator)
	}

	return C.rocksdb_iter_valid(s.iterator) != 0
}

func (s *rocksdbScanner) Key() string {
	var length C.size_t // an unsigned integer data type
	key := C.rocksdb_iter_key(s.iterator, &length)

	return C.GoString(key)
}

func (s *rocksdbScanner) Value() []byte {
	var length C.size_t
	value := C.rocksdb_iter_value(s.iterator, &length)

	// 在这个指针，读取length长度的byte
	return C.GoBytes(unsafe.Pointer(value), C.int(length))
}

func (c *rocksdbCache) NewScanner() scanner.Scanner {
	return &rocksdbScanner{C.rocksdb_create_iterator(c.db, c.ro), false}
}
