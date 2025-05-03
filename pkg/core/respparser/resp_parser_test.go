package respparser

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	e "github.com/Arthur-phys/redigo/pkg/error"
)

func Test_ParseCommand_Should_Not_Return_Err_When_Passed_Valid_Command_As_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	parser.maxBytesPerCallAllowed = 10240
	parser.totalBytesRead = 0
	arr, err := parser.ParseCommand()
	if err.Code != 3 {
		t.Errorf("An unexpected error happened! %v", err)
	} else if len(arr) != 1 {
		t.Errorf("Unexpected len for commands! %d", len(arr))
	}
}

func Test_ParseCommand_Should_Return_Err_When_Passed_Smaller_Command_As_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*3\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, err := parser.ParseCommand()
	if err.Code == 0 {
		t.Errorf("Error did not happen! %v", err)
	}
}

func Test_ParseCommand_Should_Return_Err_When_Passed_Bigger_Array_Command_As_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, err := parser.ParseCommand()
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
}

func Test_ParseCommand_Should_Return_Multiple_Functions_When_Passed_Multiple_Commands(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	commands, err := parser.ParseCommand()
	if err.Code != 3 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 2 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
}

func Test_ParseCommand_Should_Return_Multiple_Functions_When_Passed_Multiple_Commands_Sepparatedly(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	commands, err := parser.ParseCommand()
	if err.Code != 3 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 1 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
	parser.rawBufferPosition = 0
	incomingBytes = fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser.buffer.Reset(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	commands, err = parser.ParseCommand()
	if err.Code != 3 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 1 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
}

func Test_ParseBlobString_Should_Return_String_When_Passed_Valid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "$9\r\npingüino\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	str, _, err := parser.ParseBlobString()
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if str != "pingüino" {
		t.Errorf("Unexpected string! %s", str)
	}
}

func Test_ParseBlobString_Should_Return_Error_When_Encountered_Instead_Of_String(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "-Test ERROR!\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, _, err := parser.ParseBlobString()
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
	if err.ExtraContext["text"] != "Test ERROR!" {
		t.Errorf("Unexpected context! %v", err.ExtraContext["text"])
	}
}

func Test_ParseUInt_Should_Return_Int_When_Passed_Valid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, ":2779\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	i, _, err := parser.ParseUInt()
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if i != 2779 {
		t.Errorf("Unexpected int! %d", i)
	}
}

func Test_ParseUInt_Should_Return_Error_When_Encountered_Instead_Of_Int(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "-Test ERROR!\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, _, err := parser.ParseUInt()
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
	if err.ExtraContext["text"] != "Test ERROR!" {
		t.Errorf("Unexpected context! %v", err.ExtraContext["text"])
	}
}

func Test_ParseNull_Should_Return_Nil_When_Passed_Valid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "_\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, err := parser.ParseNull()
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
}

func Test_ParseNull_Should_Return_Error_When_Encountered_Instead_Of_Nil(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "-Test ERROR!\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	_, err := parser.ParseNull()
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
	if err.ExtraContext["text"] != "Test ERROR!" {
		t.Errorf("Unexpected context! %v", err.ExtraContext["text"])
	}
}

func Test_ParseArray_Should_Return_Array_When_Passed_Valid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (string, int, e.Error) {
		return r.ParseBlobString()
	})
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(arr) != 2 {
		t.Errorf("Unexpected len for array! %d", len(arr))
	}
}

func Test_ParseArray_Should_Return_Array_When_Passed_Valid_Bytes_For_UInts(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*5\r\n:22\r\n:555\r\n:127\r\n:488\r\n:999\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.ParseUInt()
	})
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(arr) != 5 {
		t.Errorf("Unexpected len for array! %d", len(arr))
	}
}

func Test_ParseArray_Should_Return_Error_When_Passed_Invalid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*10\r\n:22\r\n:555\r\n:127\r\n:488\r\n:999\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.ParseUInt()
	})
	if err.Code == 0 {
		t.Errorf("Error did not happen! %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("Unexpected len for array! %d", len(arr))
	}
}

func Test_ParseArray_Should_Return_Error_When_Passed_Bytes_With_Multiple_Types(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*4\r\n:22\r\n$3\r\nGET\r\n:488\r\n:999\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	parser.rawBufferEffectiveSize = len(incomingBytes)
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.ParseUInt()
	})
	if err.Code == 0 {
		t.Errorf("Error did not happen! %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("Unexpected len for array! %d", len(arr))
	}
}

func TestReadUntilSliceFound_Should_Find_Whole_Slice_When_Present(t *testing.T) {
	incomingBytes := []byte{'h', 'o', 'l', 'a'}
	stream := RESPParser{}
	stream.rawBuffer = incomingBytes
	stream.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	stream.rawBufferEffectiveSize = len(incomingBytes)
	bytesRead, i, err := stream.readUntilSliceFound([]byte{'l', 'a'})
	if err.Code != 0 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
	if i != 4 {
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
	incomingBytes := []byte{'h', 'o', 'l', 'a'}
	stream := RESPParser{}
	stream.rawBuffer = incomingBytes
	stream.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	stream.rawBufferEffectiveSize = len(incomingBytes)
	bytesRead, i, err := stream.readUntilSliceFound([]byte{'l', 'x'})
	if err.Code != 4 {
		t.Errorf("Unexpected error occurred! %v", err)
	}
	if i != 4 {
		t.Errorf("len for bytesRead is not 4, reader was not exhausted! %v", i)
	}
	mockBytes := []byte{'h', 'o', 'l', 'a'}
	for i := range bytesRead {
		if mockBytes[i] != bytesRead[i] {
			t.Errorf("bytesRead is not '[h,o,l,a]'! %v", bytesRead)
		}
	}
}
