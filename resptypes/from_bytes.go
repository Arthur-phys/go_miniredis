package resptypes

import (
	"bufio"
	"io"
	e "miniredis/error"
	"strconv"
)

func BlobStringFromBytes(st *bufio.Reader) (string, e.Error) {
	firstByte, err := st.ReadByte()
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
	bytesRead, err := ReadUntilSliceFound(st, []byte{'\r', '\n'})
	if err != nil {
		newErr := e.UnableToFindPattern
		newErr.From = err
		newErr.ExtraContext["pattern"] = `\n`
		return "", newErr
	}
	long, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		newErr := e.UnableToDetermineBulkArraySize
		newErr.From = err
		return "", newErr
	}
	blobString := make([]byte, long)
	_, err = io.ReadFull(st, blobString)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", newErr
	}
	_, err = st.Discard(2)
	if err != nil {
		newErr := e.UnableToReadBytes
		newErr.From = err
		return "", newErr
	}
	return string(blobString), e.Error{}
}

func ErrorFromBytes(st *bufio.Reader) e.Error {
	firstByte, err := st.ReadByte()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return newErr
	}
	if firstByte == '_' {
		restOfResponse := make([]byte, 2)
		_, err = io.ReadFull(st, restOfResponse)
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
		bytesRead, err := ReadUntilSliceFound(st, []byte{'\r', '\n'})
		if err != nil {
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

func UIntFromBytes(st *bufio.Reader) (int, e.Error) {
	firstByte, err := st.ReadByte()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return 0, newErr
	}
	if firstByte != ':' {
		newErr := e.UnexpectedFirstByte
		newErr.ExtraContext["expected"] = ":"
		newErr.ExtraContext["received"] = string(firstByte)
		return 0, newErr
	}
	bytesRead, err := ReadUntilSliceFound(st, []byte{'\r', '\n'})
	if err != nil {
		newErr := e.UnableToFindPattern
		newErr.From = err
		newErr.ExtraContext["pattern"] = `\r\n`
		return 0, newErr
	}
	num, err := strconv.Atoi(string(bytesRead))
	if err != nil {
		newErr := e.UnableToConvertLenToInt
		newErr.From = err
		return 0, newErr
	}
	return num, e.Error{}
}

func ReadUntilSliceFound(buffer *bufio.Reader, delim []byte) ([]byte, error) {
	if len(delim) == 0 {
		return []byte{}, e.Error{} // Change
	}
	var sliceFoundRecursive func([]byte, []byte) ([]byte, error)
	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, error) {
		bytes, err := buffer.ReadBytes(delim[0])
		bytesRead = append(bytesRead, bytes...)
		if err != nil {
			return bytesRead, err
		}
		for i := 1; i < len(delim); i++ {
			newByte, err := buffer.ReadByte()
			if err != nil {
				return bytesRead, err
			}
			bytesRead = append(bytesRead, newByte)
			if newByte != delim[i] {
				return sliceFoundRecursive(delim, bytesRead) // Change
			}
		}
		return bytesRead, nil
	}
	bytes, err := sliceFoundRecursive(delim, []byte{})
	if err == nil {
		bytes = bytes[:len(bytes)-len(delim)]
	}
	return bytes, err
}
