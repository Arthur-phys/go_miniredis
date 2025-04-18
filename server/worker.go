package server

import (
	"log/slog"
	"miniredis/core/coreinterface"
	rt "miniredis/resptypes"
	"net"
	"time"
)

type Worker struct {
	cacheStore        coreinterface.CacheStore
	parser            coreinterface.Parser
	connectionChannel chan net.Conn
	timeout           uint
	id                uint64
}

func NewWorkerInstantiator(
	parseInstantiator func() coreinterface.Parser,
) func(
	cacheStore coreinterface.CacheStore,
	connectionChannel chan net.Conn,
	timeout uint,
) Worker {
	var i uint64 = 0
	return func(CacheStore coreinterface.CacheStore, connectionChannel chan net.Conn, timeout uint) Worker {
		i++
		return Worker{CacheStore, parseInstantiator(), connectionChannel, timeout, i}
	}
}

func (w *Worker) handleConnection(c *net.Conn) {
	defer (*c).Close()
	(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
	finalResponse := []byte{}
	stream := rt.NewStreamFromConnection(c)
	commands, newErr := w.parser.ParseCommand(&stream)
	if newErr.Code != 0 {
		slog.Error("An error occurred while parsing the command", "ERROR", newErr,
			slog.Uint64("WORKER_ID", w.id),
			slog.String("CLIENT", (*c).RemoteAddr().String()),
		)
		_, err := (*c).Write(rt.ErrToBytes(newErr))
		if err != nil {
			slog.Error("An error occurred while sending error response to client", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
		}
		return
	}

	for _, command := range commands {
		w.cacheStore.Lock()
		res, err := command(w.cacheStore)
		finalResponse = append(finalResponse, res...)
		w.cacheStore.Unlock()
		if err.Code != 0 {
			slog.Error("An error occurred while executing client's command", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
			_, err := (*c).Write(rt.ErrToBytes(newErr))
			if err != nil {
				slog.Error("An error occurred while sending error response to client", "ERROR", err,
					slog.Uint64("WORKER_ID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
			}
			return
		}
	}
	_, err := (*c).Write(finalResponse)
	if err != nil {
		slog.Error("An error occurred while returning a response to the client", "ERROR", err,
			slog.Uint64("WORKER_ID", w.id),
			slog.String("CLIENT", (*c).RemoteAddr().String()),
		)
		return
	}
	(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
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
