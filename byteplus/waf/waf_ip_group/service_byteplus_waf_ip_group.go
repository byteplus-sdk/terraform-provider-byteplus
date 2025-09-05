package waf_ip_group

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusWafIpGroupService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafIpGroupService(c *bp.SdkClient) *ByteplusWafIpGroupService {
	return &ByteplusWafIpGroupService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafIpGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafIpGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "Page", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListAllIpGroups"

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
		results, err = bp.ObtainSdkValue("Result.IpGroupList", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.IpGroupList is not Slice")
		}

		for _, ele := range data {
			ipGroup, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" ipGroup is not Map ")
			}

			ipGroupId := int(ipGroup["IpGroupId"].(float64))

			ipGroup["IpGroupIdString"] = strconv.Itoa(ipGroupId)

			logger.Debug(logger.ReqFormat, "IpGroupIdString", ipGroup["IpGroupIdString"])

			// 查询域名详细信息
			action := "ListIpGroup"
			req := map[string]interface{}{
				"IpGroupId": ipGroupId,
			}
			logger.Debug(logger.ReqFormat, action, req)

			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
			listIpGroup, err := bp.ObtainSdkValue("Result", *resp)
			if err != nil {
				return data, err
			}
			listIpGroupMap, ok := listIpGroup.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" Result is not Map ")
			}

			ipGroup["IpList"] = listIpGroupMap["IpList"]
		}

		return data, err
	})
}

func (s *ByteplusWafIpGroupService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
		result  map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ipGroupId, err := strconv.Atoi(id)
	if err != nil {
		return data, fmt.Errorf(" ipGroupId cannot convert to int ")
	}

	req := map[string]interface{}{
		"TimeOrderBy": "DESC",
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}

		if int(data["IpGroupId"].(float64)) == ipGroupId {
			result = data
			break
		}
	}
	if len(result) == 0 {
		return result, fmt.Errorf("waf_host_group %s not exist ", id)
	}
	return result, err
}

func (s *ByteplusWafIpGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafIpGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddIpGroup",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"ip_list": {
					TargetField: "IpList",
					ConvertType: bp.ConvertJsonArray,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.IpGroupId", *resp)
				d.SetId(strconv.Itoa(int(id.(float64))))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafIpGroupService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafIpGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateIpGroup",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"name": {
					TargetField: "Name",
					ForceGet:    true,
				},
				"add_type": {
					TargetField: "AddType",
					ForceGet:    true,
				},
				"ip_list": {
					TargetField: "IpList",
					ConvertType: bp.ConvertJsonArray,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipGroupId, err := strconv.Atoi(d.Id())
				if err != nil {
					return false, fmt.Errorf(" ipGroupId cannot convert to int ")
				}
				(*call.SdkParam)["IpGroupId"] = ipGroupId
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

func (s *ByteplusWafIpGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteIpGroup",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipGroupId, err := strconv.Atoi(d.Id())
				if err != nil {
					return false, fmt.Errorf(" ipGroupId cannot convert to int ")
				}
				(*call.SdkParam)["IpGroupIds"] = []int{ipGroupId}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading waf ip group on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusWafIpGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "Name",
		IdField:      "IpGroupIdString",
		CollectField: "ip_group_list",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"related_rules": {
				TargetField: "RelatedRules",
			},
		},
	}
}

func (s *ByteplusWafIpGroupService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "waf",
		Version:     "2023-12-25",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
		RegionType:  bp.Global,
	}
}

//
//func (s *ByteplusWafIpGroupService) checkResourceUtilRemoved(d *schema.ResourceData, timeout time.Duration) error {
//	return resource.Retry(timeout, func() *resource.RetryError {
//		ipGroup, _ := s.ReadResource(d, d.Id())
//		logger.Debug(logger.RespFormat, "ipGroup", ipGroup)
//
//		// 能查询成功代表还在删除中，重试
//		ipList, ok := ipGroup["IpList"].([]string)
//		if !ok {
//			return resource.NonRetryableError(fmt.Errorf("ipList is not []string"))
//		}
//		if len(ipList) != 0 {
//			return resource.RetryableError(fmt.Errorf("resource still in removing status "))
//		} else {
//			if len(ipList) == 0 {
//				return nil
//			} else {
//				return resource.NonRetryableError(fmt.Errorf("ipGroup status is not deleted "))
//			}
//		}
//	})
//}
