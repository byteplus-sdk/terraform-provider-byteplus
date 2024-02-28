// Package query provides serialization of BYTEPLUS byteplusquery requests, and responses.
package query

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"net/url"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
	"github.com/byteplus-sdk/byteplus-go-sdk/private/protocol/query/queryutil"
)

// BuildHandler is a named request handler for building byteplusquery protocol requests
var BuildHandler = request.NamedHandler{Name: "awssdk.byteplusquery.Build", Fn: Build}

// Build builds a request for an BYTEPLUS Query service.
func Build(r *request.Request) {
	body := url.Values{
		"Action":  {r.Operation.Name},
		"Version": {r.ClientInfo.APIVersion},
	}
	if err := queryutil.Parse(body, r.Params, false); err != nil {
		r.Error = bytepluserr.New(request.ErrCodeSerialization, "failed encoding Query request", err)
		return
	}

	if !r.IsPresigned() {
		r.HTTPRequest.Method = "POST"
		r.HTTPRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		r.SetBufferBody([]byte(body.Encode()))
	} else { // This is a pre-signed request
		r.HTTPRequest.Method = "GET"
		r.HTTPRequest.URL.RawQuery = body.Encode()
	}
}
