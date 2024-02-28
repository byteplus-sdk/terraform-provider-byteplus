package byteplusquery

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"net/url"
	"strings"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/byteplusbody"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
)

// BuildHandler is a named request handler for building byteplusquery protocol requests
var BuildHandler = request.NamedHandler{Name: "byteplussdk.byteplusquery.Build", Fn: Build}

// Build builds a request for a byteplus Query service.
func Build(r *request.Request) {
	body := url.Values{
		"Action":  {r.Operation.Name},
		"Version": {r.ClientInfo.APIVersion},
	}
	//r.HTTPRequest.Header.Add("Accept", "application/json")
	//method := strings.ToUpper(r.HTTPRequest.Method)

	if r.Config.ExtraUserAgent != nil && *r.Config.ExtraUserAgent != "" {
		if strings.HasPrefix(*r.Config.ExtraUserAgent, "/") {
			request.AddToUserAgent(r, *r.Config.ExtraUserAgent)
		} else {
			request.AddToUserAgent(r, "/"+*r.Config.ExtraUserAgent)
		}

	}
	r.HTTPRequest.Host = r.HTTPRequest.URL.Host
	v := r.HTTPRequest.Header.Get("Content-Type")
	if (strings.ToUpper(r.HTTPRequest.Method) == "PUT" ||
		strings.ToUpper(r.HTTPRequest.Method) == "POST" ||
		strings.ToUpper(r.HTTPRequest.Method) == "DELETE" ||
		strings.ToUpper(r.HTTPRequest.Method) == "PATCH") &&
		strings.Contains(strings.ToLower(v), "application/json") {
		r.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")
		byteplusbody.BodyJson(&body, r)
	} else {
		byteplusbody.BodyParam(&body, r)
	}
}
