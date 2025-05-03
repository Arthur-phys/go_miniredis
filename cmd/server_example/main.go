package main

import (
	"log/slog"

	"github.com/Arthur-phys/redigo/pkg/core/caches"
	"github.com/Arthur-phys/redigo/pkg/server"
)

func main() {
	s, err := server.MakeServer(
		"127.0.0.1",
		8000,
		caches.NewSimpleCacheStore,
		server.NewWorkerInstantiator(),
		10240,
		2,
		15,
	)
	if err != nil {
		slog.Error("Fatal error occurred!", "", err)
		return
	}
	s.Run()
}
