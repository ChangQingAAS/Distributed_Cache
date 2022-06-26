package inMemory

import (
	"distribute_cache/cache/Pair"
	"distribute_cache/cache/scanner"
)

type inMemoryScanner struct {
	pair    Pair.Pair
	pairCh  chan *Pair.Pair
	closeCh chan struct{}
}

func (s *inMemoryScanner) Close() {
	close(s.closeCh)
}

func (s *inMemoryScanner) Scan() bool {
	pair, ok := <-s.pairCh
	if ok {
		s.pair.Key, s.pair.Value = pair.Key, pair.Value
	}
	return ok
}

func (s *inMemoryScanner) Key() string {
	return s.pair.Key
}

func (s *inMemoryScanner) Value() []byte {
	return s.pair.Value
}

func (c *inMemoryCache) NewScanner() scanner.Scanner {
	pairCh := make(chan *Pair.Pair)
	closeCh := make(chan struct{})

	go func() {
		defer close(pairCh)
		c.mutex.RLock()
		for key, value := range c.c {
			c.mutex.RUnlock()
			select {
			case <-closeCh:
				return
			case pairCh <- &Pair.Pair{Key: key, Value: value.val}:
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}()

	return &inMemoryScanner{Pair.Pair{}, pairCh, closeCh}
}
