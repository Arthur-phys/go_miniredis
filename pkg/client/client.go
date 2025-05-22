// client presents a simple implementation of how a client would be able to use
// the server.
//
// Take into consideration that many things in server side have been reused for this client and
// also that there is no RESP handshake, since it would not serve any purpose.
//
// A client example is provided here:
//
//	import (
//		"github.com/Arthur-phys/redigo/pkg/client"
//		"net"
//		"fmt"
//	)
//
//	func main() {
//		conn, err := net.Dial("tcp", "127.0.0.1:8000")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %v", err)
//		}
//		c := client.New(&conn)
//		err = c.Set("Arturo", "26")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		res, err := c.Get("Arturo")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		fmt.Printf("I got this! %v\n", res)
//		err = c.Set("Gene", "Le gustan los gatos")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		res, err = c.Get("Gene")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		fmt.Printf("I got this! %v\n", res)
//		err = c.LPush("Gatos", "Niji", "Anubis", "Pingüica", "Don Bigos")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		res, err = c.LPop("Gatos")
//		if err != nil {
//			fmt.Printf("Fatal error occurred! %e\n", err)
//		}
//		fmt.Printf("I got this! %v\n", res)
//		conn.Close()
//	}
package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/core/respparser"
	"github.com/Arthur-phys/redigo/pkg/redigoerr"
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

func (client *Client) Get(key string) (string, error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return "", err
	}
	_, err = client.p.Read()
	if err != nil {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if bytesDiffer(err) {
		if isRESPNull(err) {
			_, err := client.p.ParseNull()
			return "", err
		} else if isRESPError(err) {
			_, err = client.p.ParseError()
			return "", err
		}
	}
	return result, err
}

func (client *Client) Set(key string, value string) error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(value), value))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return err
	}
	_, err = client.p.Read()
	if err != nil {
		return err
	}
	_, err = client.p.ParseNull()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return err
	}
	return err
}

func (client *Client) RPush(key string, args ...string) error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nRPUSH\r\n$%d\r\n%v", len(args)+2, len(key), key))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n$%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	err := client.sendBytes(finalBytes)
	if err != nil {
		return err
	}
	_, err = client.p.Read()
	if err != nil {
		return err
	}
	_, err = client.p.ParseNull()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return err
	}
	return err
}

func (client *Client) RPop(key string) (string, error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nRPOP\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return "", err
	}
	_, err = client.p.Read()
	if err != nil {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if bytesDiffer(err) {
		if isRESPNull(err) {
			_, err := client.p.ParseNull()
			return "", err
		} else if isRESPError(err) {
			_, err = client.p.ParseError()
			return "", err
		}
	}
	return result, err
}

func (client *Client) LLen(key string) (int, error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLLEN\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return 0, err
	}
	_, err = client.p.Read()
	if err != nil {
		return 0, err
	}
	result, _, err := client.p.ParseUInt()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return 0, err
	}
	return result, err
}

func (client *Client) LPop(key string) (string, error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLPOP\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return "", err
	}
	_, err = client.p.Read()
	if err != nil {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if bytesDiffer(err) {
		if isRESPNull(err) {
			_, err := client.p.ParseNull()
			return "", err
		} else if isRESPError(err) {
			_, err = client.p.ParseError()
			return "", err
		}
	}
	return result, err
}

func (client *Client) LPush(key string, args ...string) error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nLPUSH\r\n$%d\r\n%v", len(args)+2, len(key), key))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n$%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	err := client.sendBytes(finalBytes)
	if err != nil {
		return err
	}
	_, err = client.p.Read()
	if err != nil {
		return err
	}
	_, err = client.p.ParseNull()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return err
	}
	return err
}

func (client *Client) LIndex(key string, index int) (string, error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$6\r\nLINDEX\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(fmt.Sprintf("%v", index)), index))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return "", err
	}
	_, err = client.p.Read()
	if err != nil {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if bytesDiffer(err) {
		if isRESPNull(err) {
			_, err = client.p.ParseNull()
			return "", err
		} else if isRESPError(err) {
			_, err = client.p.ParseError()
			return "", err
		}
	}
	return result, err
}

func (client *Client) Del(key string) error {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nDEL\r\n$%d\r\n%v\r\n", len(key), key))
	err := client.sendBytes(finalBytes)
	if err != nil {
		return err
	}
	_, err = client.p.Read()
	if err != nil {
		return err
	}
	_, err = client.p.ParseNull()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return err
	}
	return err
}

func (client *Client) Ping() (string, error) {
	finalBytes := fmt.Appendf([]byte{}, "*1\r\n$4\r\nPING\r\n")
	err := client.sendBytes(finalBytes)
	if err != nil {
		return "", err
	}
	_, err = client.p.Read()
	if err != nil {
		return "", err
	}
	result, _, err := client.p.ParseBlobString()
	if bytesDiffer(err) && isRESPError(err) {
		_, err = client.p.ParseError()
		return "", err
	}
	return result, nil
}

func (client *Client) sendBytes(b []byte) error {
	_, err := (*client.conn).Write(b)
	if err != nil {
		redigoError := redigoerr.UnableToSendRequestToServer
		redigoError.From = err
		return redigoError
	}
	return nil
}

func bytesDiffer(err error) bool {
	nerr, ok := err.(redigoerr.Error)
	return nerr.Code == 5 && ok
}

func isRESPNull(err error) bool {
	nerr, ok := err.(redigoerr.Error)
	return nerr.ExtraContext["received"] == "_" && ok
}

func isRESPError(err error) bool {
	nerr, ok := err.(redigoerr.Error)
	return nerr.ExtraContext["received"] == "-" && ok
}
