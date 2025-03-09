package parser

import (
	"bufio"
	"io"
	"math"
	"testing"
)

func TestReadNBytes_Should_Read_400_Bytes_When_Asked(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 1000)
	stream := Stream{bufio.NewReader(&mockReader)}
	readBytes, m, err := stream.ReadNBytes(400)
	if err != nil {
		t.Errorf("Something happened! %e", err)
	}
	if m != 400 {
		t.Errorf("Unable to read 400 bytes into buffer! %v", len(readBytes))
	}
}

func TestReadNBytes_Should_Return_Error_When_EOF(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 3)
	stream := Stream{bufio.NewReader(&mockReader)}
	readBytes, m, err := stream.ReadNBytes(400)
	if m != 3 {
		t.Errorf("Unable to fill required bytes! %v", readBytes)
	}
	if err == nil || err != io.ErrUnexpectedEOF {
		t.Errorf("Unable to create an error when exceeding capacity!: %e", err)
	}
}

func TestReadUntilFound_Should_Find_Instance(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 8)
	stream := Stream{bufio.NewReader(&mockReader)}
	readBytes, err := stream.ReadUntilFound('l')
	if err != nil {
		t.Errorf("Unable to find 'l' inside stream!: %e", err)
	}
	if len(readBytes) != 3 {
		t.Errorf("readBytes len is not 3! %v", len(readBytes))
	}
}

func TestReadUntilFound_Should_Read_Whole_Array_When_Not_Found(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 4)
	stream := Stream{bufio.NewReader(&mockReader)}
	bytesRead, err := stream.ReadUntilFound('x')
	if err != io.EOF {
		t.Errorf("Unexpected error occurred! %e", err)
	}
	if len(bytesRead) != 4 {
		t.Errorf("Unexpected len for bytes read is not 4! %v", len(bytesRead))
	}
}

func TestSkip_Should_Skip_N_Bytes_When_Asked(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 80)
	stream := Stream{bufio.NewReader(&mockReader)}
	n, err := stream.Skip(10)
	if err != nil {
		t.Errorf("Unexpected error occurred! %e", err)
	}
	if n != 10 {
		t.Errorf("Unexpected byte skip!: %d", n)
	}
}

func TestSkip_Should_Return_Error_When_N_Bigger_Than_Reader_Size(t *testing.T) {
	mockReader := NewMockReader([]byte{'h', 'o', 'l', 'a'}, 80)
	stream := Stream{bufio.NewReader(&mockReader)}
	n, err := stream.Skip(81)
	if err == nil {
		t.Errorf("Error did not happen!")
	}
	if n != 80 {
		t.Errorf("Unexpected byte skip!: %d", n)
	}
}

type MockReader struct {
	bytesArr         []byte
	currentBytesRead int
	limitBytes       int
	read             func(b []byte, mc *MockReader) (int, error)
}

func NewMockReader(bytes []byte, limitBytes int) MockReader {
	read := func(b []byte, mc *MockReader) (int, error) {
		n, i := int(math.Min(float64(len(b)), float64(len(mc.bytesArr)))), 0
		for i < n {
			b[i] = mc.bytesArr[i]
			i++
			mc.currentBytesRead++
			if mc.currentBytesRead >= mc.limitBytes {
				defer func() { mc.read = func(b []byte, mc *MockReader) (int, error) { return 0, io.EOF } }()
				return i, io.EOF
			}
		}
		return n, nil
	}
	return MockReader{bytes, 0, limitBytes, read}
}

func (mc *MockReader) Read(b []byte) (int, error) {
	return mc.read(b, mc)
}
