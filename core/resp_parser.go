package core

import (
	"bufio"
	"net"
)

type RESPParser struct {
	stream Stream
}

func (r *RESPParser) ParseCommand()

func StreamFromConnection(conn *net.Conn) RESPParser {
	return RESPParser{Stream{bufio.NewReader(*conn)}}
}
