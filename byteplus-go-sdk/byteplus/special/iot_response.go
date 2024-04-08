package special

import (
	"reflect"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/response"
)

func iotResponse(response response.ByteplusResponse, i interface{}) interface{} {
	_, ok1 := reflect.TypeOf(i).Elem().FieldByName("ResponseMetadata")
	_, ok2 := reflect.TypeOf(i).Elem().FieldByName("Result")
	if ok1 && ok2 {
		return response
	}
	return response.Result
}
