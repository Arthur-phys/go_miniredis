package parser

import (
	"bufio"
	"bytes"
	"io"
	e "miniredis/error"
)

// Simple wrapper for a bufio.Reader
type Stream struct {
	byteReader *bufio.Reader
}

func NewStream(b []byte) Stream {
	return Stream{bufio.NewReader(bytes.NewBuffer(b))}
}

func (s *Stream) ReadUntilFound(delim byte) ([]byte, error) {
	return s.byteReader.ReadBytes(delim)
}

func (s *Stream) ReadUntilSliceFound(delim []byte) ([]byte, error) {
	if len(delim) == 0 {
		return []byte{}, e.Error{} // Change
	}
	var sliceFoundRecursive func([]byte, []byte) ([]byte, error)
	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, error) {
		bytes, err := s.byteReader.ReadBytes(delim[0])
		bytesRead = append(bytesRead, bytes...)
		if err != nil {
			return bytesRead, err
		}
		for i := 1; i < len(delim); i++ {
			newByte, err := s.TakeOne()
			if err != nil {
				return bytesRead, err
			}
			bytesRead = append(bytesRead, newByte)
			if newByte != delim[i] {
				return sliceFoundRecursive(delim, bytesRead) // Change
			}
		}
		return bytesRead, nil
	}
	bytes, err := sliceFoundRecursive(delim, []byte{})
	if err == nil {
		bytes = bytes[:len(bytes)-len(delim)]
	}
	return bytes, err
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
