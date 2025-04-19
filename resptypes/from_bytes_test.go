package resptypes

// func Test_miniRedisBlobStringFromBytes_Should_Convert_When_Passed_Valid_Input(t *testing.T) {
// 	stream := bufio.NewReader(bytes.NewReader([]byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r', '\n'}))

// 	s, err := BlobStringFromBytes(stream)
// 	if err.Code != 0 {
// 		t.Errorf("Unexpected error encountered! %v", err)
// 	}
// 	if s != "a samplë" {
// 		t.Errorf("Unable to obtain string from bytes! %v", s)
// 	}

// }

// func Test_miniRedisBlobStringFromBytes_Should_Return_Error_When_Passed_Invalid_Input(t *testing.T) {
// 	stream := bufio.NewReader(bytes.NewReader([]byte{'$', '9', '\r', '\n', 'a', ' ', 's', 'a', 'm', 'p', 'l', 195, 171, '\r'}))
// 	_, err := BlobStringFromBytes(stream)
// 	if err.Code == 0 {
// 		t.Errorf("Error did not happen!")
// 	}
// }
