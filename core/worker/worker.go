package worker

import (
	"io"
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

func (w *SimpleWorker) handleConnection(c *net.Conn) {
	defer (*c).Close()
	parser := w.parseInstantiator(c)
	for {
		command, err := parser.ParseCommand()
		if err == io.EOF {
			slog.Debug("Finished attending connection",
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
			return
		} else if err != nil {
			slog.Error("An error occurred while parsing the command", "ERROR", err)
			return
		}

		w.cacheStore.Lock()
		res, err := command(w.cacheStore)
		w.cacheStore.Unlock()

		if err != nil {
			slog.Error("An error occurred while returning a response to the client", "ERROR", err)
			return
		}
		_, err = (*c).Write(res)
		if err != nil {
			slog.Error("An error occurred while returning a response to the client", "ERROR", err)
			return
		}
	}
}

func (w *SimpleWorker) Run() {
	slog.Debug("Starting Worker", slog.Uint64("WORKER_ID", w.id))
	go func() {
		for {
			if incomingConnection, ok := <-w.connectionChannel; ok {
				w.handleConnection(&incomingConnection)
			} else {
				break
			}
		}
	}()
}
