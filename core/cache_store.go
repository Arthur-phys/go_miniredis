package core

import "miniredis/server"

type ShardedCacheStore struct {
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewShardedCacheStore() server.CacheStore {
	return &ShardedCacheStore{
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

func (c *ShardedCacheStore) LPush(key string, args ...string) error {
	return nil
}

func (c *ShardedCacheStore) LPop(key string) (string, error) {
	return "", nil
}

func (c *ShardedCacheStore) LLen(key string) (string, error) {
	return "", nil
}
