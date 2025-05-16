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
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type Server struct {
	listener          net.Listener
	cacheStore        interfaces.CacheStore
	connectionChannel chan net.Conn
	osSigs            chan os.Signal
	workerChannels    []chan int64
	workerWaitGroup   *sync.WaitGroup
	shutdownTolerance int64
}

func (s *Server) Accept() {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			slog.Debug("Listener closed")
			break
		}
		if err != nil {
			slog.Error("An error occurred while accepting a new connection", "ERROR", err)
			continue
		}
		s.connectionChannel <- conn
	}
}

func (s *Server) Run() {
	signal.Notify(s.osSigs, syscall.SIGINT, syscall.SIGTERM)

	// Delegate connection acceptance to another routine to listen for syscalls
	go s.Accept()

	// Waiting for a signal to close from os
	<-s.osSigs
	slog.Info("Shutting down server, signailing workers", slog.Int64("SHUTDOWNTOLERANCE", s.shutdownTolerance))
	// Signailing every worker
	for i := range s.workerChannels {
		s.workerChannels[i] <- s.shutdownTolerance
	}
	// Signailing connection goroutine to stop
	s.listener.Close()
	// Closing connection channel, which will completely terminate workers after the grace period to attend connections
	close(s.connectionChannel)

	// Now wait for every worker to finish
	waitGroupChannel := make(chan struct{})
	go func() {
		defer close(waitGroupChannel)
		s.workerWaitGroup.Wait()
	}()
	// In case a worker takes more than the tolerance, we end the program anyway
	select {
	case <-waitGroupChannel:
		slog.Info("All workers closed, terminating server")
	case <-time.After(time.Duration(s.shutdownTolerance) + time.Second*2):
		slog.Error("Unable to close all workers, terminating server anyway")
	}
}
func New(
	ipAddress string,
	port uint16,
	cacheStoreInstantiator func() interfaces.CacheStore,
	workerInstantiator func(c interfaces.CacheStore, jobs chan net.Conn, maxBytesPerCallAllowed int, timeout int64, workerChannel chan int64, workerWaitgroup *sync.WaitGroup) Worker,
	maxBytesPerCallAllowed int,
	workerNumber uint,
	keepAlive int64,
	shutdownTolerance int64,
) (Server, error) {
	var server Server

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("[REDIGO]", "")
	logger = logger.With("IP", ipAddress)
	logger = logger.With("PORT", port)
	slog.SetDefault(logger)
	slog.Info("Initializing Server")

	listenerConfig := net.ListenConfig{}
	listenerConfig.KeepAlive = time.Duration(keepAlive) * time.Second
	slog.Debug("KeepAliveConfig configuration set", slog.Int("KEAAPLIVE", int(keepAlive)))

	listener, err := listenerConfig.Listen(context.Background(), "tcp", fmt.Sprintf("%v:%v", ipAddress, port))
	if err != nil {
		miniredisError := e.Error{}
		miniredisError.From = err
		return Server{}, miniredisError
	}
	slog.Debug("Listener created")

	server.listener = listener
	server.cacheStore = cacheStoreInstantiator()
	server.connectionChannel = make(chan net.Conn)
	server.osSigs = make(chan os.Signal, 1)
	server.workerChannels = make([]chan int64, workerNumber)
	server.shutdownTolerance = shutdownTolerance
	server.workerWaitGroup = &sync.WaitGroup{}

	for i := range workerNumber {
		workerChannel := make(chan int64, 1)
		server.workerChannels[i] = workerChannel
		worker := workerInstantiator(server.cacheStore, server.connectionChannel, maxBytesPerCallAllowed, keepAlive, workerChannel, server.workerWaitGroup)
		worker.Run()
	}

	return server, nil
}
