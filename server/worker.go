package server

import (
	"log/slog"
	"miniredis/core/coreinterface"
	"miniredis/core/parser"
	rt "miniredis/tobytes"
	"net"
	"time"
)

type Worker struct {
	cacheStore             coreinterface.CacheStore
	parseInstantiator      func(c *net.Conn, maxBytesPerCallAllowed int) *parser.RESPParser
	connectionChannel      chan net.Conn
	timeout                uint
	id                     uint64
	maxBytesPerCallAllowed int
}

func NewWorkerInstantiator() func(
	cacheStore coreinterface.CacheStore,
	connectionChannel chan net.Conn,
	maxBytesPerCallAllowed int,
	timeout uint,
) Worker {
	var i uint64 = 0
	return func(CacheStore coreinterface.CacheStore, connectionChannel chan net.Conn, maxBytesPerCallAllowed int, timeout uint) Worker {
		i++
		return Worker{CacheStore, parser.NewRESPParser, connectionChannel, timeout, i, maxBytesPerCallAllowed}
	}
}

func (w *Worker) handleConnection(c *net.Conn) {
	defer (*c).Close()
	(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
	parser := w.parseInstantiator(c, w.maxBytesPerCallAllowed)

	for {
		finalResponse := []byte{}
		_, err := parser.Read()
		if err.Code != 0 {
			// Stopped any Conn error here, incluiding EOF, Broken Pipe, etc.
			return
		}
		commands, err := parser.ParseCommand()
		// if the buffer was exhausted, do not return an error
		if err.Code != 0 && err.Code != 3 && err.Code != 4 {
			// Command malformed, return immediately
			slog.Error("An error occurred while parsing the command", "ERROR", err,
				slog.Uint64("WORKER_ID", w.id),
				slog.String("CLIENT", (*c).RemoteAddr().String()),
			)
			_, err := (*c).Write(rt.ErrToBytes(err))
			if err != nil {
				slog.Error("An error occurred while sending error response to client", "ERROR", err,
					slog.Uint64("WORKER_ID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
			}
			return
		}
		// Proceed with evaluation
		for _, command := range commands {
			w.cacheStore.Lock()
			res, err := command(w.cacheStore)
			w.cacheStore.Unlock()
			if err.Code != 0 {
				slog.Error("An error occurred while executing client's command", "ERROR", err,
					slog.Uint64("WORKER_ID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				// Errors are delivered at the end for every command
				finalResponse = append(finalResponse, rt.ErrToBytes(err)...)
				continue
			}
			finalResponse = append(finalResponse, res...)
		}

		_, nerr := (*c).Write(finalResponse)
		if nerr != nil {
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
