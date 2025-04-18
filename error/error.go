package error

import (
	"fmt"
	"log/slog"
)

var (
	KeyNotFoundInDictionary        = Error{"Key not found in dictionary", "Key not found in dictionary", 1, nil, make(map[string]string)}
	IndexOutOfRange                = Error{"Index set is out of range", "Index set is out of range", 2, nil, make(map[string]string)}
	UnableToReadFirstByte          = Error{"Unable to read first byte", "", 3, nil, make(map[string]string)}
	UnexpectedFirstByte            = Error{"First byte was different from expected", "Wrong format for command", 4, nil, make(map[string]string)}
	UnableToFindPattern            = Error{"Unable to find byte pattern in byte stream", "", 5, nil, make(map[string]string)}
	UnableToDetermineBulkArraySize = Error{"Unable to determine the size of the incoming bulk array", "", 6, nil, map[string]string{}}
	UnableToReadBytes              = Error{"Unable to read the specified number of bytes", "", 7, nil, make(map[string]string)}
	InsufficientLength             = Error{"Insufficient length for command", "Command malformed", 8, nil, make(map[string]string)}
	UnableToConvertIndexToInt      = Error{"Unable to convert the provided index to an integer", "", 9, nil, make(map[string]string)}
	FunctionNotFound               = Error{"Function provided not found for current implementation", "Command not found", 9, nil, make(map[string]string)}
	UnexpectedBytes                = Error{"Unexpected bytes encountered while processing byte sequence", "", 10, nil, make(map[string]string)}
	ErrorReceived                  = Error{"Error received as a response", "", 11, nil, make(map[string]string)}
	UnableToConvertLenToInt        = Error{"Unable to convert the given response to an integer representing length of array", "", 12, nil, make(map[string]string)}
)

type Error struct {
	Content       string
	ClientContext string
	Code          uint16
	From          error
	ExtraContext  map[string]string
}

func (e Error) Error() string {
	return fmt.Sprintf("[MiniRedisError-%d] %v\n", e.Code, e.Content)
}

func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("MiniRedisError", e.Content),
		slog.Any("From", e.From),
		slog.Int("ErrorCode", int(e.Code)),
		slog.Any("Extra Information", e.ExtraContext),
	)
}

func Unwrap(e Error) error {
	return e.From
}
