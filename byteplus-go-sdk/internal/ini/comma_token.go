package ini

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

var commaRunes = []rune(",")

func isComma(b rune) bool {
	return b == ','
}

func newCommaToken() Token {
	return newToken(TokenComma, commaRunes, NoneType)
}
