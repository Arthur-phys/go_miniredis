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
	e := c.Set("Arturo", "26")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	res, e := c.Get("Arturo")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	fmt.Printf("I got this! %v\n", res)
	e = c.Set("Gene", "Le gustan los gatos")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	res, e = c.Get("Gene")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	fmt.Printf("I got this! %v\n", res)
	e = c.LPush("Gatos", "Niji", "Anubis", "Pingüica", "Don Bigos")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	res, e = c.LPop("Gatos")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	fmt.Printf("I got this! %v\n", res)
	conn.Close()
}
