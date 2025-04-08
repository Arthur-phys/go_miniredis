package parser

import (
	"fmt"
	"io"
	"miniredis/core/coreinterface"
	e "miniredis/error"
	rt "miniredis/resptypes"
	"strconv"
)

type RESPParser struct{}

func (r *RESPParser) ParseCommand(b []byte) ([]func(d coreinterface.CacheStore) ([]byte, e.Error), e.Error) {
	commands := []func(d coreinterface.CacheStore) ([]byte, e.Error){}
	stream := rt.NewStream(b)
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
			arr[j], newErr = rt.BlobStringFromBytes(&stream)
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

func NewRESPParser() coreinterface.Parser {
	return &RESPParser{}
}
