package byteplusquery

//go:generate go run -tags codegen ../../../private/model/cli/gen-protocol-tests ../../../models/protocol_tests/output/ec2.json unmarshal_test.go

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"encoding/xml"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
	"github.com/byteplus-sdk/byteplus-go-sdk/private/protocol/xml/xmlutil"
)

// UnmarshalHandler is a named request handler for unmarshaling ec2query protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "byteplussdk.byteplusquery.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling ec2query protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "byteplussdk.byteplusquery.UnmarshalMeta", Fn: UnmarshalMeta}

// UnmarshalErrorHandler is a named request handler for unmarshaling ec2query protocol request errors
var UnmarshalErrorHandler = request.NamedHandler{Name: "byteplussdk.byteplusquery.UnmarshalError", Fn: UnmarshalError}

// Unmarshal unmarshals a response body for the EC2 protocol.
func Unmarshal(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() {
		decoder := xml.NewDecoder(r.HTTPResponse.Body)
		err := xmlutil.UnmarshalXML(r.Data, decoder, "")
		if err != nil {
			r.Error = bytepluserr.NewRequestFailure(
				bytepluserr.New(request.ErrCodeSerialization,
					"failed decoding EC2 Query response", err),
				r.HTTPResponse.StatusCode,
				r.RequestID,
			)
			return
		}
	}
}

// UnmarshalMeta unmarshals response headers for the EC2 protocol.
func UnmarshalMeta(r *request.Request) {
	r.RequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
	if r.RequestID == "" {
		// Alternative version of request id in the header
		r.RequestID = r.HTTPResponse.Header.Get("X-Amz-Request-Id")
	}
}

type xmlErrorResponse struct {
	XMLName   xml.Name `xml:"Response"`
	Code      string   `xml:"Errors>Error>Code"`
	Message   string   `xml:"Errors>Error>Message"`
	RequestID string   `xml:"RequestID"`
}

// UnmarshalError unmarshals a response error for the EC2 protocol.
func UnmarshalError(r *request.Request) {
	defer r.HTTPResponse.Body.Close()

	var respErr xmlErrorResponse
	err := xmlutil.UnmarshalXMLError(&respErr, r.HTTPResponse.Body)
	if err != nil {
		r.Error = bytepluserr.NewRequestFailure(
			bytepluserr.New(request.ErrCodeSerialization,
				"failed to unmarshal error message", err),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	r.Error = bytepluserr.NewRequestFailure(
		bytepluserr.New(respErr.Code, respErr.Message, nil),
		r.HTTPResponse.StatusCode,
		respErr.RequestID,
	)
}
