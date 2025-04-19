package resptypes

// // Simple wrapper for a bufio.Reader
// type RESPReader struct {
// 	byteReader *bufio.Reader
// 	buffer     []byte
// }

// func NewRESPReader(c *net.Conn) RESPReader {
// 	return RESPReader{bufio.NewReader(*c), make([]byte, 4096)}
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
