// server package provides both the Server struct for initializing an instance of REDIGO
// and a  Configuration struct to pass commands to the creation of a server.
//
// It also has the worker implementation, but this is not accessible to the library's user.
//
// Take into consideration that there is no RESP handshake.
//
// An example of creating a server is provided here:
//
//	 import (
//	  "github.com/Arthur-phys/redigo/pkg/client"
//	  "net"
//	  "fmt"
//	 )
//
//	 fn main() {
//	  serverConfig := server.Configuration{
//	   IpAddress:              "127.0.0.1",
//	   Port:                   8000,
//	   WorkerAmount:           1,
//	   KeepAlive:              15,
//	   MessageSizeLimit:       10240,
//	   ShutdownTolerance:      5,
//	   CacheStoreInstantiator: caches.NewCache,
//	  }
//
//	  s, err := server.New(&serverConfig)
//
//	  if err != nil {
//		fmt.Printf("Fatal error occurred - %v\n", err)
//		return
//	  }
//	  s.Run()
//	}
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

	"github.com/Arthur-phys/redigo/pkg/core/cache"
	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	"github.com/Arthur-phys/redigo/pkg/redigoerr"
)

// Server holds all information related to a server.
// Both accepts connections and orchestrates workers by initializing them and
// stopping them when signailed like so by the OS or user (Using Ctrl+C for example)
type Server struct {
	listener          net.Listener
	cacheStore        *cache.Cache
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
	case <-time.After(time.Duration(s.shutdownTolerance+1) * time.Second):
		slog.Error("Unable to close all workers, terminating server anyway")
	}
}

func New(serverConfig *Configuration) (*Server, error) {

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
		redigoError := redigoerr.UnableToCreateServer
		redigoError.From = err
		return &Server{}, redigoError
	}
	slog.Debug("Listener created")

	connections := make(chan net.Conn)
	signals := make(chan os.Signal, 1)
	workerNotifiers := make([]chan struct{}, serverConfig.WorkerAmount)
	shutdownWaiter := &sync.WaitGroup{}
	cacheStore := cache.New()

	// Creating workers and running them
	for i := range serverConfig.WorkerAmount {
		notifications := make(chan struct{}, 1)
		workerNotifiers[i] = notifications
		worker := worker{
			cacheStore:     cacheStore,
			connections:    connections,
			timeout:        serverConfig.KeepAlive,
			notifications:  notifications,
			id:             i,
			parser:         respparser.New(nil, serverConfig.MessageSizeLimit),
			shutdownWaiter: shutdownWaiter,
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
	return &server, nil
}

// Configuration is a helper struct to be more idiomatic when configuring a server.
// You can see it in action in the cmd/redigo_server/ command.
type Configuration struct {
	IpAddress         string
	Port              uint16
	WorkerAmount      uint64
	KeepAlive         int64
	MessageSizeLimit  int
	ShutdownTolerance int64
}
