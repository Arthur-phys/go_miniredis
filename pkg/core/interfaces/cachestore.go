package interfaces

import (
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type CacheStore interface {
	Get(key string) (string, e.Error)
	Set(key string, value string) e.Error
	RPush(key string, args ...string) e.Error
	RPop(key string) (string, e.Error)
	LLen(key string) (int, e.Error)
	LPop(key string) (string, e.Error)
	LPush(key string, args ...string) e.Error
	LIndex(key string, index int) (string, e.Error)
	Lock()
	Unlock()
}
