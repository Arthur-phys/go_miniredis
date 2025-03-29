package parser

import (
	"fmt"
	e "miniredis/error"
	"testing"
)

func TestBlobStringToRESP_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleStr := "a samplë"
	byteString := BlobStringToRESP(sampleStr)
	arr := []byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestIntToRESP_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleInt := -43
	byteString := IntToRESP(sampleInt)
	arr := []byte{':', '-', '4', '3', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestNullToRESP_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	byteString := NullToRESP()
	arr := []byte{'_', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestErrToRESP_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleErr := e.Error{Content: "HI", Code: 22, From: nil}
	byteString := ErrToRESP(sampleErr)
	arr := []byte{'-'}
	arr = fmt.Appendf(arr, "\r\n")
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}
