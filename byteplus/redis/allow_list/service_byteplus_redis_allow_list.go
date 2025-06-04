package allow_list

import (
	"errors"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ActionDeleteAllowList = "DeleteAllowList"
	ActionCreateAllowList = "CreateAllowList"
	ActionModifyAllowList = "ModifyAllowList"
)

type ByteplusRedisAllowListService struct {
	Client *bp.SdkClient
}

func NewRedisAllowListService(c *bp.SdkClient) *ByteplusRedisAllowListService {
	return &ByteplusRedisAllowListService{
		Client: c,
	}
}

func (s *ByteplusRedisAllowListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisAllowListService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(m, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeAllowLists"
		logger.Debug(logger.ReqFormat, action, condition)
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
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.AllowLists", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.AllowLists is not Slice")
		}

		for index, element := range data {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeAllowListDetail"), &map[string]interface{}{
				"AllowListId": element.(map[string]interface{})["AllowListId"],
			})
			if err != nil {
				return data, err
			}
			respResult, err := bp.ObtainSdkValue("Result", *resp)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, "DescribeAllowListDetail", *resp)
			// 多个地址间用英文逗号（,）隔开
			ips := respResult.(map[string]interface{})["AllowList"]
			data[index].(map[string]interface{})["AllowList"] = strings.Split(ips.(string), ",")
			data[index].(map[string]interface{})["AssociatedInstances"] = respResult.(map[string]interface{})["AssociatedInstances"]
		}
		return data, err
	})
}

func (s *ByteplusRedisAllowListService) ReadResource(resourceData *schema.ResourceData, allowlistId string) (data map[string]interface{}, err error) {
	if allowlistId == "" {
		allowlistId = s.ReadResourceId(resourceData.Id())
	}
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeAllowListDetail"), &map[string]interface{}{
		"AllowListId": allowlistId,
	})
	if err != nil {
		return data, err
	}
	respResult, err := bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	// 组合数据，DescribeAllowLists 无法查询出来
	data = respResult.(map[string]interface{})
	ips := respResult.(map[string]interface{})["AllowList"]
	data["AllowList"] = strings.Split(ips.(string), ",")
	data["AllowListIPNum"] = len(strings.Split(ips.(string), ","))
	data["AssociatedInstanceNum"] = len(data["AssociatedInstances"].([]interface{}))
	return data, err
}

func (s *ByteplusRedisAllowListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusRedisAllowListService) WithResourceResponseHandlers(allowlist map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return allowlist, map[string]bp.ResponseConvert{
			"AllowListIPNum": {
				TargetField: "allow_list_ip_num",
			},
			"VPC": {
				TargetField: "vpc",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRedisAllowListService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionCreateAllowList,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"allow_list": {
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						var res []string
						for _, ele := range i.(*schema.Set).List() {
							res = append(res, ele.(string))
						}
						return strings.Join(res, ",")
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, err := bp.ObtainSdkValue("Result.AllowListId", *resp)
				if err != nil {
					return err
				}
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisAllowListService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionModifyAllowList,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"allow_list_name": {
					ForceGet: true, // 必须传递
				},
				"allow_list_desc": {
					TargetField: "AllowListDesc",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllowListId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				if d.HasChange("allow_list") {
					data, err := s.ReadResource(d, d.Id())
					if err != nil {
						return nil, err
					}
					(*call.SdkParam)["ApplyInstanceNum"] = data["AssociatedInstanceNum"]

					var res []string
					for _, ele := range d.Get("allow_list").(*schema.Set).List() {
						res = append(res, ele.(string))
					}
					(*call.SdkParam)["AllowList"] = strings.Join(res, ",")
				}

				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisAllowListService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionDeleteAllowList,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllowListId"] = resourceData.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisAllowListService) DatasourceResources(data *schema.ResourceData, resource2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AllowListName",
		IdField:      "AllowListId",
		CollectField: "allow_lists",
		ResponseConverts: map[string]bp.ResponseConvert{
			"VPC": {
				TargetField: "vpc",
			},
			"AllowListIPNum": {
				TargetField: "allow_list_ip_num",
			},
		},
	}
}

func (s *ByteplusRedisAllowListService) ReadResourceId(id string) string {
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
