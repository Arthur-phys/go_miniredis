package server

import (
	"io"
	"log/slog"
	"miniredis/core/coreinterface"
	"net"
)

type Worker struct {
	cacheStore        coreinterface.CacheStore
	parseInstantiator func(c *net.Conn) coreinterface.Parser
	connectionChannel chan net.Conn
	id                uint64
}

func NewWorkerInstantiator(
	parseInstantiator func(c *net.Conn) coreinterface.Parser,
) func(
	cacheStore coreinterface.CacheStore,
	connectionChannel chan net.Conn,
	id uint64,
) Worker {
	return func(CacheStore coreinterface.CacheStore, connectionChannel chan net.Conn, id uint64) Worker {
		return Worker{CacheStore, parseInstantiator, connectionChannel, id}
	}
}

func (w *Worker) handleConnection(c *net.Conn) {
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
			slog.Error("An error occurred while returning a response to the client", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
			return
		}
		_, err = (*c).Write(res)
		if err != nil {
			slog.Error("An error occurred while returning a response to the client", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
			return
		}
	}
}

func (w *Worker) Run() {
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
