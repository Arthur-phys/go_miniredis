package e2e

import (
	"fmt"
	"net"
	"testing"

	"github.com/Arthur-phys/redigo/pkg/core/caches"
	"github.com/Arthur-phys/redigo/pkg/server"
)

func TestE2E_Server_Full(t *testing.T) {
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
	t.Run("Command=GET,Response=Null", e2e_Connection_That_Sends_A_GET_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=SET,Response=Null", e2e_Connection_That_Sends_A_SET_Should_Receive_Null_As_Response)
	t.Run("Command=GET,Response=String", e2e_Connection_That_Sends_A_GET_Should_Receive_String_If_Key_Is_Present)
	t.Run("Command=RPOP,Response=Null", e2e_Connection_That_Sends_An_RPOP_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=LPOP,Response=Null", e2e_Connection_That_Sends_An_LPOP_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=RPUSH,Response=Null", e2e_Connection_That_Sends_An_RPUSH_Should_Receive_Null)
	t.Run("Command=RPOP,Response=String", e2e_Connection_That_Sends_An_RPOP_Should_Receive_String_If_Key_Is_Present)
	t.Run("Command=LPUSH,Response=Null", e2e_Connection_That_Sends_An_LPUSH_Should_Receive_Null)
	t.Run("Command=LINDEX,Response=Null", e2e_Connection_That_Sends_An_LINDEX_Should_Receive_Null_If_Key_Is_Present_But_Index_Is_Invalid)
	t.Run("Command=LINDEX,Response=String", e2e_Connection_That_Sends_An_LINDEX_Should_Receive_String_If_Key_Is_Present_And_Index_Is_Valid)
	t.Run("Command=LLEN,Response=Int", e2e_Connection_That_Sends_An_LLEN_Should_Receive_List_Size_If_Key_Is_Present)
	t.Run("Command=LPOP,Response=String", e2e_Connection_That_Sends_An_LPOP_Should_Receive_String_If_Key_Is_Present)
	t.Run("Command=Malformed,Response=Error-Short", e2e_Connection_That_Sends_A_Shorter_Message_Than_It_Should_Would_Receive_Error)
	t.Run("Command=Malformed,Response=Error-Large", e2e_Connection_That_Sends_A_Larger_Message_Than_It_Should_Would_Receive_Error)
	t.Run("Command=Incomplete,Response=UntilComplete", e2e_Connection_That_Sends_A_Partial_Message_Will_Receive_Response_Until_Message_Is_Complete)
	t.Run("Command=Multiple,Response=Multiple", e2e_Connection_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses)
	t.Run("Command=Multiple,Response=Multiple_2", e2e_Connection_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses_Different_Commands)
}

func e2e_Connection_That_Sends_A_GET_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_A_SET_Should_Receive_Null_As_Response(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_A_GET_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 12 || string(response[:n]) != "$6\r\nREDIGO\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_RPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_LPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_RPUSH_Should_Receive_Null(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nRPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_RPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 12 || string(response[:n]) != "$6\r\nANUBIS\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 13 || string(response[:n]) != "$7\r\nBIGOTES\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 10 || string(response[:n]) != "$4\r\nNIJI\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nRPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 12 || string(response[:n]) != "$6\r\nREDIGO\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Connection_That_Sends_An_LPUSH_Should_Receive_Null(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*6\r\n$5\r\nLPUSH\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n$6\r\nANUBIS\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_LINDEX_Should_Receive_Null_If_Key_Is_Present_But_Index_Is_Invalid(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$6\r\nLINDEX\r\n$1\r\nR\r\n$1\r\n5\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 || string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_LINDEX_Should_Receive_String_If_Key_Is_Present_And_Index_Is_Valid(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$6\r\nLINDEX\r\n$1\r\nR\r\n$1\r\n2\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 10 || string(response[:n]) != "$4\r\nNIJI\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_LLEN_Should_Receive_List_Size_If_Key_Is_Present(t *testing.T) {
	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLLEN\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 4 || string(response[:n]) != ":4\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_An_LPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 12 || string(response[:n]) != "$6\r\nANUBIS\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 13 || string(response[:n]) != "$7\r\nBIGOTES\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 10 || string(response[:n]) != "$4\r\nNIJI\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$4\r\nLPOP\r\n$1\r\nR\r\n"))
	n, err = conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 12 || string(response[:n]) != "$6\r\nREDIGO\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, string(response))
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_A_Shorter_Message_Than_It_Should_Would_Receive_Error(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*2\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 20 || string(response[:n]) != "-Command malformed\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_A_Larger_Message_Than_It_Should_Would_Receive_Error(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*4\r\n$3\r\nSET\r\n$1\r\nR\r\n$6\r\nREDIGO\r\n$1\r\nB\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 20 || string(response[:n]) != "-Command malformed\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_A_Partial_Message_Will_Receive_Response_Until_Message_Is_Complete(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r"))
	conn.Write(fmt.Appendf([]byte{}, "\n$1\r\nC\r\n$4\r\nCATS\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 3 && string(response) != "_\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses(t *testing.T) {

	response := make([]byte, 50)

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nD\r\n$4\r\nDOGS\r\n*2\r\n$3\r\nGET\r\n$1\r\nD\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 13 && string(response) != "_\r\n$4\r\nDOGS\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Connection_That_Sends_Multiple_Messages_Will_Receive_Multiple_Responses_Different_Commands(t *testing.T) {

	response := make([]byte, 50)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	conn.Write(fmt.Appendf([]byte{}, "*5\r\n$5\r\nRPUSH\r\n$1\r\nA\r\n$7\r\nANIMALS\r\n$4\r\nNIJI\r\n$7\r\nBIGOTES\r\n*2\r\n$4\r\nRPOP\r\n$1\r\nA\r\n*2\r\n$4\r\nLLEN\r\nA\r\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	if n != 20 && string(response) != "_\r\n$7\r\nBIGOTES\r\n:2\r\n" {
		t.Errorf("Unexpected response received! n = %d - response = %v", n, response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}
