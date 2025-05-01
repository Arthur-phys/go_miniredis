package server

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"miniredis/core/caches"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func TestWorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_A_Single_One(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator()
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 10240, 1)

	var genericConn net.Conn
	newConnection := newMockConnection()
	genericConn = &newConnection
	go func() {
		newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n"))
		time.Sleep(100 * time.Millisecond)
		newConnection.Close()
	}()
	newWorker.handleConnection(&genericConn)
	if string(newConnection.readAsClient()) != "_\r\n" {
		t.Errorf("Unexpected message received! %v - %v", string(newConnection.responseArray), newConnection.responseArray)
	}
}

func TestWorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_In_A_Single_Package(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator()
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 10240, 1)

	var genericConn net.Conn
	newConnection := newMockConnection()
	genericConn = &newConnection
	go func() {
		newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
		time.Sleep(100 * time.Millisecond)
		newConnection.Close()
	}()

	newWorker.handleConnection(&genericConn)
	responseArray := newConnection.readAsClient()
	if string(responseArray) != "_\r\n$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(responseArray))
	}
}

func TestWorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_In_A_Multiple_Packages(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator()
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 10240, 1)

	var genericConn net.Conn
	newConnection := newMockConnection()
	genericConn = &newConnection
	go func() {
		newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r"))
		fmt.Println("Finished writing first time")
	}()
	newWorker.handleConnection(&genericConn)
	responseArray := newConnection.readAsClient()
	fmt.Println("Finished reading first time")
	if string(responseArray) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(responseArray))
	}
	go func() {
		newConnection.writeAsClient(fmt.Appendf([]byte{}, "\n$3\r\nGET\r\n$1\r\nB\r\n"))
		time.Sleep(100 * time.Millisecond)
		newConnection.Close()
	}()
	responseArray = newConnection.readAsClient()
	if string(responseArray) != "$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(responseArray))
	}
}

func TestWorkerhandleConnection_Should_Return_Error_To_Client_When_Sent_Multiple_Commands_With_One_Wrong(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	workerInstantiator := NewWorkerInstantiator()
	cacheStore := caches.NewSimpleCacheStore()
	channel := make(chan net.Conn)
	newWorker := workerInstantiator(cacheStore, channel, 10240, 1)

	var genericConn net.Conn
	newConnection := newMockConnection()
	genericConn = &newConnection
	go func() {
		newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*1\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
		time.Sleep(100 * time.Millisecond)
		newConnection.Close()
	}()
	newWorker.handleConnection(&genericConn)
	responseArray := newConnection.readAsClient()
	if responseArray[0] != '-' {
		t.Errorf("Unexpected response! %v", string(responseArray))
	}
}

type mockConnection struct {
	requestArray     []byte
	responseArray    []byte
	currentBytesRead int
	closed           bool
	newData          bool
	mutex            *sync.Mutex
	conditional      *sync.Cond
}

func newMockConnection() mockConnection {
	mc := mockConnection{
		requestArray:     []byte{},
		responseArray:    []byte{},
		currentBytesRead: 0,
		closed:           false,
		newData:          true,
		mutex:            &sync.Mutex{},
	}
	mc.conditional = sync.NewCond(mc.mutex)
	return mc
}

func (mc *mockConnection) Read(b []byte) (int, error) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	for !mc.closed && len(mc.requestArray) == 0 {
		mc.conditional.Wait()
	}
	if mc.closed {
		return 0, io.EOF
	}
	n := int(math.Min(float64(len(b)), float64(len(mc.requestArray))))
	copy(b, mc.requestArray[:n])
	mc.requestArray = mc.requestArray[n:]
	return n, nil
}

func (mc *mockConnection) writeAsClient(b []byte) (int, error) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	defer mc.conditional.Broadcast()
	if mc.closed {
		return 0, io.EOF
	}
	mc.requestArray = append(mc.requestArray, b...)
	return len(b), nil
}

func (mc *mockConnection) Write(b []byte) (int, error) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	defer mc.conditional.Broadcast()
	if mc.closed {
		return 0, io.EOF
	}
	mc.responseArray = append(mc.responseArray, b...)
	return len(b), nil
}

func (mc *mockConnection) readAsClient() []byte {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.responseArray
}

func (mc *mockConnection) Close() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.closed = true
	mc.conditional.Broadcast()
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
