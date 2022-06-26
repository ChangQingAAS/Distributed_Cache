package rocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/7.0.3/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/7.0.3 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"distribute_cache/cache/status"
	"regexp"
	"strconv"
	"unsafe"
)

func (c *rocksdbCache) GetStat() status.Stat {
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))
	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))
	p := C.GoString(v)
	r := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := status.Stat{}
	for _, submatches := range r.FindAllStringSubmatch(p, -1) {
		if submatches[1] == " # entries" {
			s.Count, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw key size" {
			s.KeySize, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw value size" {
			s.ValueSize, _ = strconv.ParseInt(submatches[2], 10, 64)
		}
	}
	return s
}
