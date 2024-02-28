package custom

import (
	"context"
	"net/http"
	"net/url"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/credentials"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/response"
)

type RequestMetadata struct {
	ServiceName string
	Version     string
	Action      string
	HttpMethod  string
	Region      string
	Request     *http.Request
	RawQuery    *url.Values
}

type ExtendContextWithMeta func(ctx context.Context, meta RequestMetadata) context.Context

type ExtendHttpRequest func(ctx context.Context, request *http.Request)

type ExtendHttpRequestWithMeta func(ctx context.Context, request *http.Request, meta RequestMetadata)

type ExtraHttpParameters func(ctx context.Context) map[string]string

type ExtraHttpParametersWithMeta func(ctx context.Context, meta RequestMetadata) map[string]string

type ExtraHttpJsonBody func(ctx context.Context, input *map[string]interface{}, meta RequestMetadata)

type LogAccount func(ctx context.Context) *string

type DynamicCredentials func(ctx context.Context) (*credentials.Credentials, *string)

// DynamicCredentialsIncludeError func return Credentials info and error info when error appear
type DynamicCredentialsIncludeError func(ctx context.Context) (*credentials.Credentials, *string, error)

type CustomerUnmarshalError func(ctx context.Context, meta RequestMetadata, resp response.ByteplusResponse) error

type CustomerUnmarshalData func(ctx context.Context, info RequestInfo, resp response.ByteplusResponse) interface{}

type ForceJsonNumberDecode func(ctx context.Context, info RequestInfo) bool
