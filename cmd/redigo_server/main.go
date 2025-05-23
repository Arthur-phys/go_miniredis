// A simple CLI tool for the server.
// You can check the options available by using --help
package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/Arthur-phys/redigo/pkg/server"
)

var ipAddress string
var port uint
var messageSizeLimit int
var workerAmount uint64
var keepAlive int64
var shutdownTolerance int64

func init() {
	flag.StringVar(&ipAddress, "ip", "127.0.0.1", "Binding IP address for server.")
	flag.UintVar(&port, "port", 6543, "Binding Port for server.")
	flag.IntVar(&messageSizeLimit, "message_size", 10240, "Limit in size (bytes) for a single message delivered to the server.")
	flag.Uint64Var(&workerAmount, "worker_amount", 10, "Number of workers to initialize.")
	flag.Int64Var(&keepAlive, "keep_alive", 15, "Time (in seconds) to keep a connection open if no message is received.")
	flag.Int64Var(&shutdownTolerance, "shutdown", 15, "Time (in seconds) given to workers when gracefully shutting down the server.")
}

func main() {
	flag.Parse()

	if ok, err := regexp.MatchString(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`, ipAddress); !ok || err != nil {
		fmt.Printf("Invalid IP address - %s\n", ipAddress)
		return
	}
	if port > uint(^uint16(0)) {
		fmt.Printf("Unable to convert given port number (%d) to the corresponding range 0 - 65535\n", port)
		return
	}

	serverConfig := server.Configuration{
		IpAddress:         ipAddress,
		Port:              uint16(port),
		WorkerAmount:      workerAmount,
		KeepAlive:         keepAlive,
		MessageSizeLimit:  messageSizeLimit,
		ShutdownTolerance: shutdownTolerance,
	}

	s, err := server.New(&serverConfig)
	if err != nil {
		fmt.Printf("Fatal error occurred - %v\n", err)
		return
	}
	s.Run()
}
