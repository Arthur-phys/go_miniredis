// caches has all the caches implemented on this library.
// The biggest differences are in sharding or not for multiple access to a single dictonary.
//
// When using a cache, rather than instantiating it directly on a server, use a function like the NewSimpleCache.
// This is because a cache is created for every connection opened.
package caches

import (
	"fmt"
	"sync"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

// SimpleCache may be of utility to you if you don't care about performace or have very few clients
// otherwise, try to use another of the caches provided.
type SimpleCache struct {
	internalLock sync.Mutex
	dict         map[string]interface{}
}

func NewSimpleCache() interfaces.CacheStore {
	return &SimpleCache{
		sync.Mutex{},
		make(map[string]interface{}),
	}
}

func (c *SimpleCache) Get(key string) (string, e.Error) {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.(string); ok {
			return v, e.Error{}
		}
		return "", e.WrongType
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) Set(key string, value string) e.Error {
	c.dict[key] = value
	return e.Error{}
}

func (c *SimpleCache) RPush(key string, args ...string) e.Error {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			c.dict[key] = append(v, args...)
			return e.Error{}
		}
		return e.WrongType
	} else {
		c.dict[key] = args
	}
	return e.Error{}
}

func (c *SimpleCache) RPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			x, c.dict[key] = v[len(v)-1], v[:len(v)-1]
			if len(v)-1 == 0 {
				delete(c.dict, key)
			}
			return x, e.Error{}
		}
		return "", e.WrongType
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LPush(key string, args ...string) e.Error {
	i, j := 0, len(args)-1
	for i < j {
		args[i], args[j] = args[j], args[i]
		i++
		j--
	}
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			c.dict[key] = append(args, v...)
			return e.Error{}
		}
		return e.WrongType
	} else {
		c.dict[key] = args
	}
	return e.Error{}
}

func (c *SimpleCache) LPop(key string) (string, e.Error) {
	var x string
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			x, c.dict[key] = v[0], v[1:]
			if len(v)-1 == 0 {
				delete(c.dict, key)
			}
			return x, e.Error{}
		}
		return "", e.WrongType
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LIndex(key string, index int) (string, e.Error) {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			if len(v) > index && index >= 0 {
				return v[index], e.Error{}
			} else {
				err := e.IndexOutOfRange
				err.ExtraContext = map[string]string{"index": fmt.Sprintf("%d", index)}
				return "", err
			}
		}
		return "", e.WrongType
	} else {
		err := e.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) Del(key string) e.Error {
	delete(c.dict, key)
	return e.Error{}
}

func (c *SimpleCache) LLen(key string) (int, e.Error) {
	if v, ok := c.dict[key].([]string); ok {
		return len(v), e.Error{}
	}
	return 0, e.WrongType
}

func (c *SimpleCache) Lock() {
	c.internalLock.Lock()
}

func (c *SimpleCache) Unlock() {
	c.internalLock.Unlock()
}
