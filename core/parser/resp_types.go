package parser

import "fmt"

func BlobStringToRESP(s string) []byte {
	return fmt.Appendf([]byte{'$'}, fmt.Sprintf("%d\r\n%v\r\n", len(s), s))
}

func IntToRESP(i int) []byte {
	return fmt.Appendf([]byte{':'}, fmt.Sprintf("%v\r\n", i))
}

func NullToRESP() []byte {
	return []byte{'_', '\r', '\n'}
}

func ErrToRESP(err error) []byte {
	return fmt.Appendf([]byte{'-'}, fmt.Sprintf("%v\r\n", err))
}
