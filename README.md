# Redigo

![Go](https://img.shields.io/badge/Go-1.21+-blue)
![Tests](https://github.com/arthur-phys/go_miniredis/actions/workflows/test.yml/badge.svg)
![License](https://img.shields.io/badge/License-MIT-blue)
![gosec](https://img.shields.io/badge/gosec-Audited-brightgreen)

---

## 🚀 Overview

_Redigo is a REDIS replica server (and client) created with the sole purpose of learning go. As such, it implements the RESP protocol, but only for a small subset of operations. This was a love project made in my free time to truly understand what to do and not do when writing in go._

## ✨ Features

- 📝 Compatible with commands GET, SET, DEL, LPUSH, LPOP, RPUSH, RPOP, LINDEX, LLEN and PING!
- ⚙️⏲️🛑📏 Has a fully realized server which can control the **number of goroutines spawned**, timeout for sessions, **graceful shutdown** and **maximum size for a message**!
- 📬🧩 Supports **multiple messages sent on a single request**. It even holds a buffer in case you delvier partial messages (so that you can finish sending it in the same connection at a later point)!
- 🔗🧰 Has a client derived from server-created structures and functions that can be used in any project!
- 💻🗣️ Has a REPL program built on top of the client, much like REDIS has one!

## 📦 Instalation

```sh
go install github.com/Arthur-phys/redigo/cmd/redigo_server@latest # For the server
go install github.com/Arthur-phys/redigo/cmd/redigo_cli@latest # For the CLI

```
_or clone and build:_
```sh
git clone https://github.com/Arthur-phys/redigo.git
cd redigo
go build ./cmd/redigo_server
go build ./cmd/redigo_cli
```

## 🛠️ Usage

### 🖥️ For the redigo_server

```sh
redigo_server --help
```
_Example:_
```sh
redigo_server --ip=127.0.0.1 --port=6379 --worker_amount=30 --message_size=10240 --keep_alive=3600 --shutdown=30
```

### 🗣️ For the redigo_cli

```sh
redigo_cli
```
_Use `EXIT` to exit the REPL._

## 📚 Examples

You can also use the client in your code!

```go
// Example Go usage if your project is a library
package e2e

import (
	"fmt"
	"net"

	"github.com/Arthur-phys/redigo/pkg/client"
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
	err = c.Set("Gene", "Anubis")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	res, err = c.Get("Gene")
	if err != nil {
		fmt.Printf("Fatal error occurred! %e\n", err)
	}
	fmt.Printf("I got this! %v\n", res)
	err = c.LPush("Gatos", "Niji", "Anubis", "Bigotes", "Pingüica")
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
```

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgements

- [Go](https://golang.org/)
- [Redis](https://redis.io/)
- [Layout project in Golang](https://github.com/golang-standards/project-layout/tree/master)
- [Practical Go](https://dave.cheney.net/practical-go/presentations/gophercon-israel.html)
- Everyone else in my life ❤️

---
