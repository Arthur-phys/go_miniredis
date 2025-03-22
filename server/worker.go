package server

import (
	"io"
	"log/slog"
	"miniredis/core/coreinterface"
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
	buffer := make([]byte, 10240)

	for {
		n, err := (*c).Read(buffer)
		if err != nil {
			if err == io.EOF && n == 0 {
				slog.Debug("Finished attending connection",
					slog.Uint64("WORKER_ID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				return
			} else if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					slog.Error("Connection timeout",
						slog.Uint64("WORKER_ID", w.id),
						slog.String("CLIENT", (*c).RemoteAddr().String()),
					)
				}
				return
			} else if err != io.EOF {
				slog.Error("Unknown error occurred!", "ERROR", err,
					slog.Uint64("WORKER_ID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				return
			}
		}

		command, err := w.parser.ParseCommand(buffer[:n])
		if err != nil {
			slog.Error("An error occurred while parsing the command", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
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
		(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))

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
