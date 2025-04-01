package client

import (
	"fmt"
	e "miniredis/error"
	rt "miniredis/resptypes"
)

type SimpleSender struct{}

func (ss *SimpleSender) get(key string) ([]byte, func(s *rt.Stream) (string, e.Error)) {
	return fmt.Appendf([]byte{}, fmt.Sprintf(`*2\r\n$3\r\nGET\r\n$%d\r\n%v\r\n`, len(key), key)),
		rt.BlobStringFromBytes

}

func (ss *SimpleSender) set(key string, value string) ([]byte, func(s *rt.Stream) e.Error) {
	return fmt.Appendf([]byte{}, fmt.Sprintf(`*3\r\n$3\r\nSET\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n`, len(key), key, len(value), value)), rt.ErrorFromBytes
}

func (ss *SimpleSender) rPush(key string, args ...string) ([]byte, func([]byte) error) {
	// finalStr := ""
}

func (ss *SimpleSender) rPop(key string) ([]byte, func([]byte) (string, error))              {}
func (ss *SimpleSender) lLen(key string) ([]byte, func([]byte) (uint, error))                {}
func (ss *SimpleSender) lPop(key string) ([]byte, func([]byte) (string, error))              {}
func (ss *SimpleSender) lPush(key string, args ...string) ([]byte, func([]byte) error)       {}
func (ss *SimpleSender) lIndex(key string, index int) ([]byte, func([]byte) (string, error)) {}
