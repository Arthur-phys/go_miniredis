package caches

import (
	"fmt"
	"sync"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type SimpleCache struct {
	internalLock     sync.Mutex
	simpleDictionary map[string]string
	arrayDictionary  map[string][]string
}

func NewSimpleCache() interfaces.CacheStore {
	return &SimpleCache{
		sync.Mutex{},
		make(map[string]string),
		make(map[string][]string),
	}
}

func (c *SimpleCache) Get(key string) (string, e.Error) {
	if v, ok := c.simpleDictionary[key]; ok {
		return v, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) Set(key string, value string) e.Error {
	c.simpleDictionary[key] = value
	return e.Error{}
}

func (c *SimpleCache) RPush(key string, args ...string) e.Error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(c.arrayDictionary[key], args...)
	} else {
		c.arrayDictionary[key] = args
	}
	return e.Error{}
}

func (c *SimpleCache) RPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, c.arrayDictionary[key] = v[len(v)-1], v[:len(v)-1]
		if len(v)-1 == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LPush(key string, args ...string) e.Error {
	if _, ok := c.arrayDictionary[key]; ok {
		c.arrayDictionary[key] = append(args, c.arrayDictionary[key]...)
	} else {
		c.arrayDictionary[key] = args
	}
	return e.Error{}
}

func (c *SimpleCache) LPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.arrayDictionary[key]; ok {
		x, c.arrayDictionary[key] = v[0], v[1:]
		if len(v)-1 == 0 {
			delete(c.arrayDictionary, key)
		}
		return x, e.Error{}
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LIndex(key string, index int) (string, e.Error) {
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

func (c *SimpleCache) LLen(key string) (int, e.Error) {
	return len(c.arrayDictionary[key]), e.Error{}
}

func (c *SimpleCache) Lock() {
	c.internalLock.Lock()
}

func (c *SimpleCache) Unlock() {
	c.internalLock.Unlock()
}
