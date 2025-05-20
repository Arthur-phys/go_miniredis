package caches

import (
	"sync"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
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

func (c *ShardedCache) Get(key string) (string, error) {
	return "", nil
}

func (c *ShardedCache) Set(key string, value string) error {
	return nil
}

func (c *ShardedCache) RPush(key string, args ...string) error {
	return nil
}

func (c *ShardedCache) RPop(key string) (string, error) {
	return "", nil
}

func (c *ShardedCache) LPush(key string, args ...string) error {
	return nil
}

func (c *ShardedCache) LPop(key string) (string, error) {
	return "", nil
}

func (c *ShardedCache) LIndex(key string, index int) (string, error) {
	return "", nil
}

func (c *ShardedCache) LLen(key string) (int, error) {
	return 1, nil
}

func (c *ShardedCache) Del(key string) error {
	return nil
}

func (c *ShardedCache) Lock() {
	c.internalLock.Lock()
}

func (c *ShardedCache) Unlock() {
	c.internalLock.Unlock()
}
