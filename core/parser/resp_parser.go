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
	buffer                 *bufio.Reader
	lastCommand            []byte
	lastCommandUnprocessed bool
}

func (r *RESPParser) Read() (int, e.Error) {
	r.rawBuffer = make([]byte, 4096)
	r.rawBufferPosition = 0
	n, err := (*r.conn).Read(r.rawBuffer)
	if err != nil {
		newErr := e.UnableToReadFromConnection
		newErr.From = err
		return n, newErr
	}
	if r.lastCommandUnprocessed {
		r.rawBuffer = append(r.lastCommand, r.rawBuffer...)
		r.lastCommandUnprocessed = false
	}
	r.buffer.Reset(bytes.NewReader(r.rawBuffer))
	return n, e.Error{}
}

func (r *RESPParser) ParseCommand() ([]func(d coreinterface.CacheStore) ([]byte, e.Error), e.Error) {
	commands := []func(d coreinterface.CacheStore) ([]byte, e.Error){}
	var internalParser func() e.Error

	internalParser = func() e.Error {
		var newErr e.Error
		firstByte, err := r.buffer.ReadByte()
		if err == io.EOF {
			return e.Error{}
		} else if err != nil {
			newErr = e.UnableToReadFirstByte
			newErr.From = err
			return newErr
		}
		if firstByte != '*' {
			newErr = e.UnexpectedFirstByte
			newErr.ExtraContext["expected"] = "*"
			newErr.ExtraContext["received"] = string(firstByte)
			return newErr
		}

		bytesRead, n, newErr := r.readUntilSliceFound([]byte{'\r', '\n'})
		if newErr.Code == 5 {
			return newErr
		}
		if newErr.Code != 0 {
			newErr = e.UnableToFindPattern
			newErr.From = err
			newErr.ExtraContext["pattern"] = `\r\n`
			r.lastCommand = r.rawBuffer[r.rawBufferPosition:]
			r.lastCommandUnprocessed = true
			return newErr
		}
		i, err := strconv.Atoi(string(bytesRead))
		if err != nil {
			newErr = e.UnableToDetermineBulkArraySize
			newErr.From = err
			return newErr
		}

		arr := make([]string, i)
		rawStringBytesRead := 0
		for j := range arr {
			arr[j], n, newErr = r.BlobStringFromBytes()
			if newErr.Code != 0 {
				r.lastCommand = r.rawBuffer[r.rawBufferPosition:]
				return newErr
			}
			rawStringBytesRead += n
		}
		f, newErr := selectFunction(arr)
		if newErr.Code != 0 {
			return newErr
		}
		commands = append(commands, f)
		r.rawBufferPosition += n + 1 + rawStringBytesRead
		return internalParser()
	}

	newErr := internalParser()
	return commands, newErr
}

func (r *RESPParser) BlobStringFromBytes() (string, int, e.Error) {
	bytesRead := 0
	firstByte, err := r.buffer.ReadByte()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return "", bytesRead, newErr
	}
	if firstByte != '$' {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = "$"
		newErr.ExtraContext["received"] = string(firstByte)
		return "", bytesRead, newErr
	}
	bytesRead += 1
	bytesArr, n, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		newErr := e.UnableToFindPattern
		newErr.From = err
		newErr.ExtraContext["pattern"] = `\r\n`
		return "", bytesRead, newErr
	}
	bytesRead += n
	long, err := strconv.Atoi(string(bytesArr))
	if err != nil {
		newErr := e.UnableToDetermineBulkArraySize
		newErr.From = err
		return "", bytesRead, newErr
	}
	blobString := make([]byte, long)
	_, err = io.ReadFull(r.buffer, blobString)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", bytesRead, newErr
	}
	bytesRead += long
	_, err = r.buffer.Discard(2)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", bytesRead, newErr
	}
	bytesRead += 2
	return string(blobString), bytesRead, e.Error{}
}

func (r *RESPParser) ErrorFromBytes() e.Error {
	firstByte, err := r.buffer.ReadByte()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return newErr
	}
	if firstByte == '_' {
		restOfResponse := make([]byte, 2)
		_, err = io.ReadFull(r.buffer, restOfResponse)
		if err != nil {
			newErr := e.UnableToReadBytes
			newErr.From = err
			return newErr
		}
		if restOfResponse[0] == '\r' && restOfResponse[1] == '\n' {
			return e.Error{}
		} else {
			newErr := e.UnexpectedBytes
			newErr.ExtraContext["expected"] = `\r\n`
			newErr.ExtraContext["received"] = string(restOfResponse)
			return e.UnexpectedBytes
		}
	} else if firstByte == '-' {
		bytesRead, _, err := r.readUntilSliceFound([]byte{'\r', '\n'})
		if err.Code != 0 {
			newErr := e.UnableToFindPattern
			newErr.From = err
			newErr.ExtraContext["pattern"] = `\n`
			return newErr
		}
		finalErr := e.ErrorReceived
		finalErr.ExtraContext["text"] = string(bytesRead)
		return finalErr
	} else {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = "_ OR -"
		newErr.ExtraContext["received"] = string(firstByte)
		return newErr
	}
}

func (r *RESPParser) UIntFromBytes() (int, e.Error) {
	firstByte, tmpErr := r.buffer.ReadByte()
	if tmpErr != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = tmpErr
		return 0, newErr
	}
	if firstByte != ':' {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = ":"
		newErr.ExtraContext["received"] = string(firstByte)
		return 0, newErr
	}
	bytesRead, _, err := r.readUntilSliceFound([]byte{'\r', '\n'})
	if err.Code != 0 {
		newErr := e.UnableToFindPattern
		newErr.From = err
		newErr.ExtraContext["pattern"] = `\r\n`
		return 0, newErr
	}
	num, tmpErr := strconv.Atoi(string(bytesRead))
	if tmpErr != nil {
		newErr := e.UnableToConvertLenToInt
		newErr.From = tmpErr
		return 0, newErr
	}
	return num, e.Error{}
}

func (r *RESPParser) readUntilSliceFound(delim []byte) ([]byte, int, e.Error) {
	var sliceFoundRecursive func([]byte, []byte) ([]byte, e.Error)
	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, e.Error) {
		bytes, err := r.buffer.ReadBytes(delim[0])
		bytesRead = append(bytesRead, bytes...)
		if err != nil {
			newErr := e.UnableToFindPattern
			newErr.From = err
			return bytesRead, newErr
		}
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

func NewRESPParser(conn *net.Conn) *RESPParser {
	rawBuffer := make([]byte, 4096)
	lastCommand := []byte{}
	buffer := bufio.NewReader(bytes.NewReader(rawBuffer))
	return &RESPParser{conn, rawBuffer, 0, buffer, lastCommand, false}
}
