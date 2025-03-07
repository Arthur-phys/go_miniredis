package core

import (
	"bufio"
	e "miniredis/error"
	"miniredis/server"
	"net"
	"strconv"
)

type RESPParser struct {
	stream Stream
}

func (r *RESPParser) ParseCommand() (f func(d server.CacheStore) ([]byte, error), err error) {
	firstByte, err := r.stream.TakeOne()
	if err != nil {
		return
	}
	if firstByte != '*' {
		return func(d server.CacheStore) ([]byte, error) { return []byte{}, nil }, e.Error{} // Change
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

func selectFunction(arr []string) (f func(d server.CacheStore) ([]byte, error), err error) {
	switch arr[0] {
	case "GET":
		if len(arr) != 2 {
			return // Change proper error handling
		}
		return func(d server.CacheStore) ([]byte, error) {
			if val, ok := d.Get(arr[1]); ok {
				return BlobStringToRESP(val), nil
			} else {
				return []byte{}, e.Error{} //Change
			}
		}, nil
	case "SET":
		if len(arr) != 3 {
			return // Change proper error handling
		}
		return func(d server.CacheStore) ([]byte, error) {
			err = d.Set(arr[1], arr[2])
			if err != nil {
				return []byte{}, err
			}
			return NullToRESP(), nil
		}, nil
	case "RPUSH":
		if len(arr) < 3 {
			return
		}
		return func(d server.CacheStore) ([]byte, error) {
			err = d.RPush(arr[1], arr[2:]...)
			if err != nil {
				return []byte{}, err //Propper error handling
			}
			return NullToRESP(), nil
		}, nil
	case "RPOP":
		if len(arr) != 2 {
			return
		}
		return func(d server.CacheStore) ([]byte, error) {
			val, err := d.RPop(arr[1])
			if err != nil {
				return []byte{}, err // Propper error handling
			}
			return BlobStringToRESP(val), nil
		}, nil
	case "LLEN":
		if len(arr) != 2 {
			return
		}
		return func(d server.CacheStore) ([]byte, error) {
			val, err := d.LLen(arr[1])
			if err != nil {
				return []byte{}, err // Propper error handling
			}
			return IntToRESP(val), nil
		}, nil
	default:
		return func(d server.CacheStore) ([]byte, error) {
			return ErrToRESP(e.Error{}), nil
		}, e.Error{}
	}
}

func (r *RESPParser) miniRedisBlobStringFromBytes() (s string, err error) {
	firstByte, err := r.stream.TakeOne()
	if err != nil {
		return
	}
	if firstByte != '$' {
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
