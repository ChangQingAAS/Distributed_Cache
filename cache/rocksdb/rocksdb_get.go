package rocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/7.0.3/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/7.0.3 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"unsafe"
)

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.err)
	if c.err != nil {
		return nil, errors.New(C.GoString(c.err))
	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil
}
