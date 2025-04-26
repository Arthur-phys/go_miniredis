package parser

import (
	"bufio"
	"bytes"
	"fmt"
	e "miniredis/error"
	"testing"
)

func Test_ParseCommand_Should_Not_Return_Err_When_Passed_Valid_Command_As_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
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
	commands, err = parser.ParseCommand()
	if err.Code != 3 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 1 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
}

func Test_ParseArray_Should_Return_Array_When_Passed_Valid_Bytes(t *testing.T) {
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	parser := RESPParser{}
	parser.rawBuffer = incomingBytes
	parser.buffer = bufio.NewReader(bytes.NewReader(incomingBytes))
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (string, int, e.Error) {
		return r.BlobStringFromBytes()
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
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.UIntFromBytes()
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
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.UIntFromBytes()
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
	arr, _, err := ParseArray(&parser, func(r *RESPParser) (int, int, e.Error) {
		return r.UIntFromBytes()
	})
	if err.Code == 0 {
		t.Errorf("Error did not happen! %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("Unexpected len for array! %d", len(arr))
	}
}
