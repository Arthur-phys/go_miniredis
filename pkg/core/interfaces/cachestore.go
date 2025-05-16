// interfaces provides a single interface to sepparate the cache from it's implementation
package interfaces

import (
	e "github.com/Arthur-phys/redigo/pkg/error"
)

// CacheStore is the main interface to be able to
// change stores.
//
// It's primary purpose is to be able to activate or deactivate the sharding
type CacheStore interface {
	Get(key string) (string, e.Error)
	Set(key string, value string) e.Error
	RPush(key string, args ...string) e.Error
	RPop(key string) (string, e.Error)
	LLen(key string) (int, e.Error)
	LPop(key string) (string, e.Error)
	LPush(key string, args ...string) e.Error
	LIndex(key string, index int) (string, e.Error)
	Del(key string) e.Error
	Lock()
	Unlock()
}
