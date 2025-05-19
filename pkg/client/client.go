// client presents a simple implementation of how a client would be able to use
// the server.
//
// Take into consideration that many things in server side have been reused for this client and
// also that there is no RESP handshake, since it would not serve any purpose.
package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

type Client struct {
	conn   *net.Conn
	buffer *bufio.Reader
	p      *respparser.RESPParser
}

func New(conn *net.Conn) *Client {
	parser := respparser.New(conn, 10240)
	return &Client{conn, bufio.NewReader(*conn), parser}
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
	result, _, err := client.p.ParseBlobString()
	if err.Code == 5 {
		if err.ExtraContext["received"] == "_" {
			_, err := client.p.ParseNull()
			return "", err
		} else if err.ExtraContext["received"] == "-" {
			_, err = client.p.ParseError()
			return "", err
		}
	}
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
	_, err = client.p.ParseNull()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return err
	}
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
	_, err = client.p.ParseNull()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return err
	}
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
	result, _, err := client.p.ParseBlobString()
	if err.Code == 5 {
		if err.ExtraContext["received"] == "_" {
			_, err := client.p.ParseNull()
			return "", err
		} else if err.ExtraContext["received"] == "-" {
			_, err = client.p.ParseError()
			return "", err
		}
	}
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
	result, _, err := client.p.ParseUInt()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return 0, err
	}
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
	result, _, err := client.p.ParseBlobString()
	if err.Code == 5 {
		if err.ExtraContext["received"] == "_" {
			_, err := client.p.ParseNull()
			return "", err
		} else if err.ExtraContext["received"] == "-" {
			_, err = client.p.ParseError()
			return "", err
		}
	}
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
	_, err = client.p.ParseNull()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return err
	}
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
	result, _, err := client.p.ParseBlobString()
	if err.Code == 5 {
		if err.ExtraContext["received"] == "_" {
			_, err = client.p.ParseNull()
			return "", err
		} else if err.ExtraContext["received"] == "-" {
			_, err = client.p.ParseError()
			return "", err
		}
	}
	return result, err
}

func (client *Client) Del(key string) e.Error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nDEL\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return err
	}
	_, err = client.p.ParseNull()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return err
	}
	return err
}

func (client *Client) Ping() (string, e.Error) {
	finalBytes := fmt.Appendf([]byte{}, "*1\r\n$4\r\nPING\r\n")
	err := client.sendBytes(finalBytes)
	if err.Code != 0 {
		return "", err
	}
	_, err = client.p.Read()
	if err.Code != 0 {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if err.Code == 5 && err.ExtraContext["received"] == "-" {
		_, err = client.p.ParseError()
		return "", err
	}
	return result, e.Error{}
}

func (client *Client) sendBytes(b []byte) e.Error {
	_, err := (*client.conn).Write(b)
	if err != nil {
		redigoError := e.UnableToSendRequestToServer
		redigoError.From = err
		return redigoError
	}
	return e.Error{}
}
