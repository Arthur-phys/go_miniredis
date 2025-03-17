package caches

import (
	"miniredis/core/coreinterface"
	"sync"
)

type ShardedCacheStore struct {
	internalLock     sync.Mutex
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewShardedCacheStore() coreinterface.CacheStore {
	return &ShardedCacheStore{
		sync.Mutex{},
		make(map[string]string),
		make(map[string][]string),
	}
}

func (c *ShardedCacheStore) Get(key string) (string, bool) {
	return "", true
}

func (c *ShardedCacheStore) Set(key string, value string) error {
	return nil
}

func (c *ShardedCacheStore) RPush(key string, args ...string) error {
	return nil
}

func (c *ShardedCacheStore) RPop(key string) (string, error) {
	return "", nil
}

func (c *ShardedCacheStore) LPush(key string, args ...string) error {
	return nil
}

func (c *ShardedCacheStore) LPop(key string) (string, error) {
	return "", nil
}

func (c *ShardedCacheStore) LIndex(key string, index int) (string, bool) {
	return "", true
}

func (c *ShardedCacheStore) LLen(key string) (int, error) {
	return 1, nil
}

func (c *ShardedCacheStore) Lock() {
	c.internalLock.Lock()
}

func (c *ShardedCacheStore) Unlock() {
	c.internalLock.Unlock()
}
