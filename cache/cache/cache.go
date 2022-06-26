package cache

import (
	"distribute_cache/cache/inMemory"
	"distribute_cache/cache/rocksdb"
	"distribute_cache/cache/scanner"
	"distribute_cache/cache/status"
	"log"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() status.Stat
	NewScanner() scanner.Scanner
}

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "inmemory" {
		c = inMemory.NewInMemoryCache(ttl)
	} else if typ == "rocksdb" {
		c = rocksdb.NewRocksdbCache(ttl)
	}

	if c == nil {
		panic("unknown cache type " + typ)
	}

	log.Println(typ, "ready to serve")
	return c
}
