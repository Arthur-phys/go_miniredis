package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/core/parser"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type Client struct {
	conn   *net.Conn
	buffer *bufio.Reader
	p      *parser.RESPParser
}

func NewClient(conn *net.Conn) Client {
	return Client{conn, bufio.NewReader(*conn), parser.NewRESPParser(conn, 10240)}
}

func (client *Client) Get(key string) (string, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return "", err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return "", err
	}
	result, _, err := client.p.BlobStringFromBytes()
	return result, err
}

func (client *Client) Set(key string, value string) e.Error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(value), value))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return err
	}
	_, err = client.p.NullFromBytes()
	return err
}

func (client *Client) RPush(key string, args ...string) e.Error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nRPUSH\r\n$%d\r\n%v", len(args)+2, len(key), key))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n$%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return err
	}
	_, err = client.p.NullFromBytes()
	return err
}

func (client *Client) RPop(key string) (string, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nRPOP\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return "", err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return "", err
	}
	result, _, err := client.p.BlobStringFromBytes()
	return result, err
}

func (client *Client) LLen(key string) (int, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLLEN\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return 0, err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return 0, err
	}
	result, _, err := client.p.UIntFromBytes()
	return result, err
}

func (client *Client) LPop(key string) (string, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLPOP\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return "", err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return "", err
	}
	result, _, err := client.p.BlobStringFromBytes()
	return result, err
}

func (client *Client) LPush(key string, args ...string) e.Error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nLPUSH\r\n$%d\r\n%v", len(args)+2, len(key), key))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n$%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return err
	}
	_, err = client.p.NullFromBytes()
	return err
}

func (client *Client) LIndex(key string, index int) (string, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$6\r\nLINDEX\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(fmt.Sprintf("%v", index)), index))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return "", err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return "", err
	}
	result, _, err := client.p.BlobStringFromBytes()
	return result, err
}

func (client *Client) sendBytes(b []byte) e.Error {
	_, err := (*client.conn).Write(b)
	if err != nil {
		newErr := e.UnableToSendRequestToServer
		newErr.From = err
		return newErr
	}
	return e.Error{}
}
