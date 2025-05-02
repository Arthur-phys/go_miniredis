package server

import (
	"context"
	"fmt"
	"log/slog"
	"miniredis/core/coreinterface"
	e "miniredis/error"
	"net"
	"os"
	"time"
)

type Server struct {
	listener          net.Listener
	cacheStore        coreinterface.CacheStore
	connectionChannel chan net.Conn
}

func (s *Server) Accept() error {
	conn, err := s.listener.Accept()
	if err != nil {
		miniRedisError := e.Error{} // Change
		miniRedisError.From = err
		return miniRedisError
	}
	s.connectionChannel <- conn
	return nil
}

func (s *Server) Run() {
	for {
		err := s.Accept()
		if err != nil {
			slog.Error("An error occurred while accepting a new connection", "ERROR", err)
		}
	}
}
func MakeServer(
	ipAddress string,
	port uint16,
	cacheStoreInstantiator func() coreinterface.CacheStore,
	workerInstantiator func(c coreinterface.CacheStore, jobs chan net.Conn, maxBytesPerCallAllowed int, timeout int64) Worker,
	maxBytesPerCallAllowed int,
	workerNumber uint,
	keepAlive int64,
) (Server, error) {
	var server Server

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("[MiniRedis]", "")
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

	for range workerNumber {
		worker := workerInstantiator(server.cacheStore, server.connectionChannel, maxBytesPerCallAllowed, keepAlive)
		worker.Run()
	}

	return server, nil
}
