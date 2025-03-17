package parser

import (
	"bufio"
	"io"
	"math"
	"testing"
)

func TestReadUntilFound_Should_Find_Instance(t *testing.T) {
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 8)
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
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 4)
	stream := Stream{bufio.NewReader(&mockReader)}
	bytesRead, err := stream.ReadUntilFound('x')
	if err != io.EOF {
		t.Errorf("Unexpected error occurred! %e", err)
	}
	if len(bytesRead) != 4 {
		t.Errorf("len for bytesRead is not 4! %v", len(bytesRead))
	}
}

func TestReadUntilSliceFound_Should_Find_Whole_Slice_When_Present(t *testing.T) {
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 6)
	stream := Stream{bufio.NewReader(&mockReader)}
	bytesRead, err := stream.ReadUntilSliceFound([]byte{'l', 'a'})
	if err != nil {
		t.Errorf("Unexpected error occurred! %e", err)
	}
	if len(bytesRead) != 2 {
		t.Errorf("len for bytesRead is not 2! %v", len(bytesRead))
	}
	mockBytes := []byte{'h', 'o'}
	for i := range bytesRead {
		if mockBytes[i] != bytesRead[i] {
			t.Errorf("bytesRead is not '[h,o]'! %v", bytesRead)
		}
	}
}

func TestReadUntilSliceFound_Should_Return_Whole_Slice_From_Reader_When_Slice_Looked_For_Is_Not_Present(t *testing.T) {
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 8)
	stream := Stream{bufio.NewReader(&mockReader)}
	bytesRead, err := stream.ReadUntilSliceFound([]byte{'l', 'x'})
	if err != io.EOF {
		t.Errorf("Unexpected error occurred! %e", err)
	}
	if len(bytesRead) != 8 {
		t.Errorf("len for bytesRead is not 8, reader was not exhausted! %v", len(bytesRead))
	}
	mockBytes := []byte{'h', 'o', 'l', 'a', 'h', 'o', 'l', 'a'}
	for i := range bytesRead {
		if mockBytes[i] != bytesRead[i] {
			t.Errorf("bytesRead is not '[h,o,l,a,h,o,l,a]'! %v", bytesRead)
		}
	}
}

func TestReadNBytes_Should_Read_400_Bytes_When_Asked(t *testing.T) {
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 1000)
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
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 3)
	stream := Stream{bufio.NewReader(&mockReader)}
	readBytes, m, err := stream.ReadNBytes(400)
	if m != 3 {
		t.Errorf("Unable to fill required bytes! %v", readBytes)
	}
	if err == nil || err != io.ErrUnexpectedEOF {
		t.Errorf("Unable to create an error when exceeding capacity!: %e", err)
	}
}

func TestSkip_Should_Skip_N_Bytes_When_Asked(t *testing.T) {
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 80)
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
	mockReader := NewmockReader([]byte{'h', 'o', 'l', 'a'}, 80)
	stream := Stream{bufio.NewReader(&mockReader)}
	n, err := stream.Skip(81)
	if err == nil {
		t.Errorf("Error did not happen!")
	}
	if n != 80 {
		t.Errorf("Unexpected byte skip!: %d", n)
	}
}

type mockReader struct {
	bytesArr         []byte
	currentBytesRead int
	limitBytes       int
	read             func(b []byte, mc *mockReader) (int, error)
}

func NewmockReader(bytes []byte, limitBytes int) mockReader {
	read := func(b []byte, mc *mockReader) (int, error) {
		n, i := int(math.Min(float64(len(b)), float64(len(mc.bytesArr)))), 0
		for i < n {
			b[i] = mc.bytesArr[i]
			i++
			mc.currentBytesRead++
			if mc.currentBytesRead >= mc.limitBytes {
				defer func() { mc.read = func(b []byte, mc *mockReader) (int, error) { return 0, io.EOF } }()
				return i, io.EOF
			}
		}
		return n, nil
	}
	return mockReader{bytes, 0, limitBytes, read}
}

func (mc *mockReader) Read(b []byte) (int, error) {
	return mc.read(b, mc)
}
