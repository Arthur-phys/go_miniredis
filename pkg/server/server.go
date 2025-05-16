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
func New(serverConfig *Configuration) (Server, error) {
	var server Server

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("[REDIGO]", "")
	logger = logger.With("IP", serverConfig.IpAddress)
	logger = logger.With("PORT", serverConfig.Port)
	slog.SetDefault(logger)
	slog.Info("Initializing Server")

	listenerConfig := net.ListenConfig{}
	listenerConfig.KeepAlive = time.Duration(serverConfig.KeepAlive) * time.Second
	slog.Debug("KeepAliveConfig configuration set", slog.Int("KEAAPLIVE", int(serverConfig.KeepAlive)))

	listener, err := listenerConfig.Listen(context.Background(), "tcp", fmt.Sprintf("%v:%v", serverConfig.IpAddress, serverConfig.Port))
	if err != nil {
		miniredisError := e.Error{}
		miniredisError.From = err
		return Server{}, miniredisError
	}
	slog.Debug("Listener created")

	server.listener = listener
	server.cacheStore = serverConfig.CacheStoreInstantiator()
	server.connectionChannel = make(chan net.Conn)
	server.osSigs = make(chan os.Signal, 1)
	server.workerChannels = make([]chan int64, serverConfig.WorkerSize)
	server.shutdownTolerance = serverConfig.ShutdownTolerance
	server.workerWaitGroup = &sync.WaitGroup{}

	workerInstantiator := NewWorkerInstantiator()
	for i := range serverConfig.WorkerSize {
		workerChannel := make(chan int64, 1)
		server.workerChannels[i] = workerChannel
		worker := workerInstantiator(server.cacheStore, server.connectionChannel, serverConfig.MessageSizeLimit, serverConfig.KeepAlive, workerChannel, server.workerWaitGroup)
		worker.Run()
	}

	return server, nil
}

type Configuration struct {
	IpAddress              string
	Port                   uint16
	WorkerSize             uint
	KeepAlive              int64
	MessageSizeLimit       int
	ShutdownTolerance      int64
	CacheStoreInstantiator func() interfaces.CacheStore
}
