package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/client"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/client/metadata"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/corehandlers"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/request"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/signer/byteplussign"
)

var tlsUnmarshalErrorHandler = request.NamedHandler{Name: "TlsUnmarshalErrorHandler", Fn: tlsUnmarshalError}

func (u *BypassSvc) NewTlsClient() *client.Client {
	svc := "TLS"
	config := u.Session.ClientConfig(svc)
	var (
		endpoint string
	)

	c := client.New(
		*config.Config,
		metadata.ClientInfo{
			SigningName:   config.SigningName,
			SigningRegion: config.SigningRegion,
			Endpoint:      endpoint,
			ServiceName:   svc,
			ServiceID:     svc,
		},
		config.Handlers,
	)
	c.Handlers.Build.PushBackNamed(corehandlers.SDKVersionUserAgentHandler)
	c.Handlers.Build.PushBackNamed(corehandlers.AddHostExecEnvUserAgentHandler)
	c.Handlers.Sign.PushBackNamed(byteplussign.SignRequestHandler)
	c.Handlers.Build.PushBackNamed(bypassBuildHandler)
	c.Handlers.Unmarshal.PushBackNamed(bypassUnmarshalHandler)
	c.Handlers.UnmarshalError.PushBackNamed(tlsUnmarshalErrorHandler)

	return c
}

type tlsError struct {
	ErrorCode    string
	ErrorMessage string
	RequestId    string
}

func tlsUnmarshalError(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() {
		body, err := ioutil.ReadAll(r.HTTPResponse.Body)
		if err != nil {
			fmt.Printf("read byteplusbody err, %v\n", err)
			r.Error = err
			return
		}
		tos := tlsError{}
		if err = json.Unmarshal(body, &tos); err != nil {
			fmt.Printf("Unmarshal err, %v\n", err)
			r.Error = err
			return
		}
		r.Error = bytepluserr.NewRequestFailure(
			bytepluserr.New(tos.ErrorCode, tos.ErrorMessage, nil),
			r.HTTPResponse.StatusCode,
			r.HTTPResponse.Header.Get("X-Tls-Requestid"),
		)

		return
	} else {
		r.Error = bytepluserr.NewRequestFailure(
			bytepluserr.New("ServiceUnavailableException", "service is unavailable", nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}
}
