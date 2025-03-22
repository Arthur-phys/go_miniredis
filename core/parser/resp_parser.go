package parser

import (
	"miniredis/core/coreinterface"
	e "miniredis/error"
	"strconv"
)

type RESPParser struct{}

func (r *RESPParser) ParseCommand(b []byte) (fs []func(d coreinterface.CacheStore) ([]byte, error), err error) {
	stream := NewStream(b)
	commands := make([]func(d coreinterface.CacheStore) ([]byte, error), 1)
	var internalParser func() error
	internalParser = func() error {
		firstByte, err := stream.TakeOne()
		if err != nil {
			return err
		}
		if firstByte != '*' {
			return e.Error{} // Change
		}
		bytesRead, err := stream.ReadUntilSliceFound([]byte{'\r', '\n'})
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(string(bytesRead))
		if err != nil {
			return err
		}
		arr := make([]string, i)
		for j := range arr {
			arr[j], err = r.miniRedisBlobStringFromBytes(&stream)
			if err != nil {
				return err
			}
		}
		f, err := selectFunction(arr)
		if err != nil {
			return err
		}
		commands = append(commands, f)
		return internalParser()
	}

	internalParser()
	return commands, nil
}

func selectFunction(arr []string) (f func(d coreinterface.CacheStore) ([]byte, error), err error) {
	switch arr[0] {
	case "GET":
		if len(arr) != 2 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			if val, ok := d.Get(arr[1]); ok {
				return BlobStringToRESP(val), nil
			} else {
				return []byte{}, e.Error{} //Change
			}
		}, nil
	case "SET":
		if len(arr) != 3 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			err = d.Set(arr[1], arr[2])
			if err != nil {
				return []byte{}, err
			}
			return NullToRESP(), nil
		}, nil
	case "RPUSH":
		if len(arr) < 3 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			err = d.RPush(arr[1], arr[2:]...)
			if err != nil {
				return []byte{}, err //Propper error handling
			}
			return NullToRESP(), nil
		}, nil
	case "RPOP":
		if len(arr) != 2 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			val, err := d.RPop(arr[1])
			if err != nil {
				return []byte{}, err // Propper error handling
			}
			return BlobStringToRESP(val), nil
		}, nil
	case "LPUSH":
		if len(arr) < 3 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			err = d.LPush(arr[1], arr[2:]...)
			if err != nil {
				return []byte{}, err //Propper error handling
			}
			return NullToRESP(), nil
		}, nil
	case "LPOP":
		if len(arr) != 2 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			val, err := d.LPop(arr[1])
			if err != nil {
				return []byte{}, err // Propper error handling
			}
			return BlobStringToRESP(val), nil
		}, nil
	case "LLEN":
		if len(arr) != 2 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			val, err := d.LLen(arr[1])
			if err != nil {
				return []byte{}, err // Propper error handling
			}
			return IntToRESP(val), nil
		}, nil
	case "LINDEX":
		if len(arr) != 3 {
			err = e.Error{}
			return f, err // Change proper error handling
		}
		return func(d coreinterface.CacheStore) ([]byte, error) {
			index, err := strconv.Atoi(arr[2])
			if err != nil {
				return []byte{}, err
			}
			val, b := d.LIndex(arr[1], index)
			if !b {
				return []byte{}, e.Error{} // Propper error handling
			}
			return BlobStringToRESP(val), nil
		}, nil
	default:
		return func(d coreinterface.CacheStore) ([]byte, error) {
			return ErrToRESP(e.Error{}), nil
		}, e.Error{}
	}
}

func (r *RESPParser) miniRedisBlobStringFromBytes(st *Stream) (s string, err error) {
	firstByte, err := st.TakeOne()
	if err != nil {
		return
	}
	if firstByte != '$' {
		return "", e.Error{} // Change
	}
	bytesRead, err := st.ReadUntilSliceFound([]byte{'\r', '\n'})
	if err != nil {
		return
	}
	long, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		return
	}
	blobString, _, err := st.ReadNBytes(long)
	if err != nil {
		return
	}
	_, err = st.Skip(2)
	if err != nil {
		return
	}
	return string(blobString), nil
}

func NewRESPParser() coreinterface.Parser {
	return &RESPParser{}
}
