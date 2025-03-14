package server

import (
	"context"
	"fmt"
	"log/slog"
	e "miniredis/error"
	"net"
	"os"
	"time"
)

type Server struct {
	listener          net.Listener
	cacheStore        CacheStore
	connectionChannel chan net.Conn
}

type Worker interface {
	Run()
}

type CacheStore interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	RPush(key string, args ...string) error
	RPop(key string) (string, error)
	LLen(key string) (int, error)
	LPop(key string) (string, error)
	LPush(key string, args ...string) error
	LIndex(key string, index int) (string, bool)
	Lock()
	Unlock()
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
			slog.Error("[MiniRedis]", "Error", err)
		}
	}
}
func MakeServer(
	ipAddress string,
	port uint16,
	cacheStoreInstantiator func() CacheStore,
	workerInstantiator func(c CacheStore, jobs chan net.Conn, id uint64) Worker,
	workerNumber uint,
) (Server, error) {
	var server Server

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("IP", ipAddress)
	logger = logger.With("PORT", port)
	slog.SetDefault(logger)
	slog.Info("[MiniRedis]", slog.String("Initializing Server", ""))

	listenerConfig := net.ListenConfig{}
	listenerConfig.KeepAlive = time.Duration(10) * time.Second
	slog.Debug("[MiniRedis]", slog.Int("KeepAliveConfig configuration set to seconds", 10))

	listener, err := listenerConfig.Listen(context.Background(), "tcp", fmt.Sprintf("%v:%v", ipAddress, port))
	if err != nil {
		miniredisError := e.Error{}
		miniredisError.From = err
		return Server{}, miniredisError
	}
	slog.Debug("[MiniRedis]", slog.String("Listener created", ""))

	server.listener = listener
	server.cacheStore = cacheStoreInstantiator()
	server.connectionChannel = make(chan net.Conn)

	for i := range workerNumber {
		worker := workerInstantiator(server.cacheStore, server.connectionChannel, uint64(i))
		worker.Run()
	}

	return server, nil
}
