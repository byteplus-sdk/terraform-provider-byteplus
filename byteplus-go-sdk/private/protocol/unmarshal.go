package protocol

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"io"
	"io/ioutil"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
)

// UnmarshalDiscardBodyHandler is a named request handler to empty and close a response's byteplusbody
var UnmarshalDiscardBodyHandler = request.NamedHandler{Name: "byteplussdk.shared.UnmarshalDiscardBody", Fn: UnmarshalDiscardBody}

// UnmarshalDiscardBody is a request handler to empty a response's byteplusbody and closing it.
func UnmarshalDiscardBody(r *request.Request) {
	if r.HTTPResponse == nil || r.HTTPResponse.Body == nil {
		return
	}

	io.Copy(ioutil.Discard, r.HTTPResponse.Body)
	r.HTTPResponse.Body.Close()
}
