package client

import (
	"fmt"
	"io"
	e "miniredis/error"
	rt "miniredis/resptypes"
	"net"
)

type Sender interface {
	get(key string) ([]byte, func(s *rt.Stream) (string, e.Error))
	set(key string, value string) ([]byte, func(s *rt.Stream) e.Error)
	rPush(key string, args ...string) ([]byte, func(s *rt.Stream) e.Error)
	rPop(key string) ([]byte, func(s *rt.Stream) (string, e.Error))
	lLen(key string) ([]byte, func(s *rt.Stream) (int, e.Error))
	lPop(key string) ([]byte, func(s *rt.Stream) (string, e.Error))
	lPush(key string, args ...string) ([]byte, func(s *rt.Stream) e.Error)
	lIndex(key string, index int) ([]byte, func(s *rt.Stream) (string, e.Error))
}

type Client struct {
	conn   *net.Conn
	sender Sender
}

func NewClient(conn *net.Conn, sender Sender) Client {
	return Client{conn, sender}
}

func (client *Client) Get(key string) (string, e.Error) {
	bytes, checker := client.sender.get(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}

	return checker(&returnBytes)
}

func (client *Client) Set(key string, value string) e.Error {
	bytes, checker := client.sender.set(key, value)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(&returnBytes)
}

func (client *Client) RPush(key string, args ...string) e.Error {
	bytes, checker := client.sender.rPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(&returnBytes)
}

func (client *Client) RPop(key string) (string, e.Error) {
	bytes, checker := client.sender.rPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(&returnBytes)
}

func (client *Client) LLen(key string) (int, e.Error) {
	bytes, checker := client.sender.lLen(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return 0, e.Error{}
	}
	return checker(&returnBytes)
}

func (client *Client) LPop(key string) (string, e.Error) {
	bytes, checker := client.sender.lPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(&returnBytes)
}

func (client *Client) LPush(key string, args ...string) e.Error {
	bytes, checker := client.sender.lPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return err
	}
	return checker(&returnBytes)
}

func (client *Client) LIndex(key string, index int) (string, e.Error) {
	bytes, checker := client.sender.lIndex(key, index)
	returnBytes, err := client.sendAndRead(bytes)
	if err.Code != 0 {
		return "", err
	}
	return checker(&returnBytes)
}

func (client *Client) sendAndRead(bytes []byte) (rt.Stream, e.Error) {
	_, err := (*client.conn).Write(bytes)
	if err != nil {
		return rt.NewStream([]byte{}), e.Error{}
	}
	returnBytes := make([]byte, 10240)
	_, err = (*client.conn).Read(returnBytes)
	fmt.Printf("Came back from reading bytes: %v\r\n", string(returnBytes))
	if err != io.EOF && err != nil {
		return rt.NewStream([]byte{}), e.Error{}
	}
	return rt.NewStream(returnBytes), e.Error{}
}
