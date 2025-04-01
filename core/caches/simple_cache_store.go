package caches

import (
	"fmt"
	"miniredis/core/coreinterface"
	e "miniredis/error"
	"sync"
)

type SimpleCacheStore struct {
	internalLock     sync.Mutex
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewSimpleCacheStore() coreinterface.CacheStore {
	return &SimpleCacheStore{
		sync.Mutex{},
		make(map[string]string),
		make(map[string][]string),
	}
}

func (c *SimpleCacheStore) Get(key string) (string, e.Error) {
	if v, ok := c.simpleDictionary[key]; ok {
		return v, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCacheStore) Set(key string, value string) e.Error {
	c.simpleDictionary[key] = value
	return e.Error{}
}

func (c *SimpleCacheStore) RPush(key string, args ...string) e.Error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(c.arrayDictionary[key], args...)
	} else {
		c.arrayDictionary[key] = args
	}
	return e.Error{}
}

func (c *SimpleCacheStore) RPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, v = v[len(v)-1], v[:len(v)-1]
		if len(v) == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCacheStore) LPush(key string, args ...string) e.Error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(args, c.arrayDictionary[key]...)
	} else {
		c.arrayDictionary[key] = args
	}
	return e.Error{}
}

func (c *SimpleCacheStore) LPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, v = v[0], v[1:]
		if len(v) == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCacheStore) LIndex(key string, index int) (string, e.Error) {
	if v, ok := c.arrayDictionary[key]; ok && len(v) > index && index >= 0 {
		return v[index], e.Error{}
	} else if len(v) > index || index < 0 {
		err := e.IndexOutOfRange
		err.ExtraContext = map[string]string{"index": fmt.Sprintf("%d", index)}
		return "", err
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCacheStore) LLen(key string) (int, e.Error) {
	return len(c.arrayDictionary[key]), e.Error{}
}

func (c *SimpleCacheStore) Lock() {
	c.internalLock.Lock()
}

func (c *SimpleCacheStore) Unlock() {
	c.internalLock.Unlock()
}
