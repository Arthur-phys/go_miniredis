package client

import (
	"fmt"
	e "miniredis/error"
	rt "miniredis/resptypes"
)

type SimpleSender struct{}

func (ss *SimpleSender) get(key string) ([]byte, func(s *rt.Stream) (string, e.Error)) {
	return fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%v\r\n", len(key), key)),
		rt.BlobStringFromBytes

}

func (ss *SimpleSender) set(key string, value string) ([]byte, func(s *rt.Stream) e.Error) {
	return fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(value), value)), rt.ErrorFromBytes
}

func (ss *SimpleSender) rPush(key string, args ...string) ([]byte, func(s *rt.Stream) e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nRPUSH", len(args)+1))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	return finalBytes, rt.ErrorFromBytes
}

func (ss *SimpleSender) rPop(key string) ([]byte, func(s *rt.Stream) (string, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nRPOP\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.BlobStringFromBytes
}

func (ss *SimpleSender) lLen(key string) ([]byte, func(s *rt.Stream) (int, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLLEN\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.UIntFromBytes

}

func (ss *SimpleSender) lPop(key string) ([]byte, func(s *rt.Stream) (string, e.Error)) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*2\r\n$4\r\nLPOP\r\n$%d\r\n%v\r\n", len(key), key))
	return finalBytes, rt.BlobStringFromBytes
}

func (ss *SimpleSender) lPush(key string, args ...string) ([]byte, func(s *rt.Stream) e.Error) {
	finalBytes := fmt.Appendf([]byte{}, fmt.Sprintf("*%d\r\n$5\r\nLPUSH", len(args)+1))
	for i := range args {
		finalBytes = fmt.Appendf(finalBytes, fmt.Sprintf("\r\n%d\r\n%v", len(args[i]), args[i]))
	}
	finalBytes = fmt.Appendf(finalBytes, "\r\n")
	return finalBytes, rt.ErrorFromBytes
}

func (ss *SimpleSender) lIndex(key string, index int) ([]byte, func(s *rt.Stream) (string, e.Error)) {
	return fmt.Appendf([]byte{}, fmt.Sprintf("*3\r\n$6\r\nLINDEX\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n", len(key), key, len(fmt.Sprintf("%v", index)), index)), rt.BlobStringFromBytes
}
