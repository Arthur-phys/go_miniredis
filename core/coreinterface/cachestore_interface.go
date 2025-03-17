package coreinterface

type CacheStore interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	RPush(key string, args ...string) error
	RPop(key string) (string, error)
	LLen(key string) (int, error)
	LPop(key string) (string, error)
	LPush(key string, args ...string) error
	LIndex(key string, index int) (string, bool)
	Lock()
	Unlock()
}
