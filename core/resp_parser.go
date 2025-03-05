package core

import (
	"bufio"
	"fmt"
	e "miniredis/error"
	"miniredis/server"
	"net"
	"strconv"
)

type RESPParser struct {
	stream Stream
}

func (r *RESPParser) ParseCommand() (f func(d *map[string]string) ([]byte, error), err error) {
	firstByte, err := r.stream.TakeOne()
	if err != nil {
		return
	}
	if firstByte != '*' {
		return func(d *map[string]string) ([]byte, error) { return []byte{}, nil }, e.Error{} // Change
	}
	bytesRead, err := r.stream.ReadUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		return
	}
	i, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		return
	}
	arr := make([]string, i)
	for j := range arr {
		arr[j], err = r.miniRedisBlobStringFromBytes()
		if err != nil {
			return
		}
	}
	return selectFunction(arr)
}

func selectFunction(arr []string) (f func(d *map[string]string) ([]byte, error), err error) {
	if len(arr) < 2 {
		return
	}
	switch arr[0] {
	case "GET":
		return func(d *map[string]string) ([]byte, error) {
			if val, ok := (*d)[arr[1]]; ok {
				return fmt.Appendf([]byte{}, val), nil // Proper formatting must ben ensured here. Check ToRESP function
			} else {
				return []byte{}, e.Error{} //Change
			}
		}, nil
	case "SET":
	case "LPUSH":
	case "LPOP":
	case "LLEN":
	}
	return
}

func (r *RESPParser) miniRedisBlobStringFromBytes() (s string, err error) {
	firstByte, err := r.stream.TakeOne()
	if err != nil {
		return
	}
	if firstByte != '*' {
		return "", e.Error{} // Change
	}
	bytesRead, err := r.stream.ReadUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		return
	}
	long, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		return
	}
	blobString, _, err := r.stream.ReadNBytes(long)
	if err != nil {
		return
	}
	r.stream.Skip(2)
	return string(blobString), nil
}

func NewRESPParser(conn *net.Conn) server.Parser {
	return &RESPParser{Stream{bufio.NewReader(*conn)}}
}
