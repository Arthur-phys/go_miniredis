package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"miniredis/core/coreinterface"
	e "miniredis/error"
	rt "miniredis/resptypes"
	"net"
	"strconv"
)

type RESPParser struct {
	conn                   *net.Conn
	rawBuffer              []byte
	rawBufferPosition      int
	rawBufferMaxBytes      int
	buffer                 *bufio.Reader
	lastCommand            []byte
	lastCommandUnprocessed bool
}

func NewRESPParser(conn *net.Conn) *RESPParser {
	rawBuffer := make([]byte, 4096)
	lastCommand := []byte{}
	buffer := bufio.NewReader(bytes.NewReader(rawBuffer))
	return &RESPParser{conn, rawBuffer, 0, 4096, buffer, lastCommand, false}
}

func (r *RESPParser) Read() (int, e.Error) {
	r.rawBuffer = make([]byte, 4096)
	r.rawBufferPosition = 0
	n, err := (*r.conn).Read(r.rawBuffer)
	r.rawBufferMaxBytes = n
	if err != nil {
		newErr := e.UnableToReadFromConnection
		newErr.From = err
		return n, newErr
	}
	if r.lastCommandUnprocessed {
		r.rawBuffer = append(r.lastCommand, r.rawBuffer[:n]...)
		r.lastCommandUnprocessed = false
		r.buffer.Reset(bytes.NewReader(r.rawBuffer))
	} else {
		r.buffer.Reset(bytes.NewReader(r.rawBuffer[:n]))
	}
	return n, e.Error{}
}

func (r *RESPParser) ParseCommand() ([]func(d coreinterface.CacheStore) ([]byte, e.Error), e.Error) {
	commands := []func(d coreinterface.CacheStore) ([]byte, e.Error){}
	var internalParser func() e.Error

	internalParser = func() e.Error {
		strArr, n, err := ParseArray(r, func(r *RESPParser) (string, int, e.Error) {
			return r.BlobStringFromBytes()
		})
		if err.Code != 0 && strArr == nil {
			r.lastCommand = r.rawBuffer[r.rawBufferPosition:r.rawBufferMaxBytes]
			r.lastCommandUnprocessed = true
			return err
		}
		r.rawBufferPosition += n
		f, newErr := selectFunction(strArr)
		if newErr.Code != 0 {
			return newErr
		}
		commands = append(commands, f)
		return internalParser()
	}

	newErr := internalParser()
	return commands, newErr
}

func ParseArray[T any](r *RESPParser, f func(r *RESPParser) (T, int, e.Error)) ([]T, int, e.Error) {
	totalBytesRead := 0
	var arr []T

	newErr := r.checkFirstByte('*')
	if newErr.Code != 0 {
		return nil, totalBytesRead, newErr
	}
	totalBytesRead += 1

	num, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if newErr.Code != 0 {
		return nil, totalBytesRead, newErr
	}

	i, err := strconv.Atoi(string(num))
	if err != nil {
		newErr = e.UnableToDetermineBulkArraySize
		newErr.From = err
		return nil, totalBytesRead, newErr
	}

	arr = make([]T, i)
	for j := range arr {
		var m int = 0
		arr[j], m, newErr = f(r)
		totalBytesRead += m
		if newErr.Code != 0 {
			return nil, totalBytesRead, newErr
		}
	}

	return arr, totalBytesRead, e.Error{}

}

func (r *RESPParser) BlobStringFromBytes() (string, int, e.Error) {
	totalBytesRead := 0

	newErr := r.checkFirstByte('$')
	if newErr.Code != 0 {
		return "", totalBytesRead, newErr
	}
	totalBytesRead += 1

	bytesArr, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if newErr.Code != 0 {
		return "", totalBytesRead, newErr
	}

	long, err := strconv.Atoi(string(bytesArr))
	if err != nil {
		newErr := e.UnableToDetermineRawStringSize
		newErr.From = err
		return "", totalBytesRead, newErr
	}

	blobString := make([]byte, long)
	n, err = io.ReadFull(r.buffer, blobString)
	totalBytesRead += n
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", totalBytesRead, newErr
	}

	n, err = r.buffer.Discard(2)
	totalBytesRead += n
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", totalBytesRead, newErr
	}

	return string(blobString), totalBytesRead, e.Error{}
}

func (r *RESPParser) ErrorFromBytes() (int, e.Error) {
	totalBytesRead := 0

	newErr := r.checkFirstByte('-')
	if newErr.Code != 0 {
		return totalBytesRead, newErr
	}
	totalBytesRead += 1

	errorReceived, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if newErr.Code != 0 {
		return totalBytesRead, newErr
	}

	finalErr := e.ErrorReceived
	finalErr.ExtraContext["text"] = string(errorReceived)
	return totalBytesRead, finalErr
}

func (r *RESPParser) NullFromBytes() (int, e.Error) {
	totalBytesRead := 0

	newErr := r.checkFirstByte('_')
	if newErr.Code != 0 {
		return totalBytesRead, newErr
	}
	totalBytesRead += 1

	_, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if newErr.Code != 0 {
		return totalBytesRead, newErr
	}
	if n != 2 {
		return totalBytesRead, e.NotNullFoundInPlaceOfNull
	}

	return totalBytesRead, e.Error{}
}

func (r *RESPParser) UIntFromBytes() (int, int, e.Error) {
	totalBytesRead := 0

	newErr := r.checkFirstByte(':')
	if newErr.Code != 0 {
		return 0, totalBytesRead, newErr
	}
	totalBytesRead += 1

	integerReceived, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
	totalBytesRead += n
	if newErr.Code != 0 {
		return 0, totalBytesRead, newErr
	}

	num, tmpErr := strconv.Atoi(string(integerReceived))
	if tmpErr != nil {
		newErr := e.UnableToConvertLenToInt
		newErr.From = tmpErr
		return 0, totalBytesRead, newErr
	}

	return num, totalBytesRead, e.Error{}
}

func (r *RESPParser) checkFirstByte(b byte) e.Error {
	firstByte, err := r.buffer.ReadByte()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return newErr
	}
	if firstByte != b {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = string(b)
		newErr.ExtraContext["received"] = string(firstByte)
		return newErr
	}
	return e.Error{}
}

func (r *RESPParser) readUntilSliceFound(delim []byte) ([]byte, int, e.Error) {
	var sliceFoundRecursive func([]byte, []byte) ([]byte, e.Error)

	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, e.Error) {
		bytes, err := r.buffer.ReadBytes(delim[0])
		if err != nil {
			newErr := e.UnableToFindPattern
			newErr.From = err
			newErr.ExtraContext["pattern"] = string(delim)
			return bytesRead, newErr
		}
		bytesRead = append(bytesRead, bytes...)

		for i := 1; i < len(delim); i++ {
			newByte, err := r.buffer.ReadByte()
			if err != nil {
				newErr := e.UnableToFindPattern
				newErr.From = err
				return bytesRead, newErr
			}
			bytesRead = append(bytesRead, newByte)
			if newByte != delim[i] {
				return sliceFoundRecursive(delim, bytesRead)
			}
		}
		return bytesRead, e.Error{}
	}
	bytes, err := sliceFoundRecursive(delim, []byte{})
	if err.Code == 0 {
		bytes = bytes[:len(bytes)-len(delim)]
	}
	return bytes, len(bytes) + len(delim), err
}

func selectFunction(arr []string) (func(d coreinterface.CacheStore) ([]byte, e.Error), e.Error) {
	var f func(d coreinterface.CacheStore) ([]byte, e.Error)
	var newErr e.Error
	switch arr[0] {
	case "GET":
		if len(arr) != 2 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "2"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			if val, err := d.Get(arr[1]); err.Code == 0 {
				return rt.BlobStringToBytes(val), e.Error{}
			} else {
				return []byte{}, err
			}
		}, e.Error{}
	case "SET":
		if len(arr) != 3 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "3"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			err := d.Set(arr[1], arr[2])
			if err.Code != 0 {
				return []byte{}, err
			}
			return rt.NullToBytes(), e.Error{}
		}, e.Error{}
	case "RPUSH":
		if len(arr) < 3 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = ">= 3"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			err := d.RPush(arr[1], arr[2:]...)
			if err.Code != 0 {
				return []byte{}, err
			}
			return rt.NullToBytes(), e.Error{}
		}, e.Error{}
	case "RPOP":
		if len(arr) != 2 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "2"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			val, err := d.RPop(arr[1])
			if err.Code != 0 {
				return []byte{}, err
			}
			return rt.BlobStringToBytes(val), e.Error{}
		}, e.Error{}
	case "LPUSH":
		if len(arr) < 3 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "> 3"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			err := d.LPush(arr[1], arr[2:]...)
			if err.Code != 0 {
				return []byte{}, err
			}
			return rt.NullToBytes(), e.Error{}
		}, e.Error{}
	case "LPOP":
		if len(arr) != 2 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "2"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			val, err := d.LPop(arr[1])
			if err.Code != 0 {
				return []byte{}, err // Propper error handling
			}
			return rt.BlobStringToBytes(val), e.Error{}
		}, e.Error{}
	case "LLEN":
		if len(arr) != 2 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "2"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			val, err := d.LLen(arr[1])
			if err.Code != 0 {
				return []byte{}, err
			}
			return rt.IntToBytes(val), e.Error{}
		}, e.Error{}
	case "LINDEX":
		if len(arr) != 3 {
			newErr = e.InsufficientLength
			newErr.ExtraContext["expected"] = "3"
			newErr.ExtraContext["obtained"] = fmt.Sprintf("%v", len(arr))
			return f, newErr
		}
		return func(d coreinterface.CacheStore) ([]byte, e.Error) {
			index, err := strconv.Atoi(arr[2])
			if err != nil {
				newErr := e.UnableToConvertIndexToInt
				newErr.From = err
				newErr.ExtraContext["provided"] = arr[2]
				return []byte{}, newErr
			}
			val, newErr := d.LIndex(arr[1], index)
			if newErr.Code != 0 {
				return []byte{}, newErr
			}
			return rt.BlobStringToBytes(val), e.Error{}
		}, e.Error{}
	default:
		newErr := e.FunctionNotFound
		newErr.ExtraContext["function"] = arr[0]
		return f, newErr
	}
}
