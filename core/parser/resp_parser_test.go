package parser

import (
	"fmt"
	"miniredis/resptypes"
	"testing"
)

func Test_ParseCommand_Should_Not_Return_Err_When_Passed_Valid_Command_As_Bytes(t *testing.T) {
	parser := RESPParser{}
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	stream := resptypes.NewStream(incomingBytes)
	_, err := parser.ParseCommand(&stream)
	if err.Code != 0 {
		t.Errorf("An unexpected error happened! %v", err)
	}
}

func Test_ParseCommand_Should_Return_Err_When_Passed_Smaller_Command_As_Bytes(t *testing.T) {
	parser := RESPParser{}
	incomingBytes := fmt.Appendf([]byte{}, "*3\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	stream := resptypes.NewStream(incomingBytes)
	_, err := parser.ParseCommand(&stream)
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
}

func Test_ParseCommand_Should_Return_Err_When_Passed_Bigger_Array_Command_As_Bytes(t *testing.T) {
	parser := RESPParser{}
	incomingBytes := fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n$1\r\nB\r\n")
	stream := resptypes.NewStream(incomingBytes)
	_, err := parser.ParseCommand(&stream)
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
}

func Test_ParseCommand_Should_Return_Multiple_Functions_When_Passed_Multiple_Commands(t *testing.T) {
	parser := RESPParser{}
	incomingBytes := fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")
	stream := resptypes.NewStream(incomingBytes)
	commands, err := parser.ParseCommand(&stream)
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 2 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
}
