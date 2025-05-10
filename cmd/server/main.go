package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/Arthur-phys/redigo/pkg/core/caches"
	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	"github.com/Arthur-phys/redigo/pkg/server"
)

var ipAddress string
var port uint
var strCache string
var messageSizeLimit int
var workerSize uint
var keepAlive int64

func init() {
	flag.StringVar(&ipAddress, "ip", "127.0.0.1", "Binding IP address for server.")
	flag.UintVar(&port, "port", 6543, "Binding Port for server.")
	flag.StringVar(&strCache, "cache_type", "simpleCache", "Type of cache to use. Available caches are:\n - 'simpleCache': Has a simple interface and sub-optimal alogrithms. It mostly served as a placeholder for tests in this library. Nonetheless, it can be useful if you plan to use redigo for a small app with not too much traffic\n - 'shardedCache': A cache made to be used by a large application that allows multiple goroutines to access the dictionary at the same time granted the key looked for is in different shards.")
	flag.IntVar(&messageSizeLimit, "message_size", 10240, "Limit in size (bytes) for a single message delivered to the server. This size does not correspond with the size of a given TCP segment, rather it is the sum of the size of the content of one or more segments.\nWhenever you send a command to the server which it is unable to be parsed in a single call to read (the socket), every extra call will be summed to the total size read. If this limit is surpassed, the message will be discarded and an error will be returned.\nThis also applies to multi-command messages in which one or more messages are able to be parsed but another is incomplete and requires another call to read.")
	flag.UintVar(&workerSize, "worker_size", 10, "Number of workers to initialize. Every worker will be able to attend a connection independently until this connection is closed or the keepAlive parameter is exceeded.")
	flag.Int64Var(&keepAlive, "keep_alive", 15, "Time (in seconds) to keep a connection open if no message is received. Connections will be closed if no message is received during this period. When this happens, the worker will close the connection and start serving another one.")
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
	var cache func() interfaces.CacheStore
	strCache = strings.ToLower(strCache)
	if strCache == "simplecache" {
		cache = caches.NewSimpleCache
	} else if strCache == "shardedCache" {
		cache = caches.NewShardedCache
	} else {
		fmt.Printf("Option for cache not recognized - %s\n", strCache)
		return
	}

	s, err := server.New(
		ipAddress,
		uint16(port),
		cache,
		server.NewWorkerInstantiator(),
		messageSizeLimit,
		workerSize,
		keepAlive,
	)
	if err != nil {
		fmt.Printf("Fatal error occurred - %e \n", err)
		return
	}
	s.Run()
}
