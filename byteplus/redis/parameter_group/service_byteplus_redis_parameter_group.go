package parameter_group

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

type ByteplusParameterGroupService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

type ParamValue struct {
	Name  string
	Value string
}

func NewParameterGroupService(c *bp.SdkClient) *ByteplusParameterGroupService {
	return &ByteplusParameterGroupService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusParameterGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusParameterGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeParameterGroups"

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
		results, err = bp.ObtainSdkValue("Result.ParameterGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ParameterGroups is not Slice")
		}

		for _, ele := range data {
			parameterGroups, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" ParameterGroup is not Map ")
			}

			parameterGroupId, ok := parameterGroups["ParameterGroupId"]
			if !ok {
				return data, fmt.Errorf(" ParameterGroupId is not String ")
			}
			// 查询参数模版详细信息
			action := "DescribeParameterGroupDetail"
			req := map[string]interface{}{
				"ParameterGroupId": parameterGroupId,
			}
			logger.Debug(logger.ReqFormat, action, req)

			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
			parameterGroup, err := bp.ObtainSdkValue("Result.ParameterGroupInfo", *resp)
			if err != nil {
				return data, err
			}
			parameterGroupMap, ok := parameterGroup.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" Result is not Map ")
			}
			for k, v := range parameterGroupMap {
				parameterGroups[k] = v
			}
		}
		return data, err
	})
}

func (s *ByteplusParameterGroupService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results            []interface{}
		ok                 bool
		result             map[string]interface{}
		paramValuesStructs []ParamValue
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}

		parameterGroupId, ok := data["ParameterGroupId"].(string)
		if !ok {
			return data, fmt.Errorf(" ParameterGroupId is not String ")
		}

		if parameterGroupId == id {
			result = data
			break
		}

	}
	if len(result) == 0 {
		return data, fmt.Errorf("parameter_group %s not exist ", id)
	}

	parameters, ok := result["Parameters"].([]interface{})
	if !ok {
		return data, fmt.Errorf(" Parameters is not Slice ")
	}

	for _, parameter := range parameters {
		currentValue, ok := parameter.(map[string]interface{})["CurrentValue"].(string)
		if !ok {
			return data, fmt.Errorf(" CurrentValue is not String ")
		}
		paramName, ok := parameter.(map[string]interface{})["ParamName"].(string)
		if !ok {
			return data, fmt.Errorf(" ParamName is not String ")
		}
		paramValues, ok := resourceData.Get("param_values").([]interface{})
		if !ok {
			return data, fmt.Errorf(" param_values is not Slice ")
		}

		for _, paramValue := range paramValues {
			paramValueMap, ok := paramValue.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" param_value is not Map ")
			}
			paramNameValue, ok := paramValueMap["name"].(string)
			if !ok {
				return data, fmt.Errorf(" name is not String ")
			}
			if paramNameValue == paramName {
				paramValuesStruct := ParamValue{
					Name:  paramName,
					Value: currentValue,
				}
				paramValuesStructs = append(paramValuesStructs, paramValuesStruct)
			}
		}
	}

	result["Parameters.Options"] = paramValuesStructs

	logger.Debug(logger.RespFormat, "parameters group data is ", data)

	return data, err
}

func (s *ByteplusParameterGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusParameterGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateParameterGroup",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"param_values": {
					TargetField: "ParamValues",
					ConvertType: bp.ConvertJsonObjectArray,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ParameterGroupId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusParameterGroupService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusParameterGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyParameterGroup",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"param_values": {
					TargetField: "ParamValues",
					ConvertType: bp.ConvertJsonObjectArray,
					ForceGet:    true,
				},
				"name": {
					TargetField: "Name",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ParameterGroupId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusParameterGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteParameterGroup",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"ParameterGroupId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				// 开启删除保护时，跳过 CallError
				if strings.Contains(baseErr.Error(), "can not delete redis parameter group") {
					return baseErr
				}
				// 出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading redis parameter group on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusParameterGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts:  map[string]bp.RequestConvert{},
		NameField:        "Name",
		IdField:          "ParameterGroupId",
		CollectField:     "parameter_groups",
		ResponseConverts: map[string]bp.ResponseConvert{},
	}
}

func (s *ByteplusParameterGroupService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
