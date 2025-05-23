package cache

import (
	"fmt"
	"sync"

	"container/list"

	"github.com/Arthur-phys/redigo/pkg/redigoerr"
)

type Cache struct {
	internalLock sync.Mutex
	dict         map[string]any
}

func New() *Cache {
	return &Cache{
		sync.Mutex{},
		make(map[string]any),
	}
}

func (c *Cache) Get(key string) (string, error) {
	v, ok := c.dict[key]
	if !ok {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
	if v, ok := v.(string); ok {
		return v, nil
	}
	return "", redigoerr.WrongType
}

func (c *Cache) Set(key string, value string) error {
	c.dict[key] = value
	return nil
}

func (c *Cache) RPush(key string, args ...string) error {
	v, ok := c.dict[key]
	if !ok {
		l := list.New()
		for _, arg := range args {
			l.PushBack(arg)
		}
		c.dict[key] = l
		return nil
	}
	vAsList, ok := v.(*list.List)
	if !ok {
		return redigoerr.WrongType
	}
	for _, arg := range args {
		vAsList.PushBack(arg)
	}
	return nil
}

func (c *Cache) RPop(key string) (string, error) {
	v, ok := c.dict[key]
	if !ok {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
	vAsList, ok := v.(*list.List)
	if !ok {
		return "", redigoerr.WrongType
	}
	x := vAsList.Back().Value.(string)
	vAsList.Remove(vAsList.Back())
	if vAsList.Len() == 0 {
		delete(c.dict, key)
	}
	return x, nil
}

func (c *Cache) LPush(key string, args ...string) error {
	v, ok := c.dict[key]
	if !ok {
		l := list.New()
		for _, arg := range args {
			l.PushFront(arg)
		}
		c.dict[key] = l
		return nil
	}
	vAsList, ok := v.(*list.List)
	if !ok {
		return redigoerr.WrongType
	}
	for _, arg := range args {
		vAsList.PushFront(arg)
	}
	return nil
}

func (c *Cache) LPop(key string) (string, error) {
	v, ok := c.dict[key]
	if !ok {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err
	}
	if v, ok := v.(*list.List); ok {
		x := v.Front().Value.(string)
		v.Remove(v.Front())
		if v.Len() == 0 {
			delete(c.dict, key)
		}
		return x, nil
	}
	return "", redigoerr.WrongType
}

func (c *Cache) LIndex(key string, index int) (string, error) {
	v, ok := c.dict[key]
	if !ok {
		err := redigoerr.KeyNotFoundInDictionary
		err.ExtraContext = map[string]string{"key": key}
		return "", err

	}
	vAsList, ok := v.(*list.List)
	if !ok {
		return "", redigoerr.WrongType
	}
	if vAsList.Len() <= index || index < 0 {
		err := redigoerr.IndexOutOfRangeErr
		err.ExtraContext = map[string]string{"index": fmt.Sprintf("%d", index)}
		return "", err
	}
	front := vAsList.Front()
	for range index {
		front = front.Next()
	}
	return front.Value.(string), nil
}

func (c *Cache) LLen(key string) (int, error) {
	if v, ok := c.dict[key].(*list.List); ok {
		return v.Len(), nil
	}
	return 0, redigoerr.WrongType
}

func (c *Cache) Del(key string) error {
	delete(c.dict, key)
	return nil
}

func (c *Cache) Lock() {
	c.internalLock.Lock()
}

func (c *Cache) Unlock() {
	c.internalLock.Unlock()
}
