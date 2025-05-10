package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Arthur-phys/redigo/pkg/client"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

var ipAddress string
var port uint

func init() {
	flag.StringVar(&ipAddress, "ip", "127.0.0.1", "IP address to connect to.")
	flag.UintVar(&port, "port", 6543, "Server port to connect to.")
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

	conn, connErr := net.Dial("tcp", net.JoinHostPort(ipAddress, fmt.Sprintf("%d", port)))
	if connErr != nil {
		fmt.Printf("Fatal error occurred! %v", connErr)
	}
	c := client.New(&conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, readerr := reader.ReadString('\n')
		if readerr != nil {
			fmt.Printf("Error occurred while reading - %e\n", readerr)
			continue
		}
		commands := strings.Split(line[:len(line)-1], " ")
		var result any
		var err e.Error
		switch commands[0] {
		case "GET":
			if len(commands) != 2 {
				fmt.Printf("Incorrect length for command 'GET' - %d\n", len(commands))
				continue
			}
			result, err = c.Get(commands[1])
		case "SET":
			if len(commands) != 3 {
				fmt.Printf("Incorrect length for command 'SET' - %d\n", len(commands))
				continue
			}
			err = c.Set(commands[1], commands[2])
		case "RPUSH":
			if len(commands) < 3 {
				fmt.Printf("Insufficient length for command 'RPUSH' - %d\n", len(commands))
				continue
			}
			err = c.RPush(commands[1], commands[2:]...)
		case "RPOP":
			if len(commands) != 2 {
				fmt.Printf("Incorrect length for command 'RPOP' - %d\n", len(commands))
				continue
			}
			result, err = c.RPop(commands[1])
		case "LPUSH":
			if len(commands) < 3 {
				fmt.Printf("Insufficient length for command 'LPUSH' - %d\n", len(commands))
				continue
			}
			err = c.RPush(commands[1], commands[2:]...)
		case "LPOP":
			if len(commands) != 2 {
				fmt.Printf("Incorrect length for command 'LPOP' - %d\n", len(commands))
				continue
			}
			result, err = c.LPop(commands[1])
		case "LLEN":
			if len(commands) != 2 {
				fmt.Printf("Incorrect length for command 'LLEN' - %d\n", len(commands))
				continue
			}
			result, err = c.LLen(commands[1])
		case "LINDEX":
			if len(commands) != 3 {
				fmt.Printf("Incorrect length for command 'LLINDEX' - %d\n", len(commands))
				continue
			}
			tmpInt, atoiErr := strconv.Atoi(commands[2])
			if atoiErr != nil {
				fmt.Printf("Could not convert index to integer - %e\n", atoiErr)
			}
			result, err = c.LIndex(commands[1], tmpInt)
		case "DEL":
			if len(commands) != 2 {
				fmt.Printf("Incorrect length for command 'DEL' - %d\n", len(commands))
				continue
			}
			err = c.Del(commands[1])
		default:
			fmt.Println("Command specified not found")
		}

		if err.Code != 0 {
			fmt.Printf("Error occurred while processing command - %v\n", err)
		} else if result != nil {
			fmt.Printf("- %s\n", result)
		} else {
			fmt.Println("- OK")
		}
	}
}
