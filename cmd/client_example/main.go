package main

import (
	"fmt"
	"miniredis/client"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("Fatal error occurred! %v", err)
	}
	c := client.NewClient(&conn, &client.SimpleSender{})
	e := c.Set("Arturo", "26")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	res, e := c.Get("Arturo")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	fmt.Printf("I got this! %v\n", res)
	e = c.Set("Arturo", "27")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	res, e = c.Get("Arturo")
	if e.Code != 0 {
		fmt.Printf("Fatal error occurred! %v - %e\n", e, e.From)
	}
	fmt.Printf("I got this! %v\n", res)
	conn.Close()
}
