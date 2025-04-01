package client

import (
	"io"
	"net"
)

type Sender interface {
	get(key string) ([]byte, func([]byte) (string, error))
	set(key string, value string) ([]byte, func([]byte) error)
	rPush(key string, args ...string) ([]byte, func([]byte) error)
	rPop(key string) ([]byte, func([]byte) (string, error))
	lLen(key string) ([]byte, func([]byte) (uint, error))
	lPop(key string) ([]byte, func([]byte) (string, error))
	lPush(key string, args ...string) ([]byte, func([]byte) error)
	lIndex(key string, index int) ([]byte, func([]byte) (string, error))
}

type Client struct {
	conn   *net.Conn
	sender Sender
}

func NewClient(conn *net.Conn, sender Sender) Client {
	return Client{conn, sender}
}

func (client *Client) Get(key string) (string, error) {
	bytes, checker := client.sender.get(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) Set(key string, value string) error {
	bytes, checker := client.sender.set(key, value)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) RPush(key string, args ...string) error {
	bytes, checker := client.sender.rPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) RPop(key string) (string, error) {
	bytes, checker := client.sender.rPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) LLen(key string) (uint, error) {
	bytes, checker := client.sender.lLen(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return 0, err
	}
	return checker(returnBytes)
}

func (client *Client) LPop(key string) (string, error) {
	bytes, checker := client.sender.lPop(key)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) LPush(key string, args ...string) error {
	bytes, checker := client.sender.lPush(key, args...)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return err
	}
	return checker(returnBytes)
}

func (client *Client) LIndex(key string, index int) (string, error) {
	bytes, checker := client.sender.lIndex(key, index)
	returnBytes, err := client.sendAndRead(bytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) sendAndRead(bytes []byte) ([]byte, error) {
	returnBytes := []byte{}
	_, err := (*client.conn).Write(bytes)
	if err != nil {
		return []byte{}, err
	}
	totalBytesRead := 0
	for {
		tmpBytes := make([]byte, 1024)
		n, err := (*client.conn).Read(tmpBytes)
		totalBytesRead += n
		if err == io.EOF && n == 0 {
			break
		} else if (err == io.EOF && n != 0) || err == nil {
			returnBytes = append(returnBytes, tmpBytes...)
		} else {
			return []byte{}, err
		}
		if totalBytesRead > 10240 {
			return []byte{}, err // Custom error must replace this
		}
	}
	return returnBytes, nil
}
