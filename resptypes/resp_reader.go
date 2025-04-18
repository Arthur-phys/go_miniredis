package resptypes

// // Simple wrapper for a bufio.Reader
// type RESPReader struct {
// 	byteReader *bufio.Reader
// 	buffer     []byte
// }

// func NewRESPReader(c *net.Conn) RESPReader {
// 	return RESPReader{bufio.NewReader(*c), make([]byte, 4096)}
// }

// func (s *RESPReader) ReadUntilSliceFound(delim []byte) ([]byte, error) {
// 	if len(delim) == 0 {
// 		return []byte{}, e.Error{} // Change
// 	}
// 	var sliceFoundRecursive func([]byte, []byte) ([]byte, error)
// 	sliceFoundRecursive = func(delim []byte, bytesRead []byte) ([]byte, error) {
// 		bytes, err := s.byteReader.ReadBytes(delim[0])
// 		bytesRead = append(bytesRead, bytes...)
// 		if err != nil {
// 			return bytesRead, err
// 		}
// 		for i := 1; i < len(delim); i++ {
// 			newByte, err := s.TakeOne()
// 			if err != nil {
// 				return bytesRead, err
// 			}
// 			bytesRead = append(bytesRead, newByte)
// 			if newByte != delim[i] {
// 				return sliceFoundRecursive(delim, bytesRead) // Change
// 			}
// 		}
// 		return bytesRead, nil
// 	}
// 	bytes, err := sliceFoundRecursive(delim, []byte{})
// 	if err == nil {
// 		bytes = bytes[:len(bytes)-len(delim)]
// 	}
// 	return bytes, err
// }

// func (s *RESPReader) ReadNBytes(n int) ([]byte, int, error) {
// 	readBytes := make([]byte, n)
// 	m, err := io.ReadFull(s.byteReader, readBytes)
// 	return readBytes, m, err
// }

// func (s *RESPReader) TakeOne() (byte, error) {
// 	if _, err := s.byteReader.Peek(1); err != nil {
// 		return 0, err
// 	}
// 	return s.byteReader.ReadByte()
// }
