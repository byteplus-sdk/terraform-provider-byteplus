package special

import "github.com/byteplus-sdk/byteplus-go-sdk/byteplus/response"

type ResponseSpecial func(response.ByteplusResponse, interface{}) interface{}

var responseSpecialMapping map[string]ResponseSpecial

func init() {
	responseSpecialMapping = map[string]ResponseSpecial{
		"iot": iotResponse,
	}
}

func ResponseSpecialMapping() map[string]ResponseSpecial {
	return responseSpecialMapping
}
