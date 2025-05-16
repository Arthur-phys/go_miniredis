//go:build !integration && !e2e
// +build !integration,!e2e

package tobytes

import (
	"fmt"
	"testing"

	e "github.com/Arthur-phys/redigo/pkg/error"
)

func TestBlobString_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleStr := "a samplë"
	byteString := BlobString(sampleStr)
	arr := []byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestInt_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleInt := -43
	byteString := Int(sampleInt)
	arr := []byte{':', '-', '4', '3', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestNull_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	byteString := Null()
	arr := []byte{'_', '\r', '\n'}
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}

func TestErrToBytes_Should_Return_Expected_Formatted_Bytes(t *testing.T) {
	sampleErr := e.Error{Content: "HI", Code: 22, From: nil}
	byteString := Err(sampleErr)
	arr := []byte{'-'}
	arr = fmt.Appendf(arr, "\r\n")
	for i := range byteString {
		if byteString[i] != arr[i] {
			t.Errorf("Bytes did not match! %v != %v", byteString[i], arr[i])
		}
	}
}
