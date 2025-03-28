package parser

import (
	"fmt"
	"testing"
)

func Test_miniRedisBlobStringFromBytes_Should_Convert_When_Passed_Valid_Input(t *testing.T) {
	parser := RESPParser{}
	stream := NewStream([]byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r', '\n'})

	s, err := parser.miniRedisBlobStringFromBytes(&stream)
	if err.Code != 0 {
		t.Errorf("Unexpected error encountered! %v", err)
	}
	if s != "a samplë" {
		t.Errorf("Unable to obtain string from bytes! %v", s)
	}

}

func Test_miniRedisBlobStringFromBytes_Should_Return_Error_When_Passed_Invalid_Input(t *testing.T) {
	parser := RESPParser{}
	stream := NewStream([]byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r'})

	_, err := parser.miniRedisBlobStringFromBytes(&stream)
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}

}

func Test_ParseCommand_Should_Not_Return_Err_When_Passed_Valid_Command_As_Bytes(t *testing.T) {
	parser := RESPParser{}
	_, err := parser.ParseCommand(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	if err.Code != 0 {
		t.Errorf("An unexpected error happened! %v", err)
	}
}

func Test_ParseCommand_Should_Return_Err_When_Passed_Invalid_Command_As_Bytes(t *testing.T) {
	parser := RESPParser{}
	_, err := parser.ParseCommand(fmt.Appendf([]byte{}, "*3\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	if err.Code == 0 {
		t.Errorf("Error did not happen!")
	}
}

func Test_ParseCommand_Should_Return_Multiple_Functions_When_Passed_Multiple_Commands(t *testing.T) {
	parser := RESPParser{}
	commands, err := parser.ParseCommand(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	if err.Code != 0 {
		t.Errorf("Unexpected error happened! %v", err)
	}
	if len(commands) != 2 {
		t.Errorf("Unexpected len for commands! %d", len(commands))
	}
}
