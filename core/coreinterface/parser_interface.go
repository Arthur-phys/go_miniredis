package coreinterface

import (
	e "miniredis/error"
)

type Parser interface {
	ParseCommand() ([]func(d CacheStore) ([]byte, e.Error), e.Error)
}
