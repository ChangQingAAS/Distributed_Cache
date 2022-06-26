package rocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/opt/homebrew/Cellar/rocksdb/7.0.3/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/rocksdb/7.0.3 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"distribute_cache/cache/Pair"
	"time"
	"unsafe"
)

const BATCH_SIZE = 100

func flush_batch(db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var err *C.char
	C.rocksdb_write(db, o, b, &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_writebatch_clear(b)
}

func write_func(db *C.rocksdb_t, c chan *Pair.Pair, o *C.rocksdb_writeoptions_t) {
	count := 0
	t := time.NewTimer(time.Second)
	b := C.rocksdb_writebatch_create()

	for {
		select {
		case p := <-c:
			count++
			key := C.CString(p.Key)
			value := C.CBytes(p.Value)
			C.rocksdb_writebatch_put(b, key, C.size_t(len(p.Key)), (*C.char)(value), C.size_t(len(p.Value)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			// 如果达到batch个，写入batch个，并将count归0
			if count == BATCH_SIZE {
				flush_batch(db, b, o)
				count = 0
			}
			// 重置，更新时间
			// 拿走t.C，避免残留的C进入第二个case
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		// 到达时间间隔，取出t.C
		case <-t.C:
			// 写入当前count个
			if count != 0 {
				flush_batch(db, b, o)
				count = 0
			}
			// 重置时间
			t.Reset(time.Second)
		}
	}
}

func (c *rocksdbCache) Set(key string, value []byte) error {
	c.ch <- &Pair.Pair{key, value}
	return nil
}
