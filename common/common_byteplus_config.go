package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/byteplusutil"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/credentials"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/session"
)

type Config struct {
	AccessKey              string
	SecretKey              string
	SessionToken           string
	Region                 string
	Endpoint               string
	DisableSSL             bool
	EnableStandardEndpoint bool
	StandardEndpointSuffix string
	CustomerHeaders        map[string]string
	CustomerEndpoints      map[string]string
	ProxyUrl               string
}

func (c *Config) Client() (*SdkClient, error) {
	var client SdkClient
	version := fmt.Sprintf("%s/%s", TerraformProviderName, TerraformProviderVersion)

	config := byteplus.NewConfig().
		WithRegion(c.Region).
		WithExtraUserAgent(byteplus.String(version)).
		WithCredentials(credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)).
		WithDisableSSL(c.DisableSSL).
		WithExtendHttpRequest(func(ctx context.Context, request *http.Request) {
			if len(c.CustomerHeaders) > 0 {
				for k, v := range c.CustomerHeaders {
					request.Header.Add(k, v)
				}
			}
		}).
		WithEndpoint(byteplusutil.NewEndpoint().WithCustomerEndpoint(c.Endpoint).GetEndpoint())

	if c.ProxyUrl != "" {
		u, _ := url.Parse(c.ProxyUrl)
		t := &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		httpClient := http.DefaultClient
		httpClient.Transport = t
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("session init error %w", err)
	}

	client.Region = c.Region
	client.UniversalClient = NewUniversalClient(sess, c.CustomerEndpoints, c.EnableStandardEndpoint, c.StandardEndpointSuffix)
	client.BypassSvcClient = NewBypassClient(sess)

	return &client, nil
}

func init() {
	InitLocks()
}
