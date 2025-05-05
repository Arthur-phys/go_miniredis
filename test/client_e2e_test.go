package e2e

import (
	"net"
	"testing"

	"github.com/Arthur-phys/redigo/pkg/client"
	"github.com/Arthur-phys/redigo/pkg/core/caches"
	"github.com/Arthur-phys/redigo/pkg/server"
)

func TestE2E_Client_That_Sends_A_SET_Should_Receive_Null_As_Response(t *testing.T) {

	s, err := server.New(
		"127.0.0.1",
		8000,
		caches.NewSimpleCache,
		server.NewWorkerInstantiator(),
		10240,
		1,
		15,
	)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	go func() {
		s.Run()
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)

	newErr := c.Set("R", "REDIGO")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %e", err)
	}

}

func TestE2E_Client_That_Sends_A_GET_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	s, err := server.New(
		"127.0.0.1",
		8000,
		caches.NewSimpleCache,
		server.NewWorkerInstantiator(),
		10240,
		1,
		15,
	)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	go func() {
		s.Run()
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)

	str, newErr := c.Get("R")
	if newErr.Code != 0 || str != "" {
		t.Errorf("Unexpected error occurred! %v, %v - %s", newErr, newErr.ExtraContext, str)
	}

}

// func TestE2E_Client_That_Sends_A_GET_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nR\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 12 || string(response[:n]) != "$6\r\nREDIGO\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
// 	}
// 	err = conn.Close()
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// }

// func TestE2E_Client_That_Sends_An_RPUSH_Should_Receive_Null(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nRPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// }

// func TestE2E_Client_That_Sends_An_LPUSH_Should_Receive_Null(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// }

// func TestE2E_Client_That_Sends_An_RPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_An_RPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nRPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 12 || string(response[:n]) != "$6\r\nANUBIS\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// }

// func TestE2E_Client_That_Sends_An_LPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_An_LPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 12 || string(response[:n]) != "$6\r\nANUBIS\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
// 	}

// }

// func TestE2E_Client_That_Sends_An_LINDEX_Should_Receive_Null_If_Key_Is_Present_But_Index_Is_Invalid(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$6\r\nLINDEX\r\n$1\r\nR\r\n$1\r\n5\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
// 	}

// }

// func TestE2E_Client_That_Sends_An_LINDEX_Should_Receive_String_If_Key_Is_Present_And_Index_Is_Valid(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$6\r\nLINDEX\r\n$1\r\nR\r\n$1\r\n3\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 12 || string(response[:n]) != "$6\r\nREDIGO\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
// 	}

// }

// func TestE2E_Client_That_Sends_An_LLEN_Should_Receive_List_Size_If_Key_Is_Present(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 || string(response[:n]) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLLEN\r\n$1\r\nR\r\n"))
// 	n, err = conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 4 || string(response[:n]) != ":4\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
// 	}

// }

// func TestE2E_Client_That_Sends_A_Shorter_Message_Than_It_Should_Would_Receive_Error(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 20 || string(response[:n]) != "-Command malformed\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_A_Larger_Message_Than_It_Should_Would_Receive_Error(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*4\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$1\r\nB\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 20 || string(response[:n]) != "-Command malformed\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_A_Partial_Message_Will_Receive_Response_Until_Message_Is_Complete(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r"))
// 	conn.Write(fmt.Appendf([]byte{}, "\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 3 && string(response) != "_\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n*2\r\n$3\r\nGET\r\n$1\r\nR\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 15 && string(response) != "_\r\n$6\r\nREDIGO\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }

// func TestE2E_Client_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses_Different_Commands(t *testing.T) {

// 	s, err := server.New(
// 		"127.0.0.1",
// 		8000,
// 		caches.NewSimpleCache,
// 		server.NewWorkerInstantiator(),
// 		10240,
// 		1,
// 		15,
// 	)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	response := make([]byte, 50)
// 	go func() {
// 		s.Run()
// 	}()
// 	conn, err := net.Dial("tcp", "127.0.0.1:8000")
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}

// 	conn.Write(fmt.Appendf([]byte{}, "*5\r\n$5\r\nRPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n*2\r\n$4\r\nLLEN\r\nR\r\n"))
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		t.Errorf("An unexpected error occurred! %e", err)
// 	}
// 	if n != 20 && string(response) != "_\r\n$7\r\nBIGOTES\r\n:2\r\n" {
// 		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
// 	}

// }
