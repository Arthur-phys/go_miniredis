package main

import (
	"log/slog"
	"miniredis/core"
	"miniredis/server"
)

func main() {
	s, err := server.MakeServer(
		"127.0.0.1",
		8000,
		core.NewRESPParser,
		core.NewShardedCacheStore,
		map[string]uint{"maxGoRoutines": 10, "keepAlive": 10},
	)
	if err != nil {
		slog.Error("Fatal error occurred!", "", err)
		return
	}
	s.Run()
}
