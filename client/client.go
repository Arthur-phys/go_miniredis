package client

import "net"

type Sender interface {
	get(key string) ([]byte, func([]byte) (string, error))
	set(key string, value string) ([]byte, func([]byte) (string, error))
	rPush(key string, args ...string) ([]byte, func([]byte) (string, error))
	rPop(key string) ([]byte, func([]byte) (string, error))
	lLen(key string) ([]byte, func([]byte) (string, error))
	lPop(key string) ([]byte, func([]byte) (string, error))
	lPush(key string, args ...string) ([]byte, func([]byte) (string, error))
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
	_, err := (*client.conn).Write(bytes)
	if err != nil {
		return "", err
	}
	returnBytes := make([]byte, 1024) // Change to a for with a 10kb limit!
	_, err = (*client.conn).Read(returnBytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) Set(key string, value string) (string, error) {
	bytes, checker := client.sender.set(key, value)
	_, err := (*client.conn).Write(bytes)
	if err != nil {
		return "", err
	}
	returnBytes := make([]byte, 1024) // Change to a for with a 10kb limit!
	_, err = (*client.conn).Read(returnBytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}

func (client *Client) RPush(key string, args ...string) (string, error) {
	bytes, checker := client.sender.rPush(key, args...)
	_, err := (*client.conn).Write(bytes)
	if err != nil {
		return "", err
	}
	returnBytes := make([]byte, 1024) // Change to a for with a 10kb limit!
	_, err = (*client.conn).Read(returnBytes)
	if err != nil {
		return "", err
	}
	return checker(returnBytes)
}
