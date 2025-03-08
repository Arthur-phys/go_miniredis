package caches

import "testing"

func TestSet_Should_Return_NIL_When_Set_Is_Done(t *testing.T) {
	cs := NewSimpleCacheStore()
	err := cs.Set("key", "value")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
}

func TestGet_Should_Return_False_When_Key_Not_Present(t *testing.T) {
	cs := NewSimpleCacheStore()
	if _, ok := cs.Get("key"); ok {
		t.Errorf("A key is present when it should not be!")
	}
}

func TestGet_Should_Return_Value_And_True_When_Key_Is_Present(t *testing.T) {
	cs := NewSimpleCacheStore()
	err := cs.Set("key", "value")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if _, ok := cs.Get("key"); !ok {
		t.Errorf("A key is not present when it should be!")
	}
}

func TestRPush_Should_Create_New_List_In_Cache_When_Not_Present(t *testing.T) {
	cs := NewSimpleCacheStore()
	err := cs.RPush("key", "value")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err != nil || val != "value" {
		t.Errorf("An error occured! %e - %v", err, val)
	}
}

func TestRPush_Should_Add_To_List_In_Cache_When_Present(t *testing.T) {
	cs := NewSimpleCacheStore()
	err := cs.RPush("key", "value")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("key", "value_2")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("key"); err != nil || val != "value_2" {
		t.Errorf("An error occured! %e - %v", err, val)
	}
}

func TestRPop_Should_Return_Error_When_List_Is_Not_Present(t *testing.T) {
	cs := NewSimpleCacheStore()
	if _, err := cs.RPop("key"); err == nil {
		t.Errorf("Expected error but obtained nil! %e", err)
	}
}
