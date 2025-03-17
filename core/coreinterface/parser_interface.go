package coreinterface

type Parser interface {
	ParseCommand() (func(d CacheStore) ([]byte, error), error)
}
