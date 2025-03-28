package caches

import (
	"miniredis/core/coreinterface"
	e "miniredis/error"
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

func (c *ShardedCacheStore) Get(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCacheStore) Set(key string, value string) e.Error {
	return e.Error{}
}

func (c *ShardedCacheStore) RPush(key string, args ...string) e.Error {
	return e.Error{}
}

func (c *ShardedCacheStore) RPop(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCacheStore) LPush(key string, args ...string) e.Error {
	return e.Error{}
}

func (c *ShardedCacheStore) LPop(key string) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCacheStore) LIndex(key string, index int) (string, e.Error) {
	return "", e.Error{}
}

func (c *ShardedCacheStore) LLen(key string) (int, e.Error) {
	return 1, e.Error{}
}

func (c *ShardedCacheStore) Lock() {
	c.internalLock.Lock()
}

func (c *ShardedCacheStore) Unlock() {
	c.internalLock.Unlock()
}
