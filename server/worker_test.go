package server

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"miniredis/core/caches"
	"miniredis/core/parser"
	"net"
	"os"
	"testing"
	"time"
)

func TestWorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_A_Single_One(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator(parser.NewRESPParser)
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 1)

	var genericConn net.Conn
	newConnection := newMockConnection(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n"), 33)
	genericConn = &newConnection
	newWorker.handleConnection(&genericConn)
	if string(newConnection.writeArr) != "_\r\n" {
		t.Errorf("Unexpected message received! %v - %v", string(newConnection.writeArr), newConnection.writeArr)
	}
}

func TestWorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_In_A_Single_Package(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator(parser.NewRESPParser)
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 1)

	var genericConn net.Conn
	newConnection := newMockConnection(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"), 53)
	genericConn = &newConnection
	newWorker.handleConnection(&genericConn)
	if string(newConnection.writeArr) != "_\r\n$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(newConnection.writeArr))
	}
}

func TestWorkerhandleConnection_Should_Return_Error_To_Client_When_Sent_Multiple_Commands_With_One_Wrong(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator(parser.NewRESPParser)
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 1)

	var genericConn net.Conn
	newConnection := newMockConnection(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*1\r\n$3\r\nGET\r\n$1\r\nB\r\n"), 53)
	genericConn = &newConnection
	newWorker.handleConnection(&genericConn)
	if newConnection.writeArr[0] != '-' {
		t.Errorf("Unexpected response! %v", newConnection.writeArr)
	}
}

type mockConnection struct {
	bytesArr         []byte
	writeArr         []byte
	currentBytesRead int
	limitBytes       int
	read             func(b []byte, mc *mockConnection) (int, error)
	write            func(b []byte, mc *mockConnection) (int, error)
}

func newMockConnection(bytes []byte, limitBytes int) mockConnection {
	read := func(b []byte, mc *mockConnection) (int, error) {
		n, i := int(math.Min(float64(len(b)), float64(len(mc.bytesArr)))), 0
		for i < n {
			b[i] = mc.bytesArr[i]
			i++
			mc.currentBytesRead++
			if mc.currentBytesRead >= mc.limitBytes {
				return i, io.EOF
			}
		}
		return n, nil
	}
	write := func(b []byte, mc *mockConnection) (int, error) {
		mc.writeArr = append(mc.writeArr, b...)
		return len(mc.writeArr), nil
	}
	return mockConnection{bytes, []byte{}, 0, limitBytes, read, write}
}

func (mc *mockConnection) Read(b []byte) (int, error) {
	return mc.read(b, mc)
}

func (mc *mockConnection) Write(b []byte) (int, error) {
	return mc.write(b, mc)
}

func (mc *mockConnection) Close() error {
	mc.read = func(b []byte, mc *mockConnection) (int, error) { return 0, io.EOF }
	mc.write = func(b []byte, mc *mockConnection) (int, error) { return 0, io.EOF }
	return nil
}

func (mc *mockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
func (mc *mockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) LocalAddr() net.Addr {
	return mockAddr{}
}

func (mc *mockConnection) RemoteAddr() net.Addr {
	return mockAddr{}
}

type mockAddr struct{}

func (ma mockAddr) Network() string {
	return "no network"
}

func (ma mockAddr) String() string {
	return "test"
}
