package parser

import (
	"bufio"
	"io"
	e "miniredis/error"
)

// Simple wrapper for a bufio.Reader
type Stream struct {
	byteReader *bufio.Reader
}

func (s *Stream) ReadUntilFound(delim byte) ([]byte, error) {
	return s.byteReader.ReadBytes(delim)
}

func (s *Stream) ReadUntilSliceFound(delim []byte) (bytes []byte, err error) {
	if len(delim) == 0 {
		return []byte{}, e.Error{} // Change
	}
	bytes, err = s.byteReader.ReadBytes(delim[0])
	if err != nil {
		return
	}
	bytes = bytes[:len(bytes)-1]
	for i := 1; i < len(delim); i++ {
		newByte, err := s.TakeOne()
		if err != nil {
			return bytes, err
		}
		if newByte != delim[i] {
			return bytes, e.Error{}
		}
	}
	return
}

func (s *Stream) ReadNBytes(n int) ([]byte, int, error) {
	readBytes := make([]byte, n)
	m, err := io.ReadFull(s.byteReader, readBytes)
	return readBytes, m, err
}

func (s *Stream) TakeOne() (byte, error) {
	return s.byteReader.ReadByte()
}

func (s *Stream) Skip(n int) (int, error) {
	return s.byteReader.Discard(n)
}
