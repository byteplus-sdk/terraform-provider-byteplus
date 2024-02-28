package query

// This File is modify from https://github.com/aws/aws-sdk-go/blob/main/private/protocol/query/unmarshal.go

import (
	"encoding/xml"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
	"github.com/byteplus-sdk/byteplus-go-sdk/private/protocol/xml/xmlutil"
)

// UnmarshalHandler is a named request handler for unmarshaling byteplusquery protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "awssdk.byteplusquery.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling byteplusquery protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "awssdk.byteplusquery.UnmarshalMeta", Fn: UnmarshalMeta}

// Unmarshal unmarshals a response for an BYTEPLUS Query service.
func Unmarshal(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() {
		decoder := xml.NewDecoder(r.HTTPResponse.Body)
		err := xmlutil.UnmarshalXML(r.Data, decoder, r.Operation.Name+"Result")
		if err != nil {
			r.Error = bytepluserr.NewRequestFailure(
				bytepluserr.New(request.ErrCodeSerialization, "failed decoding Query response", err),
				r.HTTPResponse.StatusCode,
				r.RequestID,
			)
			return
		}
	}
}

// UnmarshalMeta unmarshals header response values for an BYTEPLUS Query service.
func UnmarshalMeta(r *request.Request) {
	r.RequestID = r.HTTPResponse.Header.Get("X-Top-Requestid")
}
