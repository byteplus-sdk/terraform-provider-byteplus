package cdn_kv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnKvService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnKvService(c *bp.SdkClient) *ByteplusCdnKvService {
	return &ByteplusCdnKvService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnKvService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnKvService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageOffsetQuery(m, "Limit", "Page", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListKvKey"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))

		results, err = bp.ObtainSdkValue("Result.NamespaceKeys", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.NamespaceKeys is not Slice")
		}
		return data, err
	})

	for _, v := range data {
		namespaceKey, ok := v.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf(" The Sparrow of Result is not map ")
		}

		// 查询 value
		valueAction := "GetKeyValue"
		valueReq := map[string]interface{}{
			"NamespaceId": namespaceKey["NamespaceId"],
			"Namespace":   namespaceKey["Namespace"],
			"Key":         namespaceKey["Key"],
		}
		logger.Debug(logger.ReqFormat, valueAction, valueReq)
		valueResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(valueAction), &valueReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, valueAction, valueResp)
		value, err := bp.ObtainSdkValue("Result.Value", *valueResp)
		if err != nil {
			return data, err
		}
		namespaceKey["Value"] = value
	}

	return data, err
}

func (s *ByteplusCdnKvService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 3 {
		return data, fmt.Errorf(" Invalid CdnKvNamespaceKey Id %s ", id)
	}

	req := map[string]interface{}{
		"NamespaceId": ids[0],
		"Namespace":   ids[1],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		var namespaceKey map[string]interface{}
		if namespaceKey, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if namespaceKey["Key"].(string) == ids[2] {
			data = namespaceKey
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("cdn_kv %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnKvService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (ByteplusCdnKvService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"DDL": {
				TargetField: "ddl",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnKvService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateKeyValue",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"ttl": {
					TargetField: "TTL",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				namespaceId := d.Get("namespace_id").(string)
				namespace := d.Get("namespace").(string)
				key := d.Get("key").(string)
				d.SetId(fmt.Sprintf(namespaceId + ":" + namespace + ":" + key))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnKvService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateKeyValue",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"value": {
					TargetField: "Value",
				},
				"ttl": {
					TargetField: "TTL",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					ids := strings.Split(d.Id(), ":")
					if len(ids) != 3 {
						return false, fmt.Errorf(" Invalid CdnKvNamespaceKey Id %s ", d.Id())
					}

					(*call.SdkParam)["NamespaceId"] = ids[0]
					(*call.SdkParam)["Namespace"] = ids[1]
					(*call.SdkParam)["Key"] = ids[2]
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnKvService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteKvKey",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 3 {
					return false, fmt.Errorf(" Invalid CdnKvNamespaceKey Id %s ", d.Id())
				}

				(*call.SdkParam)["NamespaceId"] = ids[0]
				(*call.SdkParam)["Namespace"] = ids[1]
				(*call.SdkParam)["Key"] = ids[2]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading kv namespace key on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnKvService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "Key",
		CollectField: "namespace_keys",
		ResponseConverts: map[string]bp.ResponseConvert{
			"DDL": {
				TargetField: "ddl",
			},
		},
	}
}

func (s *ByteplusCdnKvService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func getPostUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
