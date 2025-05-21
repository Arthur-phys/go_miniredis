package error

import (
	"fmt"
	"io"
	"log/slog"
)

var (
	KeyNotFoundInDictionary = Error{"Key not found in dictionary", "", 1, nil, make(map[string]string)}
	IndexOutOfRangeErr      = Error{"Index set is out of range", "", 2, nil, make(map[string]string)}

	// Not cheked errors
	UnableToReadFirstByte          = Error{"Unable to read first byte", "", 3, nil, make(map[string]string)}
	UnableToFindPattern            = Error{"Unable to find byte pattern in byte stream", "", 4, nil, make(map[string]string)}
	UnexpectedFirstByte            = Error{"First byte was different from expected", "Command malformed", 5, nil, make(map[string]string)}
	UnableToDetermineBulkArraySize = Error{"Unable to determine the size of the incoming bulk array", "Command malformed", 6, nil, map[string]string{}}
	UnableToDetermineRawStringSize = Error{"Unable to determine the size of the incoming raw string", "Command malformed", 7, nil, map[string]string{}}
	UnableToReadBytes              = Error{"Unable to read the specified number of bytes", "Command malformed", 8, nil, make(map[string]string)}
	InsufficientLength             = Error{"Insufficient length for command", "Command malformed", 9, nil, make(map[string]string)}
	FunctionNotFound               = Error{"Function provided not found for current implementation", "Command not found", 10, nil, make(map[string]string)}
	UnableToConvertIndexToInt      = Error{"Unable to convert the provided index to an integer", "", 11, nil, make(map[string]string)}
	NotNullFoundInPlaceOfNull      = Error{"Null-like stream processed with not null received (len of content is bigger than 2)", "", 12, nil, make(map[string]string)}
	ErrorReceived                  = Error{"Error received as a response", "", 13, nil, make(map[string]string)}
	UnableToConvertLenToInt        = Error{"Unable to convert the given response to an integer representing length of array", "", 14, nil, make(map[string]string)}
	UnableToSendRequestToServer    = Error{"Unable to send request to miniredis server", "", 16, nil, make(map[string]string)}
	MaxSizePerCallExceeded         = Error{"Max size per call exceeded the marked threshold", "Call exceeded size allowed", 17, nil, make(map[string]string)}
	WrongType                      = Error{"Operation against a key holding the wrong kind of value", "Operation against a key holding the wrong kind of value", 18, nil, make(map[string]string)}
	UnableToCreateServer           = Error{"Unable to create the redigo server", "", 19, nil, make(map[string]string)}
)

type Error struct {
	Content       string
	ClientContext string
	Code          uint16
	From          error
	ExtraContext  map[string]string
}

func (e Error) Error() string {
	return fmt.Sprintf("{CODE: %d -- CONTENT: %v -- FROM: %e}", e.Code, e.Content, e.From)
}

func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ERROR", e.Content),
		slog.Any("FROM", e.From),
		slog.Int("CODE", int(e.Code)),
		slog.Any("INFORMATION", e.ExtraContext),
	)
}

func ConnectionRelated(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

func IndexOutOfRange(e error) bool {
	err, ok := e.(Error)
	return err.Code == 2 && ok
}

func KeyNotFound(e error) bool {
	err, ok := e.(Error)
	return err.Code == 1 && ok
}

func ExceededMaxSize(e error) bool {
	err, ok := e.(Error)
	return err.Code == 17 && ok
}

func BufferExhausted(e error) bool {
	err, ok := e.(Error)
	return err.Code == 0 || err.Code == 3 || err.Code == 4 || err.Code == 8 && ok
}
