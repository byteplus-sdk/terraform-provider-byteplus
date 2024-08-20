package classic_cdn_shared_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnSharedConfigService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnSharedConfigService(c *bp.SdkClient) *ByteplusCdnSharedConfigService {
	return &ByteplusCdnSharedConfigService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnSharedConfigService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnSharedConfigService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNum", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListSharedConfig"

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
		results, err = bp.ObtainSdkValue("Result.ConfigData", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ConfigData is not Slice")
		}
		for index, d := range data {
			config := d.(map[string]interface{})
			query := map[string]interface{}{
				"ConfigName": config["ConfigName"],
			}
			action = "DescribeSharedConfig"
			logger.Debug(logger.ReqFormat, action, query)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, query, *resp)
			allowIp, err := bp.ObtainSdkValue("Result.AllowIpAccessRule", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["AllowIpAccessRule"] = allowIp
			denyIp, err := bp.ObtainSdkValue("Result.DenyIpAccessRule", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["DenyIpAccessRule"] = denyIp
			allowReferer, err := bp.ObtainSdkValue("Result.AllowRefererAccessRule", *resp)
			if err != nil {
				return data, err
			}
			if allowReferer != nil {
				allowReferer.(map[string]interface{})["CommonType"] = []interface{}{allowReferer.(map[string]interface{})["CommonType"]}
			}
			data[index].(map[string]interface{})["AllowRefererAccessRule"] = allowReferer
			denyReferer, err := bp.ObtainSdkValue("Result.DenyRefererAccessRule", *resp)
			if err != nil {
				return data, err
			}
			if denyReferer != nil {
				denyReferer.(map[string]interface{})["CommonType"] = []interface{}{denyReferer.(map[string]interface{})["CommonType"]}
			}
			data[index].(map[string]interface{})["DenyRefererAccessRule"] = denyReferer
			common, err := bp.ObtainSdkValue("Result.CommonMatchList", *resp)
			if err != nil {
				return data, err
			}
			if common != nil {
				common.(map[string]interface{})["CommonType"] = []interface{}{common.(map[string]interface{})["CommonType"]}
			}
			data[index].(map[string]interface{})["CommonMatchList"] = common
		}
		return data, err
	})
}

func (s *ByteplusCdnSharedConfigService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ConfigName": id,
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
		return data, fmt.Errorf("cdn_shared_config %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnSharedConfigService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
	}
}

func (s *ByteplusCdnSharedConfigService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddSharedConfig",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_ip_access_rule": {
					Ignore: true,
				},
				"deny_ip_access_rule": {
					Ignore: true,
				},
				"project_name": {
					TargetField: "project",
				},
				"allow_referer_access_rule": {
					Ignore: true,
				},
				"deny_referer_access_rule": {
					Ignore: true,
				},
				"common_match_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if allowIp, ok := d.GetOk("allow_ip_access_rule"); ok {
					result := make(map[string]interface{})
					if list, ok := allowIp.([]interface{}); ok && len(list) > 0 {
						ipMap := list[0].(map[string]interface{})
						rules := ipMap["rules"]
						result["Rules"] = rules.(*schema.Set).List()
					}
					(*call.SdkParam)["AllowIpAccessRule"] = result
				}

				if denyIp, ok := d.GetOk("deny_ip_access_rule"); ok {
					result := make(map[string]interface{})
					if list, ok := denyIp.([]interface{}); ok && len(list) > 0 {
						ipMap := list[0].(map[string]interface{})
						rules := ipMap["rules"]
						result["Rules"] = rules.(*schema.Set).List()
					}
					(*call.SdkParam)["DenyIpAccessRule"] = result
				}

				if allowRefererAccessRule, ok := d.GetOk("allow_referer_access_rule"); ok {
					result := make(map[string]interface{})
					if list, ok := allowRefererAccessRule.([]interface{}); ok && len(list) > 0 {
						allowRefererAccessRuleMap := list[0].(map[string]interface{})
						if allowEmpty, ok := allowRefererAccessRuleMap["allow_empty"]; ok {
							result["AllowEmpty"] = allowEmpty
						}
						commonTypeResult := make(map[string]interface{})
						commonType := allowRefererAccessRuleMap["common_type"]
						if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
							commonTypeMap := commonTypeList[0].(map[string]interface{})
							if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
								commonTypeResult["IgnoreCase"] = ignoreCase
							}
							rules := commonTypeMap["rules"]
							commonTypeResult["Rules"] = rules.(*schema.Set).List()
							result["CommonType"] = commonTypeResult
						}
					}
					(*call.SdkParam)["AllowRefererAccessRule"] = result
				}

				if denyRefererAccessRule, ok := d.GetOk("deny_referer_access_rule"); ok {
					result := make(map[string]interface{})
					if list, ok := denyRefererAccessRule.([]interface{}); ok && len(list) > 0 {
						denyRefererAccessRuleMap := list[0].(map[string]interface{})
						if allowEmpty, ok := denyRefererAccessRuleMap["allow_empty"]; ok {
							result["AllowEmpty"] = allowEmpty
						}
						commonTypeResult := make(map[string]interface{})
						commonType := denyRefererAccessRuleMap["common_type"]
						if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
							commonTypeMap := commonTypeList[0].(map[string]interface{})
							if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
								commonTypeResult["IgnoreCase"] = ignoreCase
							}
							rules := commonTypeMap["rules"]
							commonTypeResult["Rules"] = rules.(*schema.Set).List()
							result["CommonType"] = commonTypeResult
						}
					}
					(*call.SdkParam)["DenyRefererAccessRule"] = result
				}

				if common, ok := d.GetOk("common_match_list"); ok {
					result := make(map[string]interface{})
					if list, ok := common.([]interface{}); ok && len(list) > 0 {
						commonMap := list[0].(map[string]interface{})
						commonTypeResult := make(map[string]interface{})
						commonType := commonMap["common_type"]
						if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
							commonTypeMap := commonTypeList[0].(map[string]interface{})
							if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
								commonTypeResult["IgnoreCase"] = ignoreCase
							}
							rules := commonTypeMap["rules"]
							commonTypeResult["Rules"] = rules.(*schema.Set).List()
							result["CommonType"] = commonTypeResult
						}
					}
					(*call.SdkParam)["CommonMatchList"] = result
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
				id := d.Get("config_name")
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusCdnSharedConfigService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"Project": {
				TargetField: "project_name",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnSharedConfigService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateSharedConfig",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ConfigName"] = d.Id()

				if d.HasChange("allow_ip_access_rule") {
					if allowIp, ok := d.GetOk("allow_ip_access_rule"); ok {
						result := make(map[string]interface{})
						if list, ok := allowIp.([]interface{}); ok && len(list) > 0 {
							ipMap := list[0].(map[string]interface{})
							rules := ipMap["rules"]
							result["Rules"] = rules.(*schema.Set).List()
						}
						(*call.SdkParam)["AllowIpAccessRule"] = result
					} else {
						(*call.SdkParam)["AllowIpAccessRule"] = map[string]interface{}{}
					}
				}

				if d.HasChange("deny_ip_access_rule") {
					if denyIp, ok := d.GetOk("deny_ip_access_rule"); ok {
						result := make(map[string]interface{})
						if list, ok := denyIp.([]interface{}); ok && len(list) > 0 {
							ipMap := list[0].(map[string]interface{})
							rules := ipMap["rules"]
							result["Rules"] = rules.(*schema.Set).List()
						}
						(*call.SdkParam)["DenyIpAccessRule"] = result
					} else {
						(*call.SdkParam)["DenyIpAccessRule"] = map[string]interface{}{}
					}
				}

				if d.HasChange("allow_referer_access_rule") {
					// common type 必传
					if allowRefererAccessRule, ok := d.GetOk("allow_referer_access_rule"); ok {
						result := make(map[string]interface{})
						if list, ok := allowRefererAccessRule.([]interface{}); ok && len(list) > 0 {
							allowRefererAccessRuleMap := list[0].(map[string]interface{})
							if allowEmpty, ok := allowRefererAccessRuleMap["allow_empty"]; ok {
								result["AllowEmpty"] = allowEmpty
							}
							commonTypeResult := make(map[string]interface{})
							commonType := allowRefererAccessRuleMap["common_type"]
							if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
								commonTypeMap := commonTypeList[0].(map[string]interface{})
								if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
									commonTypeResult["IgnoreCase"] = ignoreCase
								}
								rules := commonTypeMap["rules"]
								commonTypeResult["Rules"] = rules.(*schema.Set).List()
								result["CommonType"] = commonTypeResult
							}
						}
						(*call.SdkParam)["AllowRefererAccessRule"] = result
					} else {
						(*call.SdkParam)["AllowRefererAccessRule"] = map[string]interface{}{}
					}
				}

				if d.HasChange("deny_referer_access_rule") {
					// common type 必传
					if denyRefererAccessRule, ok := d.GetOk("deny_referer_access_rule"); ok {
						result := make(map[string]interface{})
						if list, ok := denyRefererAccessRule.([]interface{}); ok && len(list) > 0 {
							denyRefererAccessRuleMap := list[0].(map[string]interface{})
							if allowEmpty, ok := denyRefererAccessRuleMap["allow_empty"]; ok {
								result["AllowEmpty"] = allowEmpty
							}
							commonTypeResult := make(map[string]interface{})
							commonType := denyRefererAccessRuleMap["common_type"]
							if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
								commonTypeMap := commonTypeList[0].(map[string]interface{})
								if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
									commonTypeResult["IgnoreCase"] = ignoreCase
								}
								rules := commonTypeMap["rules"]
								commonTypeResult["Rules"] = rules.(*schema.Set).List()
								result["CommonType"] = commonTypeResult
							}
						}
						(*call.SdkParam)["DenyRefererAccessRule"] = result
					} else {
						(*call.SdkParam)["DenyRefererAccessRule"] = map[string]interface{}{}
					}
				}

				if d.HasChange("common_match_list") {
					// common type 必传
					if common, ok := d.GetOk("common_match_list"); ok {
						result := make(map[string]interface{})
						if list, ok := common.([]interface{}); ok && len(list) > 0 {
							commonMap := list[0].(map[string]interface{})
							commonTypeResult := make(map[string]interface{})
							commonType := commonMap["common_type"]
							if commonTypeList, ok := commonType.([]interface{}); ok && len(commonTypeList) > 0 {
								commonTypeMap := commonTypeList[0].(map[string]interface{})
								if ignoreCase, ok := commonTypeMap["ignore_case"]; ok {
									commonTypeResult["IgnoreCase"] = ignoreCase
								}
								rules := commonTypeMap["rules"]
								commonTypeResult["Rules"] = rules.(*schema.Set).List()
								result["CommonType"] = commonTypeResult
							}
						}
						(*call.SdkParam)["CommonMatchList"] = result
					} else {
						(*call.SdkParam)["CommonMatchList"] = map[string]interface{}{}
					}
				}
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

func (s *ByteplusCdnSharedConfigService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSharedConfig",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"ConfigName": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnSharedConfigService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"project_name": {
				TargetField: "Project",
			},
			"config_type_list": {
				TargetField: "ConfigTypeList",
				ConvertType: bp.ConvertListN,
			},
		},
		NameField:    "ConfigName",
		IdField:      "ConfigName",
		ContentType:  bp.ContentTypeJson,
		CollectField: "config_data",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Project": {
				TargetField: "project_name",
			},
		},
	}
}

func (s *ByteplusCdnSharedConfigService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
