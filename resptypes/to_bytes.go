package resptypes

import (
	"fmt"
	e "miniredis/error"
)

func BlobStringToBytes(s string) []byte {
	return fmt.Appendf([]byte{'$'}, fmt.Sprintf("%d\r\n%v\r\n", len(s), s))
}

func IntToBytes(i int) []byte {
	return fmt.Appendf([]byte{':'}, fmt.Sprintf("%v\r\n", i))
}

func NullToBytes() []byte {
	return []byte{'_', '\r', '\n'}
}

func ErrToBytes(err e.Error) []byte {
	return fmt.Appendf([]byte{'-'}, fmt.Sprintf("%v\r\n", err.ClientContext))
}
