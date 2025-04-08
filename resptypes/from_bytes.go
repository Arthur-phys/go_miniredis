package resptypes

// More than meets the eye !!!

import (
	e "miniredis/error"
	"strconv"
)

func BlobStringFromBytes(st *Stream) (string, e.Error) {
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

func ErrorFromBytes(st *Stream) e.Error {
	firstByte, err := st.TakeOne()
	if err != nil {
		newErr := e.UnableToReadFirstByte
		newErr.From = err
		return newErr
	}
	if firstByte == '_' {
		restOfResponse, _, err := st.ReadNBytes(2)
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
		bytesRead, err := st.ReadUntilSliceFound([]byte{'\r', '\n'})
		if err != nil {
			newErr := e.UnableToFindPattern
			newErr.From = err
			newErr.ExtraContext["pattern"] = `\r\n`
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

func UIntFromBytes(st *Stream) (int, e.Error) {
	firstByte, err := st.TakeOne()
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
	bytesRead, err := st.ReadUntilSliceFound([]byte{'\r', '\n'})
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
