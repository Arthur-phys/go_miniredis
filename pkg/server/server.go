package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
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
	shutdownTolerance int64
}

func (s *Server) Run() {
	signal.Notify(s.osSigs)
out:
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error("An error occurred while accepting a new connection", "ERROR", err)
		}
		select {
		case s.connectionChannel <- conn:
			continue
		case sig := <-s.osSigs:
			fmt.Println("Received", sig)
			slog.Info("Shutting down server, signailing workers", slog.Int64("SHUTDOWNTOLERANCE", s.shutdownTolerance))
			for i := range s.workerChannels {
				s.workerChannels[i] <- s.shutdownTolerance
			}
			close(s.connectionChannel)
			break out
		}
	}

	s.listener.Close()
}
func New(
	ipAddress string,
	port uint16,
	cacheStoreInstantiator func() interfaces.CacheStore,
	workerInstantiator func(c interfaces.CacheStore, jobs chan net.Conn, maxBytesPerCallAllowed int, timeout int64, workerChannel chan int64) Worker,
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

	for i := range workerNumber {
		workerChannel := make(chan int64, 1)
		server.workerChannels[i] = workerChannel
		worker := workerInstantiator(server.cacheStore, server.connectionChannel, maxBytesPerCallAllowed, keepAlive, workerChannel)
		worker.Run()
	}

	return server, nil
}
