// interfaces provides a single interface to sepparate the cache from it's implementation
package interfaces

// CacheStore is the main interface to be able to
// change stores.
//
// It's primary purpose is to be able to activate or deactivate the sharding
type CacheStore interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	RPush(key string, args ...string) error
	RPop(key string) (string, error)
	LLen(key string) (int, error)
	LPop(key string) (string, error)
	LPush(key string, args ...string) error
	LIndex(key string, index int) (string, error)
	Del(key string) error
	Lock()
	Unlock()
}
