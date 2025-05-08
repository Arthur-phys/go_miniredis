//go:build e2e
// +build e2e

package e2e

import (
	"net"
	"testing"

	"github.com/Arthur-phys/redigo/pkg/client"
	"github.com/Arthur-phys/redigo/pkg/core/caches"
	"github.com/Arthur-phys/redigo/pkg/server"
)

func TestE2E_Client_Full(t *testing.T) {

	s, err := server.New(
		"127.0.0.1",
		8001,
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
	t.Run("Command=GET,Response=Null", e2e_Client_That_Sends_A_GET_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=SET,Response=Null", e2e_Client_That_Sends_A_SET_Should_Receive_Null_As_Response)
	t.Run("Command=GET,Response=String", e2e_Client_That_Sends_A_GET_Should_Receive_String_If_Key_Is_Present)
	t.Run("Command=RPOP,Response=Null", e2e_Client_That_Sends_An_RPOP_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=LPOP,Response=Null", e2e_Client_That_Sends_An_LPOP_Should_Receive_Null_If_Key_Is_Not_Present)
	t.Run("Command=RPUSH,Response=Null", e2e_Client_That_Sends_An_RPUSH_Should_Receive_Null)
	t.Run("Command=RPOP,Response=String", e2e_Client_That_Sends_An_RPOP_Should_Receive_String_If_Key_Is_Present)
	t.Run("Command=LPUSH,Response=Null", e2e_Client_That_Sends_An_LPUSH_Should_Receive_Null)
	t.Run("Command=LINDEX,Response=Null", e2e_Client_That_Sends_An_LINDEX_Should_Receive_Null_If_Key_Is_Present_But_Index_Is_Invalid)
	t.Run("Command=LINDEX,Response=String", e2e_Client_That_Sends_An_LINDEX_Should_Receive_String_If_Key_Is_Present_And_Index_Is_Valid)
	t.Run("Command=LLEN,Response=Int", e2e_Client_That_Sends_An_LLEN_Should_Receive_List_Size_If_Key_Is_Present)
	t.Run("Command=LPOP,Response=String", e2e_Client_That_Sends_An_LPOP_Should_Receive_String_If_Key_Is_Present)

}

func e2e_Client_That_Sends_A_GET_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)

	str, newErr := c.Get("R")
	if newErr.Code != 0 || str != "" {
		t.Errorf("Unexpected error occurred! %v, %v - %s", newErr, newErr.ExtraContext, str)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("Unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_A_SET_Should_Receive_Null_As_Response(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)

	newErr := c.Set("R", "REDIGO")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %e", err)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("Unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_A_GET_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	response, newErr := c.Get("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "REDIGO" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_An_RPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	response, newErr := c.RPop("V")

	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_An_LPOP_Should_Receive_Null_If_Key_Is_Not_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	response, newErr := c.LPop("V")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_An_RPUSH_Should_Receive_Null(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	newErr := c.RPush("R", "REDIGO", "NIJI", "BIGOTES", "ANUBIS")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

}

func e2e_Client_That_Sends_An_RPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)
	response, newErr := c.RPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "ANUBIS" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	response, newErr = c.RPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "BIGOTES" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	response, newErr = c.RPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "NIJI" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	response, newErr = c.RPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "REDIGO" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Client_That_Sends_An_LPUSH_Should_Receive_Null(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	newErr := c.LPush("R", "REDIGO", "NIJI", "BIGOTES", "ANUBIS")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Client_That_Sends_An_LINDEX_Should_Receive_Null_If_Key_Is_Present_But_Index_Is_Invalid(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	response, newErr := c.LIndex("R", 5)
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Client_That_Sends_An_LINDEX_Should_Receive_String_If_Key_Is_Present_And_Index_Is_Valid(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}

	c := client.New(&conn)
	response, newErr := c.LIndex("R", 3)
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "REDIGO" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Client_That_Sends_An_LLEN_Should_Receive_List_Size_If_Key_Is_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)
	response, newErr := c.LLen("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != 4 {
		t.Errorf("Unexpected list size retrieved! %d", response)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}

func e2e_Client_That_Sends_An_LPOP_Should_Receive_String_If_Key_Is_Present(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
	c := client.New(&conn)
	response, newErr := c.LPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "ANUBIS" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}
	response, newErr = c.LPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "BIGOTES" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}
	response, newErr = c.LPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "NIJI" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}
	response, newErr = c.LPop("R")
	if newErr.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", newErr)
	}
	if response != "REDIGO" {
		t.Errorf("Unexpected string retrieved! %s", response)
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("An unexpected error occurred! %e", err)
	}
}
