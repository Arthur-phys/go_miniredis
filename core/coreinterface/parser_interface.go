package coreinterface

import (
	e "miniredis/error"
)

type Parser interface {
	ParseCommand(bytes []byte) ([]func(d CacheStore) ([]byte, e.Error), e.Error)
}
