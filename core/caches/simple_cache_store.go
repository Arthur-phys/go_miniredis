package caches

import (
	"miniredis/core/worker"
	e "miniredis/error"
	"sync"
)

type SimpleCacheStore struct {
	internalLock     sync.Mutex
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewSimpleCacheStore() worker.CacheStore {
	return &SimpleCacheStore{
		sync.Mutex{},
		make(map[string]string),
		make(map[string][]string),
	}
}

func (c *SimpleCacheStore) Get(key string) (v string, ok bool) {
	v, ok = c.simpleDictionary[key]
	return
}

func (c *SimpleCacheStore) Set(key string, value string) error {
	c.simpleDictionary[key] = value
	return nil
}

func (c *SimpleCacheStore) RPush(key string, args ...string) error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(c.arrayDictionary[key], args...)
	} else {
		c.arrayDictionary[key] = args
	}
	return nil
}

func (c *SimpleCacheStore) RPop(key string) (string, error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, v = v[len(v)-1], v[:len(v)-1]
		if len(v) == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, nil
	} else {
		return "", e.Error{}
	}
}

func (c *SimpleCacheStore) LPush(key string, args ...string) error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(args, c.arrayDictionary[key]...)
	} else {
		c.arrayDictionary[key] = args
	}
	return nil
}

func (c *SimpleCacheStore) LPop(key string) (string, error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, v = v[0], v[1:]
		if len(v) == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, nil
	} else {
		return "", e.Error{} //Change
	}
}

func (c *SimpleCacheStore) LIndex(key string, index int) (string, bool) {
	if v, ok := c.arrayDictionary[key]; ok && len(v) > index && index >= 0 {
		return v[index], true
	}
	return "", false //Change
}

func (c *SimpleCacheStore) LLen(key string) (int, error) {
	return len(c.arrayDictionary[key]), nil
}

func (c *SimpleCacheStore) Lock() {
	c.internalLock.Lock()
}

func (c *SimpleCacheStore) Unlock() {
	c.internalLock.Unlock()
}
