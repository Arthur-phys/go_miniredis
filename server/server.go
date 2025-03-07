package server

import (
	"context"
	"fmt"
	"log/slog"
	e "miniredis/error"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	listener           net.Listener
	maxGoRoutines      uint16
	currentGoRoutines  uint16
	cacheStore         CacheStore
	currentThreadsLock *sync.Mutex
	parserInstantiator func(c *net.Conn) Parser
}

type CacheStore interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	RPush(key string, args ...string) error
	RPop(key string) (string, error)
	LLen(key string) (int, error)
	Lock()
	Unlock()
}

func (s *Server) Accept() (Connection, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		miniRedisError := e.Error{} // Change
		miniRedisError.From = err
		return Connection{}, miniRedisError
	}
	s.currentThreadsLock.Lock()
	s.currentGoRoutines += 1
	slog.Debug("[MiniRedis]", slog.Int("Current GoRoutines", int(s.currentGoRoutines)))
	if s.maxGoRoutines < s.currentGoRoutines {
		conn.Close()
		return Connection{}, e.Error{} // Change
	}
	s.currentThreadsLock.Unlock()
	return s.newConnection(&conn), nil
}

func (s *Server) Run() {
	for {
		conn, err := s.Accept()
		if err != nil {
			slog.Error("[MiniRedis]", "Error", err)
		}
		go conn.Answer()
	}
}

func MakeServer(ipAddress string, port uint16, parserInstantiator func(c *net.Conn) Parser, cacheStoreInstantiator func() CacheStore, options map[string]uint) (Server, error) {
	var server Server

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With("IP", ipAddress)
	logger = logger.With("PORT", port)
	slog.SetDefault(logger)
	slog.Info("[MiniRedis]", slog.String("Initializing Server", ""))

	listenerConfig := net.ListenConfig{}
	if keepAlive, ok := options["KeepAlive"]; ok {
		listenerConfig.KeepAlive = time.Duration(keepAlive) * time.Second
		// listenerConfig.KeepAliveConfig.Enable = true
		slog.Debug("[MiniRedis]", slog.Int("KeepAliveConfig configuration set to seconds", int(keepAlive)))
	}

	listener, err := listenerConfig.Listen(context.Background(), "tcp", fmt.Sprintf("%v:%v", ipAddress, port))
	if err != nil {
		miniredisError := e.Error{}
		miniredisError.From = err
		return Server{}, miniredisError
	}
	slog.Debug("[MiniRedis]", slog.String("Listener created", ""))

	if maxGoRoutines, ok := options["maxGoRoutines"]; ok {
		server.maxGoRoutines = uint16(maxGoRoutines)
	} else {
		server.maxGoRoutines = 1024
	}
	slog.Debug("[MiniRedis]", slog.Int("MaxGoRoutines set to", int(server.maxGoRoutines)))

	server.listener = listener
	server.currentThreadsLock = &sync.Mutex{}
	server.cacheStore = cacheStoreInstantiator()
	server.parserInstantiator = parserInstantiator

	return server, nil
}

func (s *Server) newConnection(conn *net.Conn) Connection {
	connection := Connection{}
	connection.conn = conn
	connection.currentGoRoutines = &s.currentGoRoutines
	connection.currentThreadsLock = s.currentThreadsLock
	connection.cacheStore = s.cacheStore
	connection.parser = s.parserInstantiator(conn)

	return connection
}
