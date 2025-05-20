package respparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/Arthur-phys/redigo/pkg/core/interfaces"
	"github.com/Arthur-phys/redigo/pkg/core/tobytes"
	e "github.com/Arthur-phys/redigo/pkg/error"
)

// RESPParser is responsible for holding all the methods to abstract
// the extraction of commands out of a byte stream (a tcp connection).
//
// This means that instead of reading a raw connection or buffer to parse types,
// RESPParser can be passed the connection an will do so for you with one of it's methods.
//
// Be aware: not all parser operations for RESP types are implemented, only those needed. This behaviour should
// be easy to extend with a simple new function for the given type.
//
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
type RESPParser struct {
	conn                   *net.Conn
	rawBuffer              []byte
	rawBufferPosition      int
	rawBufferEffectiveSize int
	totalBytesRead         int
	messageSizeLimit       int
	buffer                 *bufio.Reader
	lastCommand            []byte
	lastCommandUnprocessed bool
}

func New(conn *net.Conn, maxBytesAllowed int) *RESPParser {
	return &RESPParser{conn, []byte{}, 0, 0, 0, maxBytesAllowed, &bufio.Reader{}, []byte{}, false}
}

func (r *RESPParser) NewConnection(conn *net.Conn) {
	r.conn = conn
	r.rawBuffer = []byte{}
	r.rawBufferPosition = 0
	r.rawBufferEffectiveSize = 0
	r.totalBytesRead = 0
	r.buffer = &bufio.Reader{}
	r.lastCommand = []byte{}
	r.lastCommandUnprocessed = false
}

// Read reads from the connection and into an internal bufio.Reader taking into account possible
// previous commands not parsed.
func (r *RESPParser) Read() (int, error) {
	// Read in chunks of 4 kilobytes
	r.rawBuffer = make([]byte, 4096)
	r.rawBufferPosition = 0
	n, err := (*r.conn).Read(r.rawBuffer)
	if err != nil {
		redigoError := e.UnableToReadFromConnection
		redigoError.From = err
		return n, redigoError
	}
	// From 4096 bytes, how many are actually non empty?
	r.rawBufferEffectiveSize = n

	// Whenever the last command was not complete, add that to the buffer
	if r.lastCommandUnprocessed {
		r.totalBytesRead += n
		r.rawBuffer = append(r.lastCommand, r.rawBuffer[:n]...)
		r.rawBufferEffectiveSize += len(r.lastCommand)
		// Make sure to remove the last command and also the unprocessed flag!
		r.lastCommand = []byte{}
		r.lastCommandUnprocessed = false
		// The internal buffer will be correctly sized and passed into a bufio.Reader
		r.buffer.Reset(bytes.NewReader(r.rawBuffer))
	} else {
		// Otherwise everything was normal
		r.totalBytesRead = n
		r.buffer.Reset(bytes.NewReader(r.rawBuffer[:n]))
	}

	// Did we exceed the size for a single message?
	// This takes into account any unprocessed commands sent previously
	if r.totalBytesRead > r.messageSizeLimit {
		r.lastCommand = []byte{}
		r.lastCommandUnprocessed = false
		r.totalBytesRead = 0
		redigoError := e.MaxSizePerCallExceeded
		redigoError.ExtraContext["maxSize"] = fmt.Sprintf("%d", r.messageSizeLimit)
		redigoError.ExtraContext["currentSize"] = fmt.Sprintf("%d", r.totalBytesRead)
		return n, redigoError
	}
	return n, nil
}

// ParseCommand will use the RESPParser to parse as many commands as possible from the given internal buffer.
//
// It returns all commands able to be parsed at once to the client, incluiding any errors.
func (r *RESPParser) ParseCommand() ([]func(d interfaces.CacheStore) ([]byte, error), error) {
	var (
		// To create the array of strings this function needs to call itself
		internalParser func() error
		commands       []func(d interfaces.CacheStore) ([]byte, error)
	)

	internalParser = func() error {
		blobStrings, n, err := ParseArray(r, func(r *RESPParser) (string, int, error) {
			return r.ParseBlobString()
		})
		if n == 0 {
			return err
		}
		// If no command was able to be formed, this means the only command available is incomplete,
		// return and try putting more into the buffer
		if err != nil && blobStrings == nil {
			r.lastCommand = r.rawBuffer[r.rawBufferPosition:r.rawBufferEffectiveSize]
			r.lastCommandUnprocessed = true
			return err
		}
		// Now for every blobString array representing a command, we select the function and
		// Call the parser again
		r.rawBufferPosition += n
		f, err := selectFunction(blobStrings)
		if err != nil {
			return err
		}
		commands = append(commands, f)
		// Now go for the next command in the same buffer
		return internalParser()
	}

	err := internalParser()
	return commands, err
}

// ParseArray is recursive and uses any of the other Parse functions to create an array of that type.
//
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
func ParseArray[T any](r *RESPParser, transformer func(r *RESPParser) (T, int, error)) ([]T, int, error) {
	// Every function returns the total amount read in case it is necessary for whom it calls it
	var totalBytesRead int

	err := r.checkFirstByte('*')
	if err != nil {
		return nil, totalBytesRead, err
	}
	totalBytesRead += 1

	// Determines the size of the array by converting the given string into a number
	num, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if err != nil {
		return nil, totalBytesRead, err
	}
	i, err := strconv.Atoi(string(num))
	if err != nil {
		redigoError := e.UnableToDetermineBulkArraySize
		redigoError.From = err
		return nil, totalBytesRead, redigoError
	}

	// Apply the transformer to every element of the byteStream array to obtain an array of T
	arr := make([]T, i)
	for j := range arr {
		var m int = 0
		arr[j], m, err = transformer(r)
		totalBytesRead += m
		if err != nil {
			return nil, totalBytesRead, err
		}
	}
	return arr, totalBytesRead, nil
}

// ParseBlobString uses RESP Protocol to convert bytes into a string.
//
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
func (r *RESPParser) ParseBlobString() (string, int, error) {
	totalBytesRead := 0

	err := r.checkFirstByte('$')
	if err != nil {
		return "", totalBytesRead, err
	}
	totalBytesRead += 1

	bytesArr, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if err != nil {
		return "", totalBytesRead, err
	}

	long, err := strconv.Atoi(string(bytesArr))
	if err != nil {
		redigoError := e.UnableToDetermineRawStringSize
		redigoError.From = err
		return "", totalBytesRead, redigoError
	}

	blobString := make([]byte, long)
	n, err = io.ReadFull(r.buffer, blobString)
	totalBytesRead += n
	if err != nil {
		redigoError := e.UnableToReadBytes
		redigoError.From = err
		return "", totalBytesRead, redigoError
	}

	n, err = r.buffer.Discard(2)
	totalBytesRead += n
	if err != nil {
		redigoError := e.UnableToReadBytes
		redigoError.From = err
		return "", totalBytesRead, redigoError
	}
	return string(blobString), totalBytesRead, nil
}

// ParseNull uses RESP Protocol to convert Null response into an empty Error
//
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
func (r *RESPParser) ParseNull() (int, error) {
	totalBytesRead := 0

	err := r.checkFirstByte('_')
	if err != nil {
		return totalBytesRead, err
	}
	totalBytesRead += 1

	_, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}
	if n != 2 {
		return totalBytesRead, e.NotNullFoundInPlaceOfNull
	}

	return totalBytesRead, nil
}

// ParseUInt uses RESP Protocol to convert bytes into an int.
//
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
func (r *RESPParser) ParseUInt() (int, int, error) {
	totalBytesRead := 0

	err := r.checkFirstByte(':')
	if err != nil {
		return 0, totalBytesRead, err
	}
	totalBytesRead += 1

	integerReceived, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if err != nil {
		return 0, totalBytesRead, err
	}

	num, err := strconv.Atoi(string(integerReceived))
	if err != nil {
		redigoError := e.UnableToConvertLenToInt
		redigoError.From = err
		return 0, totalBytesRead, redigoError
	}
	return num, totalBytesRead, nil
}

// ParseError uses RESP Protocol to convert an error into an Error
// See RESP protocol
// https://github.com/redis/redis-specifications/blob/master/protocol/RESP3.md
func (r *RESPParser) ParseError() (int, error) {
	totalBytesRead := 0

	err := r.checkFirstByte('-')
	if err != nil {
		return totalBytesRead, err
	}
	totalBytesRead += 1
	errorReceived, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		return totalBytesRead, err
	}
	totalBytesRead += n
	finalErr := e.ErrorReceived
	finalErr.ExtraContext["text"] = string(errorReceived)
	return totalBytesRead, finalErr
}

func (r *RESPParser) checkFirstByte(b byte) error {
	firstByte, err := r.buffer.ReadByte()
	if err != nil {
		redigoError := e.UnableToReadFirstByte
		redigoError.From = err
		return redigoError
	}
	if firstByte != b {
		err := r.buffer.UnreadByte()
		redigoError := e.UnexpectedFirstByte
		redigoError.From = err
		redigoError.ExtraContext["expected"] = string(b)
		redigoError.ExtraContext["received"] = string(firstByte)
		return redigoError
	}
	return nil
}

// readUntilSliceFound is a helper function to recursively read a buffer until finding a chain of bytes.
// All of them have to match (in order & presence) for the function to return a value satisfactory
func (r *RESPParser) readUntilSliceFound(delim []byte) ([]byte, int, error) {
	var sliceFoundRecursive func([]byte, []byte) ([]byte, error)

	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, error) {
		bytes, err := r.buffer.ReadBytes(delim[0])
		if err != nil {
			redigoError := e.UnableToFindPattern
			redigoError.From = err
			redigoError.ExtraContext["pattern"] = string(delim)
			return bytesRead, redigoError
		}
		bytesRead = append(bytesRead, bytes...)

		for i := 1; i < len(delim); i++ {
			newByte, err := r.buffer.ReadByte()
			if err != nil {
				redigoError := e.UnableToFindPattern
				redigoError.From = err
				return bytesRead, redigoError
			}
			bytesRead = append(bytesRead, newByte)
			if newByte != delim[i] {
				return sliceFoundRecursive(delim, bytesRead)
			}
		}
		return bytesRead, nil
	}
	bytes, err := sliceFoundRecursive(delim, []byte{})
	totalBytesRead := len(bytes)
	if err == nil {
		bytes = bytes[:len(bytes)-len(delim)]
	}
	return bytes, totalBytesRead, err
}

// selectFunction will read an array of strings and return a command to be run on the cache.
//
// Here's where you would implement a new command.
func selectFunction(arr []string) (func(d interfaces.CacheStore) ([]byte, error), error) {
	var f func(d interfaces.CacheStore) ([]byte, error)
	var redigoError e.Error
	switch arr[0] {
	case "GET":
		if len(arr) != 2 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "2"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			if val, err := d.Get(arr[1]); err == nil {
				return tobytes.BlobString(val), nil
			} else if e.KeyNotFound(err) {
				return tobytes.Null(), nil
			} else {
				return []byte{}, err
			}
		}, nil
	case "SET":
		if len(arr) != 3 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "3"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			err := d.Set(arr[1], arr[2])
			if err != nil {
				return []byte{}, err
			}
			return tobytes.Null(), nil
		}, nil
	case "RPUSH":
		if len(arr) < 3 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = ">= 3"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			err := d.RPush(arr[1], arr[2:]...)
			if err != nil {
				return []byte{}, err
			}
			return tobytes.Null(), nil
		}, nil
	case "RPOP":
		if len(arr) != 2 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "2"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			val, err := d.RPop(arr[1])
			if err == nil {
				return tobytes.BlobString(val), nil
			} else if e.KeyNotFound(err) {
				return tobytes.Null(), nil
			}
			return []byte{}, err
		}, nil
	case "LPUSH":
		if len(arr) < 3 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "> 3"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			err := d.LPush(arr[1], arr[2:]...)
			if err != nil {
				return []byte{}, err
			}
			return tobytes.Null(), nil
		}, nil
	case "LPOP":
		if len(arr) != 2 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "2"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			val, err := d.LPop(arr[1])
			if err == nil {
				return tobytes.BlobString(val), nil
			} else if e.KeyNotFound(err) {
				return tobytes.Null(), nil
			}
			return []byte{}, err
		}, nil
	case "LLEN":
		if len(arr) != 2 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "2"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			val, err := d.LLen(arr[1])
			if err != nil {
				return []byte{}, err
			}
			return tobytes.Int(val), nil
		}, nil
	case "LINDEX":
		if len(arr) != 3 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "3"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			index, err := strconv.Atoi(arr[2])
			if err != nil {
				redigoError := e.UnableToConvertIndexToInt
				redigoError.From = err
				redigoError.ExtraContext["provided"] = arr[2]
				return []byte{}, redigoError
			}
			val, err := d.LIndex(arr[1], index)
			if e.KeyNotFound(err) || e.IndexOutOfRange(err) {
				return tobytes.Null(), nil
			} else if err != nil {
				return []byte{}, err
			}
			return tobytes.BlobString(val), nil
		}, nil
	case "DEL":
		if len(arr) != 2 {
			redigoError = e.InsufficientLength
			redigoError.ExtraContext["expected"] = "2"
			redigoError.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, redigoError
		}
		return func(d interfaces.CacheStore) ([]byte, error) {
			err := d.Del(arr[1])
			if err != nil {
				return []byte{}, err
			}
			return tobytes.Null(), nil
		}, nil
	case "PING":
		return func(d interfaces.CacheStore) ([]byte, error) {
			return tobytes.Pong(), nil
		}, nil
	default:
		redigoError := e.FunctionNotFound
		redigoError.ExtraContext["function"] = arr[0]
		return f, redigoError
	}
}
