package coreinterface

type Parser interface {
	ParseCommand(buffer []byte) ([]func(d CacheStore) ([]byte, error), error)
}
