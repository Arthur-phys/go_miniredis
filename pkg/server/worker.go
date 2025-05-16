package server

import (
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	rt "github.com/Arthur-phys/redigo/pkg/core/tobytes"
)

type Worker struct {
	cacheStore             interfaces.CacheStore
	parseInstantiator      func(c *net.Conn, maxBytesPerCallAllowed int) *respparser.RESPParser
	connectionChannel      chan net.Conn
	timeout                int64
	id                     uint64
	maxBytesPerCallAllowed int
	workerChannel          chan int64
	shutdown               bool
	workerWaitGroup        *sync.WaitGroup
}

func NewWorkerInstantiator() func(
	cacheStore interfaces.CacheStore,
	connectionChannel chan net.Conn,
	maxBytesPerCallAllowed int,
	timeout int64,
	workerChannel chan int64,
	workerWaitGroup *sync.WaitGroup,
) Worker {
	var i uint64 = 0
	return func(CacheStore interfaces.CacheStore, connectionChannel chan net.Conn, maxBytesPerCallAllowed int, timeout int64, workerChannel chan int64, workerWaitgroup *sync.WaitGroup) Worker {
		i++
		return Worker{CacheStore, respparser.New, connectionChannel, timeout, i, maxBytesPerCallAllowed, workerChannel, false, workerWaitgroup}
	}
}

func (w *Worker) handleConnection(c *net.Conn) {
	defer (*c).Close()
	(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
	respparser := w.parseInstantiator(c, w.maxBytesPerCallAllowed)

	for {
		select {
		case fullTimeout := <-w.workerChannel:
			(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(fullTimeout)))
			slog.Debug("Starting shutdown for worker", slog.Uint64("WORKERID", w.id))
			w.shutdown = true

		default:
			finalResponse := []byte{}
			_, err := respparser.Read()
			if err.Code == 15 {
				// Stopped any Conn error here, incluiding EOF, Broken Pipe, etc.
				return
			} else if err.Code == 17 {
				// Too big of a command
				_, err := (*c).Write(rt.ErrToBytes(err))
				if err != nil {
					slog.Error("An error occurred while sending error response to client", "ERROR", err,
						slog.Uint64("WORKERID", w.id),
						slog.String("CLIENT", (*c).RemoteAddr().String()),
					)
				}
				continue
			}
			commands, err := respparser.ParseCommand()
			// if the buffer was exhausted, do not return an error
			if err.Code != 0 && err.Code != 3 && err.Code != 4 && err.Code != 8 {
				// Command malformed, return immediately
				slog.Error("An error occurred while parsing the command", "ERROR", err,
					slog.Uint64("WORKERID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				_, err := (*c).Write(rt.ErrToBytes(err))
				if err != nil {
					slog.Error("An error occurred while sending error response to client", "ERROR", err,
						slog.Uint64("WORKERID", w.id),
						slog.String("CLIENT", (*c).RemoteAddr().String()),
					)
				}
				return
			}
			for _, command := range commands {
				w.cacheStore.Lock()
				res, err := command(w.cacheStore)
				w.cacheStore.Unlock()
				if err.Code != 0 {
					slog.Error("An error occurred while executing client's command", "ERROR", err,
						slog.Uint64("WORKERID", w.id),
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
					slog.Uint64("WORKERID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				return
			}
			if !w.shutdown {
				(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
			}
		}
	}
}

func (w *Worker) Run() {
	w.workerWaitGroup.Add(1)
	defer w.workerWaitGroup.Done()
	slog.Info("Starting Worker", slog.Uint64("WORKERID", w.id))
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
