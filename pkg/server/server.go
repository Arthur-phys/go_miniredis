// server package provides both the Server struct for initializing an instance of REDIGO
// and a  Configuration struct to pass commands to the creation of a server.
//
// It also has the worker implementation, but this is not accessible to the library's user.
//
// Take into consideration that there is no RESP handshake.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

// Server holds all information related to a server.
// Both accepts connections and orchestrates workers by initializing them and
// stopping them when signailed like so by the OS or user (Using Ctrl+C for example)
type Server struct {
	listener          net.Listener
	cacheStore        interfaces.CacheStore
	connections       chan net.Conn
	signals           chan os.Signal
	workerNotifiers   []chan struct{}
	shutdownWaiter    *sync.WaitGroup
	shutdownTolerance int64
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			// Whenever signailed to close the server, do so
			slog.Info("Listener closed")
			break
		} else if err != nil {
			// Continue trying to accept connections even if one fails
			slog.Error("An error occurred while accepting a new connection", "ERROR", err)
			continue
		}
		s.connections <- conn
	}
}

func (s *Server) Run() {
	// Ask to be notified when program is to be shutdown, disables go normal behaviour when Ctrl+C
	signal.Notify(s.signals, syscall.SIGINT, syscall.SIGTERM)

	// Delegate connection acceptance to another routine to listen for syscalls
	go s.accept()

	// Waiting for a signal to close from os
	<-s.signals
	slog.Info("Shutting down server, signailing workers", slog.Int64("SHUTDOWNTOLERANCE", s.shutdownTolerance))
	// Signailing every worker
	for i := range s.workerNotifiers {
		s.workerNotifiers[i] <- struct{}{}
	}
	// Signailing connection goroutine to stop
	s.listener.Close()
	// Closing connection channel, which will completely terminate workers after the grace period to attend connections
	close(s.connections)

	// Now wait for every worker to finish
	shutdownSignailer := make(chan struct{})
	go func() {
		defer close(shutdownSignailer)
		s.shutdownWaiter.Wait()
	}()
	// In case a worker takes more than the tolerance, we end the program anyway
	select {
	case <-shutdownSignailer:
		slog.Info("All workers closed, terminating server")
		// Give an extra 5 secs for workers to do stuff before being left behind
	case <-time.After(time.Duration(s.shutdownTolerance) + time.Second*5):
		slog.Error("Unable to close all workers, terminating server anyway")
	}
}

func New(serverConfig *Configuration) (*Server, e.Error) {

	// Configure global logger to use ip, port and REDIGO as values in log output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("[REDIGO]", "")
	logger = logger.With("IP", serverConfig.IpAddress)
	logger = logger.With("PORT", serverConfig.Port)
	slog.SetDefault(logger)
	slog.Info("Initializing Server")

	// keepalive via TCP probes is disabled, every connection checks it on its own
	listenerConfig := net.ListenConfig{KeepAlive: -1}
	listener, err := listenerConfig.Listen(context.Background(), "tcp", net.JoinHostPort(serverConfig.IpAddress, fmt.Sprintf("%d", serverConfig.Port)))
	if err != nil {
		redigoError := e.UnableToCreateServer
		redigoError.From = err
		return &Server{}, redigoError
	}
	slog.Debug("Listener created")

	connections := make(chan net.Conn)
	signals := make(chan os.Signal, 1)
	workerNotifiers := make([]chan struct{}, serverConfig.WorkerAmount)
	shutdownWaiter := &sync.WaitGroup{}
	cacheStore := serverConfig.CacheStoreInstantiator()

	// Creating workers and running them
	for i := range serverConfig.WorkerAmount {
		notifications := make(chan struct{}, 1)
		workerNotifiers[i] = notifications
		worker := worker{
			cacheStore:        cacheStore,
			connections:       connections,
			timeout:           serverConfig.KeepAlive,
			notifications:     notifications,
			id:                i,
			shutdownTolerance: serverConfig.ShutdownTolerance,
			parser:            respparser.New(nil, serverConfig.MessageSizeLimit),
			shutdownWaiter:    shutdownWaiter,
		}
		go worker.run()
	}

	// Creating server
	server := Server{
		listener:          listener,
		cacheStore:        cacheStore,
		connections:       connections,
		signals:           signals,
		workerNotifiers:   workerNotifiers,
		shutdownTolerance: serverConfig.ShutdownTolerance,
		shutdownWaiter:    shutdownWaiter,
	}
	return &server, e.Error{}
}

// Configuration is a helper struct to be more idiomatic when configuring a server.
// You can see it in action in the cmd/redigo_server/ command.
type Configuration struct {
	IpAddress              string
	Port                   uint16
	WorkerAmount           uint64
	KeepAlive              int64
	MessageSizeLimit       int
	ShutdownTolerance      int64
	CacheStoreInstantiator func() interfaces.CacheStore
}
