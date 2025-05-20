// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
//
// tobytes transforms a given go type into a byte array following RESP protocol
package tobytes

import (
	"fmt"

	e "github.com/Arthur-phys/redigo/pkg/error"
)

func BlobString(s string) []byte {
	return fmt.Appendf([]byte{'$'}, fmt.Sprintf("%d\r\n%v\r\n", len(s), s))
}

func Int(i int) []byte {
	return fmt.Appendf([]byte{':'}, fmt.Sprintf("%v\r\n", i))
}

func Null() []byte {
	return []byte{'_', '\r', '\n'}
}

func Err(err error) []byte {
	redigoError, ok := err.(e.Error)
	if !ok {
		return fmt.Appendf([]byte{'-'}, "Internal Server Error\r\n")
	}
	return fmt.Appendf([]byte{'-'}, fmt.Sprintf("%v\r\n", redigoError.ClientContext))
}

func Pong() []byte {
	return fmt.Appendf([]byte{'$'}, "4\r\nPONG\r\n")
}
