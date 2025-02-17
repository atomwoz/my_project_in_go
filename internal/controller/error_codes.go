package controller

// Struct to hold error codes, not magic numbers
// Manually written hex'es to maintain readability and grep'ability

var ERRORS = struct {
	ErrCodeInternalDatabase    int
	ErrCodeSwiftRecordNotFound int
	ErrCodeCountryCodeNotFound int
}{
	ErrCodeInternalDatabase:    0x1,
	ErrCodeSwiftRecordNotFound: 0x2,
	ErrCodeCountryCodeNotFound: 0x3,
}
