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
	"github.com/Arthur-phys/redigo/pkg/redigoerr"
)

// SimpleCache may be of utility to you if you don't care about performace or have very few clients
// otherwise, try to use another of the caches provided.
type SimpleCache struct {
	internalLock sync.Mutex
	dict         map[string]any
}

func NewSimpleCache() interfaces.CacheStore {
	return &SimpleCache{
		sync.Mutex{},
		make(map[string]any),
	}
}

func (c *SimpleCache) Get(key string) (string, error) {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.(string); ok {
			return v, nil
		}
		return "", redigoerr.WrongType
	} else {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) Set(key string, value string) error {
	c.dict[key] = value
	return nil
}

func (c *SimpleCache) RPush(key string, args ...string) error {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			c.dict[key] = append(v, args...)
			return nil
		}
		return redigoerr.WrongType
	} else {
		c.dict[key] = args
	}
	return nil
}

func (c *SimpleCache) RPop(key string) (string, error) {
	var x string
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			x, c.dict[key] = v[len(v)-1], v[:len(v)-1]
			if len(v)-1 == 0 {
				delete(c.dict, key)
			}
			return x, nil
		}
		return "", redigoerr.WrongType
	} else {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LPush(key string, args ...string) error {
	i, j := 0, len(args)-1
	for i < j {
		args[i], args[j] = args[j], args[i]
		i++
		j--
	}
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			c.dict[key] = append(args, v...)
			return nil
		}
		return redigoerr.WrongType
	} else {
		c.dict[key] = args
	}
	return nil
}

func (c *SimpleCache) LPop(key string) (string, error) {
	var x string
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			x, c.dict[key] = v[0], v[1:]
			if len(v)-1 == 0 {
				delete(c.dict, key)
			}
			return x, nil
		}
		return "", redigoerr.WrongType
	} else {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) LIndex(key string, index int) (string, error) {
	if v, ok := c.dict[key]; ok {
		if v, ok := v.([]string); ok {
			if len(v) > index && index >= 0 {
				return v[index], nil
			} else {
				err := redigoerr.IndexOutOfRangeErr
				err.ExtraContext = map[string]string{"index": fmt.Sprintf("%d", index)}
				return "", err
			}
		}
		return "", redigoerr.WrongType
	} else {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
}

func (c *SimpleCache) Del(key string) error {
	delete(c.dict, key)
	return nil
}

func (c *SimpleCache) LLen(key string) (int, error) {
	if v, ok := c.dict[key].([]string); ok {
		return len(v), nil
	}
	return 0, redigoerr.WrongType
}

func (c *SimpleCache) Lock() {
	c.internalLock.Lock()
}

func (c *SimpleCache) Unlock() {
	c.internalLock.Unlock()
}
