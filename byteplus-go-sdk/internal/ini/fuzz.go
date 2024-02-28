//go:build gofuzz
// +build gofuzz

package ini

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"bytes"
)

func Fuzz(data []byte) int {
	b := bytes.NewReader(data)

	if _, err := Parse(b); err != nil {
		return 0
	}

	return 1
}
