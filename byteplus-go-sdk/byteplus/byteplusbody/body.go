package byteplusbody

// Copy from https://github.com/aws/aws-sdk-go
// May have been modified by Byteplus.

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/byteplusutil"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/custom"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
	"github.com/byteplus-sdk/byteplus-go-sdk/private/protocol"
	"github.com/byteplus-sdk/byteplus-go-sdk/private/protocol/query/queryutil"
)

func BodyParam(body *url.Values, r *request.Request) {
	var (
		isForm bool
	)
	contentType := r.HTTPRequest.Header.Get("Content-Type")
	newBody := body
	if strings.ToUpper(r.HTTPRequest.Method) == "POST" || (len(contentType) > 0 && strings.Contains(strings.ToLower(contentType), "x-www-form-urlencoded")) {
		isForm = true
		newBody = &url.Values{}
	}

	if !isForm && len(contentType) > 0 {
		r.Error = bytepluserr.New("SerializationError", "not support such content-type", nil)
		return
	}

	if reflect.TypeOf(r.Params) == reflect.TypeOf(&map[string]interface{}{}) {
		m := *(r.Params).(*map[string]interface{})
		for k, v := range m {
			if reflect.TypeOf(v).String() == "string" {
				newBody.Add(k, v.(string))
			} else {
				newBody.Add(k, fmt.Sprintf("%v", v))
			}
		}
	} else if err := queryutil.Parse(*newBody, r.Params, false); err != nil {
		r.Error = bytepluserr.New("SerializationError", "failed encoding Query request", err)
		return
	}

	//extra process
	if r.Config.ExtraHttpParameters != nil {
		extra := r.Config.ExtraHttpParameters(r.Context())
		if extra != nil {
			for k, value := range extra {
				newBody.Add(k, value)
			}
		}
	}
	if r.Config.ExtraHttpParametersWithMeta != nil {
		extra := r.Config.ExtraHttpParametersWithMeta(r.Context(), custom.RequestMetadata{
			ServiceName: r.ClientInfo.ServiceName,
			Version:     r.ClientInfo.APIVersion,
			Action:      r.Operation.Name,
			HttpMethod:  r.Operation.HTTPMethod,
			Region:      *r.Config.Region,
			Request:     r.HTTPRequest,
			RawQuery:    body,
		})
		if extra != nil {
			for k, value := range extra {
				newBody.Add(k, value)
			}
		}
	}

	if isForm {
		r.HTTPRequest.URL.RawQuery = body.Encode()
		r.HTTPRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		r.SetBufferBody([]byte(newBody.Encode()))
		return
	}

	r.Input = byteplusutil.ParameterToMap(body.Encode(), r.Config.LogSensitives,
		r.Config.LogLevel.Matches(byteplus.LogInfoWithInputAndOutput) || r.Config.LogLevel.Matches(byteplus.LogDebugWithInputAndOutput))

	r.HTTPRequest.URL.RawQuery = newBody.Encode()
}

func BodyJson(body *url.Values, r *request.Request) {
	method := strings.ToUpper(r.HTTPRequest.Method)
	if v := r.HTTPRequest.Header.Get("Content-Type"); len(v) == 0 {
		r.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if v := r.HTTPRequest.Header.Get("Content-Type"); !strings.Contains(strings.ToLower(v), "application/json") || method == "GET" {
		return
	}

	input := make(map[string]interface{})

	pt := reflect.ValueOf(r.Params)

	if pt.Kind() == reflect.Ptr {
		value := pt.Elem()
		if value.Kind() == reflect.Struct {
			t := value.Type()
			for i := 0; i < value.NumField(); i++ {
				elemValue := queryutil.ElemOf(value.Field(i))
				field := t.Field(i)
				if field.Name == "ClientToken" && field.Type.Elem().Kind() == reflect.String {
					if !elemValue.IsValid() {
						token := protocol.GetIdempotencyToken()
						value.Field(i).Set(reflect.ValueOf(&token))
					}
				}
			}
		}
	}

	b, _ := json.Marshal(r.Params)

	_ = json.Unmarshal(b, &input)
	if r.Config.ExtraHttpJsonBody != nil {
		r.Config.ExtraHttpJsonBody(r.Context(), &input, custom.RequestMetadata{
			ServiceName: r.ClientInfo.ServiceName,
			Version:     r.ClientInfo.APIVersion,
			Action:      r.Operation.Name,
			HttpMethod:  r.Operation.HTTPMethod,
			Region:      *r.Config.Region,
			Request:     r.HTTPRequest,
			RawQuery:    body,
		})
		b, _ = json.Marshal(input)
	}
	r.SetStringBody(string(b))

	r.HTTPRequest.URL.RawQuery = body.Encode()
	r.IsJsonBody = true

	r.Input = byteplusutil.BodyToMap(input, r.Config.LogSensitives,
		r.Config.LogLevel.Matches(byteplus.LogInfoWithInputAndOutput) || r.Config.LogLevel.Matches(byteplus.LogDebugWithInputAndOutput))
	r.Params = nil
	r.HTTPRequest.Header.Set("Accept", "application/json")
}
