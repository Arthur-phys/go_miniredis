package server

import (
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	"github.com/Arthur-phys/redigo/pkg/core/tobytes"
)

// worker accepts new tcp connections and responds to clients
// by parsing their commands.
type worker struct {
	cacheStore        interfaces.CacheStore
	parseInstantiator func(c *net.Conn, messageSizeLimit int) *respparser.RESPParser
	connections       chan net.Conn
	timeout           int64
	id                uint64
	messageSizeLimit  int
	notifications     chan int64
	shutdown          bool
	shutdownWaiter    *sync.WaitGroup
}

// handleConnection answer a single client until the connection closes or a timeout happens
func (w *worker) handleConnection(c *net.Conn) {
	// Never forget to close the connection!
	defer (*c).Close()
	// Setting max deadline for reading or writing
	(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
	respparser := w.parseInstantiator(c, w.messageSizeLimit)

	// A worker sticks with a connection until it closes, therefore just one worker attends a given connection
	for {
		select {
		// When signailed to stop, give the connection a last chance to be read and receive an answer
		case fullTimeout := <-w.notifications:
			(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(fullTimeout)))
			slog.Debug("Starting shutdown for worker, finishing any active connections", slog.Uint64("WORKERID", w.id))
			w.shutdown = true

		default:
			finalResponse := []byte{}
			_, err := respparser.Read()
			if err.Code == 15 {
				// Stopped any Conn error here, incluiding EOF, Broken Pipe, etc.
				slog.Debug("The connection was closed", "REASON", err.From,
					slog.Uint64("WORKERID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()))
				return
			} else if err.Code == 17 {
				// Too big of a command
				if _, err := (*c).Write(tobytes.Err(err)); err != nil {
					slog.Error("An error occurred while sending error response to client", "ERROR", err,
						slog.Uint64("WORKERID", w.id),
						slog.String("CLIENT", (*c).RemoteAddr().String()),
					)
				}
				continue
			}

			// Was able to read, now parse commands
			commands, err := respparser.ParseCommand()
			// If the buffer was exhausted, do not return an error, which is true for cases 0,3,4 & 8
			if err.Code != 0 && err.Code != 3 && err.Code != 4 && err.Code != 8 {
				// Command malformed, return immediately
				slog.Error("An error occurred while parsing the command", "ERROR", err,
					slog.Uint64("WORKERID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				_, err := (*c).Write(tobytes.Err(err))
				if err != nil {
					slog.Error("An error occurred while sending error response to client", "ERROR", err,
						slog.Uint64("WORKERID", w.id),
						slog.String("CLIENT", (*c).RemoteAddr().String()),
					)
				}
				return
			}

			// Interpret & evaluate commands
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
					finalResponse = append(finalResponse, tobytes.Err(err)...)
					continue
				}
				finalResponse = append(finalResponse, res...)
			}

			// Return all responses at once
			_, nerr := (*c).Write(finalResponse)
			if nerr != nil {
				slog.Error("An error occurred while returning a response to the client", "ERROR", err,
					slog.Uint64("WORKERID", w.id),
					slog.String("CLIENT", (*c).RemoteAddr().String()),
				)
				return
			}

			// Restart timer if not signailed to stop
			if w.shutdown {
				continue
			}
			(*c).SetDeadline(time.Now().Add(time.Second * time.Duration(w.timeout)))
		}
	}
}

// run is the main process of a worker
func (w *worker) run() {
	// Allow the server to wait on this worker for some time
	w.shutdownWaiter.Add(1)
	defer w.shutdownWaiter.Done()

	// Accept connections until the connection channel is closed
	slog.Info("Starting worker", slog.Uint64("WORKERID", w.id))
	for {
		if conn, ok := <-w.connections; ok {
			w.handleConnection(&conn)
		} else {
			break
		}
	}
}
