package universal

import "github.com/byteplus-sdk/byteplus-go-sdk/byteplus/session"

type Universal struct {
	Session *session.Session
}

type RequestUniversal struct {
	ServiceName string
	Action      string
	Version     string
	HttpMethod  HttpMethod
	ContentType ContentType
}
