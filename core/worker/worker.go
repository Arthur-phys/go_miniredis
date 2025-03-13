package worker

import (
	"log/slog"
	"miniredis/server"
	"net"
)

type SimpleWorker struct {
	cacheStore        CacheStore
	parseInstantiator func(c *net.Conn) Parser
	connectionChannel chan net.Conn
	id                uint64
}

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

type Parser interface {
	ParseCommand() (func(d CacheStore) ([]byte, error), error)
}

func NewSimpleWorker(server *server.Server, id uint64) server.Worker {
	return &SimpleWorker{server.cacheStore, parseInstantiator, connectionChannel, id}
}

func (w *SimpleWorker) HandleConnection(c net.Conn) {
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
			incomingConnection := <-w.connectionChannel
			w.handleConnection(incomingConnection)
		}
	}()
}
