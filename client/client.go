package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"miniredis/core/parser"
	e "miniredis/error"
	"net"
)

type Client struct {
	conn   *net.Conn
	buffer *bufio.Reader
	p      *parser.RESPParser
}

func NewClient(conn *net.Conn) Client {
	return Client{conn, bufio.NewReader(*conn), parser.NewRESPParser(conn)}
}

func (client *Client) Get(key string) (string, e.Error) {
	bytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%v\r\n", len(key), key))
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return client.p.BlobStringFromBytes(returnBytes)
}

func (client *Client) Set(key string, value string) e.Error {
	bytes, checker := client.sender.set(key, value)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) RPush(key string, args ...string) e.Error {
	bytes, checker := client.sender.rPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) RPop(key string) (string, e.Error) {
	bytes, checker := client.sender.rPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) LLen(key string) (int, e.Error) {
	bytes, checker := client.sender.lLen(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return 0, e.Error{}
	}
	return checker(returnBytes)
}

func (client *Client) LPop(key string) (string, e.Error) {
	bytes, checker := client.sender.lPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) LPush(key string, args ...string) e.Error {
	bytes, checker := client.sender.lPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) LIndex(key string, index int) (string, e.Error) {
	bytes, checker := client.sender.lIndex(key, index)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) sendAndRead(b []byte) (*bufio.Reader, e.Error) {
	_, err := (*client.conn).Write(b)
	if err != nil {
		return bufio.NewReader(bytes.NewReader([]byte{})), e.Error{}
	}
	returnBytes := make([]byte, 1024)
	n, err := (*client.buffer).Read(returnBytes)
	if err != io.EOF && err != nil {
		return bufio.NewReader(bytes.NewReader([]byte{})), e.Error{}
	}
	return bufio.NewReader(bytes.NewReader(returnBytes[:n])), e.Error{}
}

func (ss *SimpleSender) set(key string, value string) ([]byte, func(s *bufio.Reader) e.Error) {
	return fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(value), value)), rt.ErrorFromBytes
}

func (ss *SimpleSender) rPush(key string, args ...string) ([]byte, func(s *bufio.Reader) e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nRPUSH", len(args)+1))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	return finalBytes, rt.ErrorFromBytes
}

func (ss *SimpleSender) rPop(key string) ([]byte, func(s *bufio.Reader) (string, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nRPOP\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.BlobStringFromBytes
}

func (ss *SimpleSender) lLen(key string) ([]byte, func(s *bufio.Reader) (int, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLLEN\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.UIntFromBytes

}

func (ss *SimpleSender) lPop(key string) ([]byte, func(s *bufio.Reader) (string, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLPOP\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.BlobStringFromBytes
}

func (ss *SimpleSender) lPush(key string, args ...string) ([]byte, func(s *bufio.Reader) e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nLPUSH", len(args)+1))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	return finalBytes, rt.ErrorFromBytes
}

func (ss *SimpleSender) lIndex(key string, index int) ([]byte, func(s *bufio.Reader) (string, e.Error)) {
	return fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$6\r\nLINDEX\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(fmt.Sprintf("%v", index)), index)), rt.BlobStringFromBytes
}
