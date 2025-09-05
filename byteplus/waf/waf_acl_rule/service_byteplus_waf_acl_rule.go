package waf_acl_rule

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusWafAclRuleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafAclRuleService(c *bp.SdkClient) *ByteplusWafAclRuleService {
	return &ByteplusWafAclRuleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafAclRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafAclRuleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "Page", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListAclRule"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.Rules", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Rules is not Slice")
		}
		for _, ele := range data {
			aclRule, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" aclRule is not Map ")
			}

			aclRule["AclRuleID"] = strconv.Itoa(int(aclRule["ID"].(float64)))

			logger.Debug(logger.ReqFormat, "AclRuleID", aclRule["AclRuleID"])

		}
		return data, err
	})
}

func (s *ByteplusWafAclRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
		ruleTag string
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return data, fmt.Errorf("format of waf acl rule resource id is invalid,%s", id)
	}
	aclRuleId := parts[0]
	aclType := parts[1]
	ruleId, err := strconv.Atoi(aclRuleId)
	tag := fmt.Sprintf("%012d", ruleId)
	if aclType == "" {
		return data, errors.New("acl_type is null")
	}

	if aclType == "Allow" {
		ruleTag = "A" + tag
	}

	if aclType == "Block" {
		ruleTag = "B" + tag
	}
	req := map[string]interface{}{
		"RuleTag": ruleTag,
		"AclType": aclType,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("waf_acl_rule %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusWafAclRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafAclRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAclRule",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"host_group_id": {
					TargetField: "HostGroupId",
					ConvertType: bp.ConvertJsonArray,
				},
				"ip_location_country": {
					TargetField: "IpLocationCountry",
					ConvertType: bp.ConvertJsonArray,
				},
				"ip_location_subregion": {
					TargetField: "IpLocationSubregion",
					ConvertType: bp.ConvertJsonArray,
				},
				"accurate_group": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "AccurateGroup",
					NextLevelConvert: map[string]bp.RequestConvert{
						"accurate_rules": {
							ConvertType: bp.ConvertJsonObjectArray,
							TargetField: "AccurateRules",
							NextLevelConvert: map[string]bp.RequestConvert{
								"http_obj": {
									TargetField: "HttpObj",
								},
								"obj_type": {
									TargetField: "ObjType",
								},
								"opretar": {
									TargetField: "Opretar",
								},
								"property": {
									TargetField: "Property",
								},
								"value_string": {
									TargetField: "ValueString",
								},
							},
						},
						"logic": {
							TargetField: "Logic",
						},
					},
				},
				"host_list": {
					TargetField: "HostList",
					ConvertType: bp.ConvertJsonArray,
				},
				"ip_group_id": {
					TargetField: "IpGroupId",
					ConvertType: bp.ConvertJsonArray,
				},
				"ip_list": {
					TargetField: "IpList",
					ConvertType: bp.ConvertJsonArray,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				enable, ok := d.Get("enable").(int)
				if !ok {
					return false, errors.New("enable is not int")
				}
				if enable == 0 {
					(*call.SdkParam)["Enable"] = 0
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				aclTypeString, ok := d.Get("acl_type").(string)
				if !ok {
					return errors.New("acl_type is not string")
				}
				d.SetId(fmt.Sprintf("%s:%s", strconv.Itoa(int(id.(float64))), aclTypeString))

				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafAclRuleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"ID": {
				TargetField: "id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafAclRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	// 服务返回问题 暂时不适配
	//callback := bp.Callback{
	//	Call: bp.SdkCall{
	//		Action:      "UpdateAclRule",
	//		ConvertMode: bp.RequestConvertAll,
	//		ContentType: bp.ContentTypeJson,
	//		Convert: map[string]bp.RequestConvert{
	//			"action": {
	//				TargetField: "Action",
	//				ForceGet:    true,
	//			},
	//			"prefix_switch": {
	//				TargetField: "PrefixSwitch",
	//				ForceGet:    true,
	//			},
	//			"description": {
	//				TargetField: "Description",
	//				ForceGet:    true,
	//			},
	//			"advanced": {
	//				TargetField: "Advanced",
	//				ForceGet:    true,
	//			},
	//			"ip_add_type": {
	//				TargetField: "IpAddType",
	//				ForceGet:    true,
	//			},
	//			"host_add_type": {
	//				TargetField: "HostAddType",
	//				ForceGet:    true,
	//			},
	//			"url": {
	//				TargetField: "Url",
	//				ForceGet:    true,
	//			},
	//			"enable": {
	//				TargetField: "Enable",
	//				ForceGet:    true,
	//			},
	//			"name": {
	//				TargetField: "Name",
	//				ForceGet:    true,
	//			},
	//			"acl_type": {
	//				TargetField: "AclType",
	//				ForceGet:    true,
	//			},
	//			"ip_location_country": {
	//				TargetField: "IpLocationCountry",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"ip_list": {
	//				TargetField: "IpList",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"ip_group_id": {
	//				TargetField: "IpGroupId",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"host_list": {
	//				TargetField: "HostList",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"host_group_id": {
	//				TargetField: "HostGroupId",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"ip_location_subregion": {
	//				TargetField: "IpLocationSubregion",
	//				ConvertType: bp.ConvertJsonArray,
	//				ForceGet:    true,
	//			},
	//			"accurate_group": {
	//				ConvertType: bp.ConvertJsonObject,
	//				TargetField: "AccurateGroup",
	//				ForceGet:    true,
	//				NextLevelConvert: map[string]bp.RequestConvert{
	//					"accurate_rules": {
	//						ConvertType: bp.ConvertJsonObjectArray,
	//						TargetField: "AccurateRules",
	//						NextLevelConvert: map[string]bp.RequestConvert{
	//							"http_obj": {
	//								TargetField: "HttpObj",
	//							},
	//							"obj_type": {
	//								TargetField: "ObjType",
	//							},
	//							"opretar": {
	//								TargetField: "Opretar",
	//							},
	//							"property": {
	//								TargetField: "Property",
	//							},
	//							"value_string": {
	//								TargetField: "ValueString",
	//							},
	//						},
	//					},
	//					"logic": {
	//						TargetField: "Logic",
	//					},
	//				},
	//			},
	//		},
	//		BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
	//			aclTypeId, err := strconv.Atoi(d.Id())
	//			if err != nil {
	//				return false, fmt.Errorf(" aclTypeId cannot convert to int ")
	//			}
	//			(*call.SdkParam)["ID"] = aclTypeId
	//
	//			logic, ok := d.Get("accurate_group.0.logic").(int)
	//			if !ok {
	//				return false, fmt.Errorf("accurate_group.0.logic cannot convert to int ")
	//			}
	//
	//			if logic == 0 {
	//				delete(*call.SdkParam, "AccurateGroup.Logic")
	//			}
	//
	//			return true, nil
	//		},
	//		ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
	//			logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
	//			resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
	//			logger.Debug(logger.RespFormat, call.Action, resp, err)
	//			return resp, err
	//		},
	//	},
	//}
	return []bp.Callback{}
}

func (s *ByteplusWafAclRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAclRule",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {

				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf acl rule resource id is invalid,%s", d.Id())
				}
				aclRuleId := parts[0]
				aclType := parts[1]

				aclTypeId, err := strconv.Atoi(aclRuleId)
				if err != nil {
					return false, fmt.Errorf(" aclTypeId cannot convert to int ")
				}
				(*call.SdkParam)["ID"] = aclTypeId
				(*call.SdkParam)["AclType"] = aclType
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading waf acl rule on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusWafAclRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"action": {
				TargetField: "Action",
				ConvertType: bp.ConvertJsonArray,
			},
			"defence_host": {
				TargetField: "DefenceHost",
				ConvertType: bp.ConvertJsonArray,
			},
			"enable": {
				TargetField: "Enable",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		NameField:    "Name",
		IdField:      "AclRuleID",
		CollectField: "rules",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"ID": {
				TargetField: "id",
			},
		},
	}
}

func (s *ByteplusWafAclRuleService) ReadResourceId(id string) string {
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
