package allowlist

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlAllowListService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func (s *ByteplusRdsMysqlAllowListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlAllowListService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp        *map[string]interface{}
		results     interface{}
		ok          bool
		allowListId string
	)
	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
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
		results, err = bp.ObtainSdkValue("Result.AllowLists", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.AllowLists is not slice ")
		}

		if id, exist := condition["AllowListId"]; exist {
			allowListId = id.(string)
		}
		for index, ele := range data {
			allowList := ele.(map[string]interface{})

			if allowListId == "" || allowListId == allowList["AllowListId"].(string) {
				query := map[string]interface{}{
					"AllowListId": allowList["AllowListId"],
				}
				action = "DescribeAllowListDetail"
				logger.Debug(logger.ReqFormat, action, query)
				resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
				if err != nil {
					return data, err
				}
				logger.Debug(logger.RespFormat, action, query, *resp)
				instances, err := bp.ObtainSdkValue("Result.AssociatedInstances", *resp)
				if err != nil {
					return data, err
				}
				data[index].(map[string]interface{})["AssociatedInstances"] = instances
				allowListIp, err := bp.ObtainSdkValue("Result.AllowList", *resp)
				if err != nil {
					return data, err
				}
				allowListIpArr := strings.Split(allowListIp.(string), ",")
				data[index].(map[string]interface{})["AllowList"] = allowListIpArr
				userAllowList, err := bp.ObtainSdkValue("Result.UserAllowList", *resp)
				if err != nil {
					return data, err
				}
				userAllowListArr := strings.Split(userAllowList.(string), ",")
				data[index].(map[string]interface{})["UserAllowList"] = userAllowListArr
			}
		}
		return data, err
	})
}

func (s *ByteplusRdsMysqlAllowListService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"RegionId":    s.Client.Region,
		"AllowListId": id,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		result, ok := v.(map[string]interface{})
		if !ok {
			return data, errors.New("Value is not map ")
		}
		if result["AllowListId"].(string) == id {
			data = result
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Rds instance %s not exist ", id)
	}
	if sgIds, ok := resourceData.GetOk("security_group_ids"); ok {
		data["SecurityGroupIds"] = sgIds.(*schema.Set).List()
	}
	return data, err
}

func (s *ByteplusRdsMysqlAllowListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusRdsMysqlAllowListService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, map[string]bp.ResponseConvert{}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlAllowListService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAllowList",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_list": {
					Ignore: true,
				},
				"security_group_ids": {
					ConvertType: bp.ConvertListN,
				},
				"security_group_bind_infos": {
					Ignore: true,
				},
				"user_allow_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				allowListsInterface, ok := d.GetOk("allow_list")
				if ok {
					allowListsSet := allowListsInterface.(*schema.Set)
					var allowStrings []string
					allowLists := allowListsSet.List()
					for _, list := range allowLists {
						allowStrings = append(allowStrings, list.(string))
					}
					lists := strings.Join(allowStrings, ",")
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam, lists)
					(*call.SdkParam)["AllowList"] = lists
				}

				userAllowListsInterface, ok := d.GetOk("user_allow_list")
				if ok {
					userAllowListsSet := userAllowListsInterface.(*schema.Set)
					var userAllowStrings []string
					userAllowLists := userAllowListsSet.List()
					for _, list := range userAllowLists {
						userAllowStrings = append(userAllowStrings, list.(string))
					}
					userLists := strings.Join(userAllowStrings, ",")
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam, userLists)
					(*call.SdkParam)["UserAllowList"] = userLists
				}

				securityGroupBindInfoInterface, ok := d.GetOk("security_group_bind_infos")
				if ok {
					securityGroupBindInfoSet := securityGroupBindInfoInterface.(*schema.Set)
					var securityGroupBindInfos []map[string]interface{}
					securityGroupBindInfoList := securityGroupBindInfoSet.List()
					for _, list := range securityGroupBindInfoList {
						securityGroupBindInfo := list.(map[string]interface{})
						result := make(map[string]interface{})
						result["BindMode"] = securityGroupBindInfo["bind_mode"]
						result["SecurityGroupId"] = securityGroupBindInfo["security_group_id"]
						securityGroupBindInfos = append(securityGroupBindInfos, result)
					}
					(*call.SdkParam)["SecurityGroupBindInfos"] = securityGroupBindInfos
					logger.Debug(logger.ReqFormat, call.Action+" SecurityGroupBindInfos", call.SdkParam, securityGroupBindInfoList)
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.AllowListId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlAllowListService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyAllowList",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_list": {
					Ignore: true,
				},
				"apply_instance_num": {
					Ignore: true,
				},
				"allow_list_desc": {
					ForceGet: true,
				},
				"security_group_ids": {
					Ignore: true,
				},
				"security_group_bind_infos": {
					Ignore: true,
				},
				"allow_list_category": {
					TargetField: "AllowListCategory",
				},
				"user_allow_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// 修改allowList必须传ApplyInstanceNum
				resp, err := s.ReadResource(d, d.Id())
				if err != nil {
					return false, err
				}
				num := resp["AssociatedInstanceNum"].(float64)
				(*call.SdkParam)["ApplyInstanceNum"] = int(num)
				if d.HasChange("allow_list") {
					var allowStrings []string
					allowListsSet := d.Get("allow_list").(*schema.Set)
					allowLists := allowListsSet.List()
					for _, list := range allowLists {
						allowStrings = append(allowStrings, list.(string))
					}
					lists := strings.Join(allowStrings, ",")
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam, lists)
					(*call.SdkParam)["AllowList"] = lists
				}
				// 这里逻辑需要改一下
				// 接口逻辑是即使这里没改，也得把它和安全组信息传过去，否则会覆盖删除
				if _, ok := d.GetOk("user_allow_list"); ok {
					var userAllowStrings []string
					userAllowListsSet := d.Get("user_allow_list").(*schema.Set)
					userAllowLists := userAllowListsSet.List()
					for _, list := range userAllowLists {
						userAllowStrings = append(userAllowStrings, list.(string))
					}
					userLists := strings.Join(userAllowStrings, ",")
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam, userLists)
					(*call.SdkParam)["UserAllowList"] = userLists
				}

				if securityGroupIdsInterface, ok := d.GetOk("security_group_ids"); ok {
					securityGroupIdsSet := securityGroupIdsInterface.(*schema.Set)
					securityGroupIds := securityGroupIdsSet.List()
					(*call.SdkParam)["SecurityGroupIds"] = securityGroupIds
				}

				if _, ok := d.GetOk("security_group_bind_infos"); ok {
					securityGroupBindInfoInterface := d.Get("security_group_bind_infos")
					securityGroupBindInfoSet := securityGroupBindInfoInterface.(*schema.Set)
					var securityGroupBindInfos []map[string]interface{}
					securityGroupBindInfoList := securityGroupBindInfoSet.List()
					for _, list := range securityGroupBindInfoList {
						securityGroupBindInfo := list.(map[string]interface{})
						result := make(map[string]interface{})
						result["BindMode"] = securityGroupBindInfo["bind_mode"]
						result["SecurityGroupId"] = securityGroupBindInfo["security_group_id"]
						securityGroupBindInfos = append(securityGroupBindInfos, result)
					}
					(*call.SdkParam)["SecurityGroupBindInfos"] = securityGroupBindInfos
				}

				(*call.SdkParam)["UpdateSecurityGroup"] = true
				return true, nil
			},
			SdkParam: &map[string]interface{}{
				"AllowListId":   data.Id(),
				"AllowListName": data.Get("allow_list_name").(string),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlAllowListService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAllowList",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"AllowListId": data.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlAllowListService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AllowListName",
		IdField:      "AllowListId",
		CollectField: "allow_lists",
		ResponseConverts: map[string]bp.ResponseConvert{
			"AllowListIPNum": {
				TargetField: "allow_list_ip_num",
			},
			"VPC": {
				TargetField: "vpc",
			},
		},
	}
}

func (s *ByteplusRdsMysqlAllowListService) ReadResourceId(id string) string {
	return id
}

func NewRdsMysqlAllowListService(client *bp.SdkClient) *ByteplusRdsMysqlAllowListService {
	return &ByteplusRdsMysqlAllowListService{
		Client:     client,
		Dispatcher: &bp.Dispatcher{},
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_mysql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
