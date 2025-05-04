//go:build !integration && !e2e
// +build !integration,!e2e

package caches

import (
	"testing"
)

func TestSet_Should_Return_NIL_When_Set_Is_Done(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.Set("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
}

func TestGet_Should_Return_Error_When_Key_Not_Present(t *testing.T) {
	cs := NewSimpleCache()
	if _, err := cs.Get("key"); err.Code == 0 {
		t.Errorf("A key is present when it should not be!")
	}
}

func TestGet_Should_Return_Value_And_Nil_When_Key_Is_Present(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.Set("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if _, err := cs.Get("key"); err.Code != 0 {
		t.Errorf("A key is not present when it should be!")
	}
}

func TestRPush_Should_Create_NewSimpleCache_List_In_Cache_When_Not_Present(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.RPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err.Code != 0 || val != "value" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestRPush_Should_Add_To_List_In_Cache_When_Present(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.RPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err.Code != 0 || val != "value_2" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestRPop_Should_Return_Error_When_List_Is_Not_Present(t *testing.T) {
	cs := NewSimpleCache()
	if _, err := cs.RPop("key"); err.Code == 0 {
		t.Errorf("Expected error but obtained nil! %v", err)
	}
}

func TestRPop_Should_Delete_List_When_Empty(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.RPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err.Code != 0 || val != "value" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.RPop("key"); err.Code != 1 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestRPop_Should_Remove_Elements_In_Succession(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.RPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err.Code != 0 || val != "value_2" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if val, err := cs.RPop("key"); err.Code != 0 || val != "value" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.RPop("key"); err.Code != 1 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPush_Should_Add_To_List_In_Cache_When_Present(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.LPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("key"); err.Code != 0 || val != "value_2" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestLPop_Should_Delete_List_When_Empty(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.LPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("key"); err.Code != 0 || val != "value" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.LPop("key"); err.Code != 1 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPop_Should_Remove_Elements_In_Succession(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.LPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("key"); err.Code != 0 || val != "value_2" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if val, err := cs.LPop("key"); err.Code != 0 || val != "value" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.LPop("key"); err.Code != 1 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPop_Should_Return_Error_When_List_Is_Not_Present(t *testing.T) {
	cs := NewSimpleCache()
	if _, err := cs.LPop("key"); err.Code == 0 {
		t.Errorf("Expected error but obtained nil! %v", err)
	}
}

func TestLIndex_Should_Return_Element_When_Present_In_List(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.LPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("key", "value_3")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if str, err := cs.LIndex("key", 0); err.Code != 0 || str != "value_3" {
		t.Errorf("Unable to retrieve first value (0) from 'key' list! %v", err)
	}
}

func TestLIndex_Should_Return_Error_When_Index_Not_Present_In_List(t *testing.T) {
	cs := NewSimpleCache()
	err := cs.LPush("key", "value")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("key", "value_2")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("key", "value_3")
	if err.Code != 0 {
		t.Errorf("An error occurred! %v", err)
	}
	if s, err := cs.LIndex("key", 5); err.Code != 1 && err.Code != 2 {
		t.Errorf("Was able to retrieve unexistant value! %v - %s", err, s)
	}
}
