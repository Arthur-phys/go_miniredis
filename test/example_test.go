//go:build e2e
// +build e2e

// An example of creating & using a REDIGO client
package e2e

import (
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/client"
	"github.com/Arthur-phys/redigo/pkg/server"
)

func ExampleClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("Fatal error occurred! %v", err)
	}
	c := client.New(&conn)
	err = c.Set("Arturo", "26")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	res, err := c.Get("Arturo")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	fmt.Printf("I got this! %v\n", res)
	err = c.Set("Gene", "Le gustan los gatos")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	res, err = c.Get("Gene")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	fmt.Printf("I got this! %v\n", res)
	err = c.LPush("Gatos", "Niji", "Anubis", "Don Bigos", "Pingüica")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	res, err = c.LPop("Gatos")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	fmt.Printf("I got this! %v\n", res)
	conn.Close()
}

func ExampleServer() {
	serverConfig := server.Configuration{
		IpAddress:         "127.0.0.1",
		Port:              8000,
		WorkerAmount:      1,
		KeepAlive:         15,
		MessageSizeLimit:  10240,
		ShutdownTolerance: 5,
	}

	s, err := server.New(&serverConfig)
	if err != nil {
		fmt.Printf("Fatal error occurred - %v\n", err)
		return
	}
	s.Run()
}
