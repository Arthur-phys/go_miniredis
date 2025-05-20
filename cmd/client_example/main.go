// An example of creating & using a REDIGO client
package main

import (
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/client"
)

func main() {
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
	err = c.LPush("Gatos", "Niji", "Anubis", "Pingüica", "Don Bigos")
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
