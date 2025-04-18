package coreinterface

import (
	e "miniredis/error"
	"miniredis/resptypes"
)

type Parser interface {
	ParseCommand(stream *resptypes.Stream) ([]func(d CacheStore) ([]byte, e.Error), e.Error)
}
