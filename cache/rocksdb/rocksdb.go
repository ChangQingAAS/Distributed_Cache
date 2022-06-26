package rocksdb

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/7.0.3/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/7.0.3 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"distribute_cache/cache/Pair"
	"runtime"
)

type rocksdbCache struct {
	db  *C.rocksdb_t
	ro  *C.rocksdb_readoptions_t
	wo  *C.rocksdb_writeoptions_t
	err *C.char
	ch  chan *Pair.Pair
}

func NewRocksdbCache(ttl int) *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var err *C.char
	db := C.rocksdb_open_with_ttl(options, C.CString("/tmp/rocksdb"), C.int(ttl), &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_options_destroy(options)
	c := make(chan *Pair.Pair, 5000)
	wo := C.rocksdb_writeoptions_create()
	go write_func(db, c, wo)
	return &rocksdbCache{
		db,
		C.rocksdb_readoptions_create(),
		wo,
		err,
		c,
	}
}
