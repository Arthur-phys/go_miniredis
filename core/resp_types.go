package core

import "fmt"

func BlobStringToRESP(s string) []byte {
	return fmt.Appendf([]byte{'$'}, fmt.Sprintf("%d\r\n%v\r\n", len(s), s))
}

func NullToRESP() []byte {
	return []byte{'_', '\r', '\n'}
}
