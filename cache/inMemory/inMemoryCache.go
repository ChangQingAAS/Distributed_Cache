package inMemory

import (
	"distribute_cache/cache/status"
	"fmt"
	"sync"
	"time"
)

type value struct {
	val     []byte
	created time.Time
}

type inMemoryCache struct {
	c      map[string]value
	mutex  sync.RWMutex
	status status.Stat
	ttl    time.Duration
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	//fmt.Printf("Set(): key: %s, value: %s\n", k, v)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//tmp, exist := c.c[k]
	//if exist {
	//	c.status.DelStatus(k, tmp)
	//}
	c.c[k] = value{v, time.Now()}
	/*
		每次Set的键值对是最先超时的，也就是先进先出的缓存淘汰策略（FIFO）
		如果每次Get也去更新缓存时间，那么等于实现了一个最近最少使用策略（LRU）
		如果不用时间而用计数器，在每次GET时都让键值对的计数值加一，那么可以实现一个最少使用频率策略（Least Frequently Used,LFU)
	*/
	c.status.AddStatus(k, v)

	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.c[k].val, nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, exist := c.c[k]
	if exist {
		delete(c.c, k)
		fmt.Println("before status is ", c.status)
		c.status.DelStatus(k, value.val)
		fmt.Println("after status is ", c.status)
	}
	return nil
}

func (c *inMemoryCache) GetStat() status.Stat {
	return c.status
}

func NewInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{
		c:      make(map[string]value),
		mutex:  sync.RWMutex{},
		status: status.Stat{},
		ttl:    time.Duration(ttl) * time.Second,
	}

	if ttl > 0 {
		go c.expirer()
	}

	return c
}

func (c *inMemoryCache) expirer() {
	for {
		time.Sleep(c.ttl)
		c.mutex.Lock()
		for k, v := range c.c {
			c.mutex.Unlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.Lock()
		}
		c.mutex.Unlock()
	}
}
