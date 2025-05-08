package caches

import (
	"sync"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type ShardedCache struct {
	internalLock     sync.Mutex
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewShardedCache() interfaces.CacheStore {
	return &ShardedCache{
		sync.Mutex{},
		make(map[string]string),
		make(map[string][]string),
	}
}

func (c *ShardedCache) Get(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCache) Set(key string, value string) e.Error {
	return e.Error{}
}

func (c *ShardedCache) RPush(key string, args ...string) e.Error {
	return e.Error{}
}

func (c *ShardedCache) RPop(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCache) LPush(key string, args ...string) e.Error {
	return e.Error{}
}

func (c *ShardedCache) LPop(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCache) LIndex(key string, index int) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCache) LLen(key string) (int, e.Error) {
	return 1, e.Error{}
}

func (c *ShardedCache) Del(key string) e.Error {
	return e.Error{}
}

func (c *ShardedCache) Lock() {
	c.internalLock.Lock()
}

func (c *ShardedCache) Unlock() {
	c.internalLock.Unlock()
}
