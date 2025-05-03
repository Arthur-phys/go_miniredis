//go:build !integration
// +build !integration

package tobytes

import (
	"fmt"
	"testing"

	e "github.com/Arthur-phys/redigo/pkg/error"
)

func TestBlobStringToBytes_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleStr := "a samplë"
	byteString := BlobStringToBytes(sampleStr)
	arr := []byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestIntToBytes_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleInt := -43
	byteString := IntToBytes(sampleInt)
	arr := []byte{':', '-', '4', '3', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestNullToBytes_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	byteString := NullToBytes()
	arr := []byte{'_', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestErrToBytes_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleErr := e.Error{Content: "HI", Code: 22, From: nil}
	byteString := ErrToBytes(sampleErr)
	arr := []byte{'-'}
	arr = fmt.Appendf(arr, "\r\n")
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}
