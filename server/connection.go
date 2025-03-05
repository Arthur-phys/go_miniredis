package server

import (
	"log/slog"
	"net"
	"sync"
)

type Parser interface {
	ParseCommand() (func(d *map[string]string) ([]byte, error), error)
}

type Connection struct {
	conn               *net.Conn
	dictUsedLock       *sync.Mutex
	currentThreadsLock *sync.Mutex
	dict               *map[string]string
	currentGoRoutines  *uint16
	parser             Parser
}

func (c *Connection) Answer() {
	var err error
	defer func() {
		c.currentThreadsLock.Lock()
		*c.currentGoRoutines -= 1
		c.currentThreadsLock.Unlock()
		(*c.conn).Close()
	}()
	slog.Debug("[MiniRedis]", slog.Any("Answering connection", (*c.conn).RemoteAddr()))
	command, err := c.parser.ParseCommand()
	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while parsing the command", err)
		return
	}

	c.dictUsedLock.Lock()
	res, err := command(c.dict)
	c.dictUsedLock.Unlock()
	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while returning a response to the client", err)
		return
	}
	_, err = (*c.conn).Write(res)
	if err != nil {
		slog.Error("[MiniRedis]", "An error occurred while returning a response to the client", err)
		return
	}
}
