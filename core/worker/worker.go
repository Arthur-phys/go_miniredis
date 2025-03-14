package worker

import (
	"log/slog"
	"miniredis/server"
	"net"
)

type SimpleWorker struct {
	cacheStore        server.CacheStore
	parseInstantiator func(c *net.Conn) Parser
	connectionChannel chan net.Conn
	id                uint64
}

type Parser interface {
	ParseCommand() (func(d server.CacheStore) ([]byte, error), error)
}

func NewSimpleWorkerInstantiator(
	parseInstantiator func(c *net.Conn) Parser,
) func(
	cacheStore server.CacheStore,
	connectionChannel chan net.Conn,
	id uint64,
) server.Worker {
	return func(cacheStore server.CacheStore, connectionChannel chan net.Conn, id uint64) server.Worker {
		return &SimpleWorker{cacheStore, parseInstantiator, connectionChannel, id}
	}
}

func (w *SimpleWorker) handleConnection(c net.Conn) {
	defer c.Close()
	parser := w.parseInstantiator(&c)
	command, err := parser.ParseCommand()
	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while parsing the command", err)
		return
	}

	w.cacheStore.Lock()
	res, err := command(w.cacheStore)
	w.cacheStore.Unlock()

	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while returning a response to the client", err)
		return
	}
	_, err = c.Write(res)
	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while returning a response to the client", err)
		return
	}
}

func (w *SimpleWorker) Run() {
	slog.Debug("[MiniRedis]", slog.Uint64("Starting Worker with ID", w.id))
	go func() {
		for {
			slog.Debug("[MiniRedis]", slog.String("Status of routine:", "Waiting for connection"))
			incomingConnection := <-w.connectionChannel
			w.handleConnection(incomingConnection)
		}
	}()
}
