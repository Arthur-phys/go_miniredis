package parser

import (
	"fmt"
	"io"
	"miniredis/core/coreinterface"
	e "miniredis/error"
	"strconv"
)

type RESPParser struct{}

func (r *RESPParser) ParseCommand(b []byte) ([]func(d coreinterface.CacheStore) ([]byte, e.Error), e.Error) {
	commands := []func(d coreinterface.CacheStore) ([]byte, e.Error){}
	stream := NewStream(b)
	var internalParser func() e.Error
	internalParser = func() e.Error {
		var newErr e.Error
		firstByte, err := stream.TakeOne()
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
		bytesRead, err := stream.ReadUntilSliceFound([]byte{'\r', '\n'})
		if err != nil {
			newErr = e.UnableToFindPattern
			newErr.From = err
			newErr.ExtraContext["pattern"] = `\r\n`
			return newErr
		}
		i, err := strconv.Atoi(string(bytesRead))
		if err != nil {
			newErr = e.UnableToDetermineBulkArraySize
			newErr.From = err
			return newErr
		}
		arr := make([]string, i)
		for j := range arr {
			arr[j], newErr = r.miniRedisBlobStringFromBytes(&stream)
			if newErr.Code != 0 {
				return newErr
			}
		}
		f, newErr := selectFunction(arr)
		if newErr.Code != 0 {
			return newErr
		}
		commands = append(commands, f)
		return internalParser()
	}

	newErr := internalParser()
	return commands, newErr
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
				return BlobStringToRESP(val), e.Error{}
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
			return NullToRESP(), e.Error{}
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
			return NullToRESP(), e.Error{}
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
			return BlobStringToRESP(val), e.Error{}
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
			return NullToRESP(), e.Error{}
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
			return BlobStringToRESP(val), e.Error{}
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
			return IntToRESP(val), e.Error{}
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
			return BlobStringToRESP(val), e.Error{}
		}, e.Error{}
	default:
		newErr := e.FunctionNotFound
		newErr.ExtraContext["function"] = arr[0]
		return f, newErr
	}
}

func (r *RESPParser) miniRedisBlobStringFromBytes(st *Stream) (string, e.Error) {
	firstByte, err := st.TakeOne()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return "", newErr
	}
	if firstByte != '$' {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = "$"
		newErr.ExtraContext["received"] = string(firstByte)
		return "", newErr
	}
	bytesRead, err := st.ReadUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		newErr := e.UnableToFindPattern
		newErr.From = err
		newErr.ExtraContext["pattern"] = `\r\n`
		return "", newErr
	}
	long, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		newErr := e.UnableToDetermineBulkArraySize
		newErr.From = err
		return "", newErr
	}
	blobString, _, err := st.ReadNBytes(long)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", newErr
	}
	_, err = st.Skip(2)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", newErr
	}
	return string(blobString), e.Error{}
}

func NewRESPParser() coreinterface.Parser {
	return &RESPParser{}
}
