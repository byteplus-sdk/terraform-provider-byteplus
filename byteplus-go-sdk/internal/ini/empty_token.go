package ini

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

// emptyToken is used to satisfy the Token interface
var emptyToken = newToken(TokenNone, []rune{}, NoneType)
