package waf_cc_rule

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

type ByteplusWafCcRuleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafCcRuleService(c *bp.SdkClient) *ByteplusWafCcRuleService {
	return &ByteplusWafCcRuleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafCcRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafCcRuleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "Page", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListCCRule"
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
		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Data is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusWafCcRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return data, fmt.Errorf("format of waf acl rule resource id is invalid,%s", id)
	}
	ccRuleId := parts[0]
	host := parts[1]

	ruleId, err := strconv.Atoi(ccRuleId)
	tag := fmt.Sprintf("%012d", ruleId)
	ruleTag := "E" + tag
	req := map[string]interface{}{
		"RuleTag": ruleTag,
		"Host":    host,
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
		return data, fmt.Errorf("waf_cc_rule %s not exist ", id)
	}
	if ruleGroups, ruleGroupsExist := data["RuleGroup"]; ruleGroupsExist {
		for _, ruleGroup := range ruleGroups.([]interface{}) {
			if rules, rulesExist := ruleGroup.(map[string]interface{})["Rules"]; rulesExist {
				for _, rule := range rules.([]interface{}) {
					if accurateGroup, accurateGroupExist := rule.(map[string]interface{})["AccurateGroup"]; accurateGroupExist {
						rule.(map[string]interface{})["AccurateGroup"] = []interface{}{
							accurateGroup,
						}
					}
				}
			}
		}
	}
	if ruleGroups, ruleGroupsExist := data["RuleGroup"]; ruleGroupsExist {
		for _, ruleGroup := range ruleGroups.([]interface{}) {
			if group, groupExist := ruleGroup.(map[string]interface{})["Group"]; groupExist {
				ruleGroup.(map[string]interface{})["Group"] = []interface{}{
					group,
				}
			}
		}
	}

	if ruleGroups, ruleGroupsExist := data["RuleGroup"]; ruleGroupsExist {
		for _, ruleGroup := range ruleGroups.([]interface{}) {
			if rules, rulesExist := ruleGroup.(map[string]interface{})["Rules"]; rulesExist {
				for _, rule := range rules.([]interface{}) {
					if cronEnable, cronEnableExist := rule.(map[string]interface{})["CronEnable"]; cronEnableExist {
						data["CronEnable"] = cronEnable
					}
					if exemptionTime, exemptionTimeExist := rule.(map[string]interface{})["ExemptionTime"]; exemptionTimeExist {
						data["ExemptionTime"] = exemptionTime
					}
					if cronConfs, cronConfsExist := rule.(map[string]interface{})["CronConfs"]; cronConfsExist {
						data["CronConfs"] = cronConfs
					}
					if name, nameExist := rule.(map[string]interface{})["Name"]; nameExist {
						data["Name"] = name
					}
					if ccType, ccTypeExist := rule.(map[string]interface{})["CCType"]; ccTypeExist {
						data["CCType"] = ccType
					}
					if advancedEnable, advancedEnableExist := rule.(map[string]interface{})["AdvancedEnable"]; advancedEnableExist {
						data["AdvancedEnable"] = advancedEnable
					}
					if countTime, countTimeExist := rule.(map[string]interface{})["CountTime"]; countTimeExist {
						data["CountTime"] = countTime
					}
					if effectTime, effectTimeExist := rule.(map[string]interface{})["EffectTime"]; effectTimeExist {
						data["EffectTime"] = effectTime
					}
					if enable, enableExist := rule.(map[string]interface{})["Enable"]; enableExist {
						data["Enable"] = enable
					}
					if field, fieldExist := rule.(map[string]interface{})["Field"]; fieldExist {
						data["Field"] = field
					}
					if pathThreshold, pathThresholdExist := rule.(map[string]interface{})["PathThreshold"]; pathThresholdExist {
						data["PathThreshold"] = pathThreshold
					}
					if rulePriority, rulePriorityExist := rule.(map[string]interface{})["RulePriority"]; rulePriorityExist {
						data["RulePriority"] = rulePriority
					}
					if singleThreshold, singleThresholdExist := rule.(map[string]interface{})["SingleThreshold"]; singleThresholdExist {
						data["SingleThreshold"] = singleThreshold
					}
					if accurateGroup, accurateGroupExist := rule.(map[string]interface{})["AccurateGroup"]; accurateGroupExist {
						data["AccurateGroup"] = accurateGroup
					}
				}
			}
		}
	}

	return data, err
}

func (s *ByteplusWafCcRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafCcRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCCRule",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"cc_type": {
					TargetField: "CCType",
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
				"cron_confs": {
					ConvertType: bp.ConvertJsonObjectArray,
					TargetField: "CronConfs",
					NextLevelConvert: map[string]bp.RequestConvert{
						"crontab": {
							TargetField: "Crontab",
						},
						"path_threshold": {
							TargetField: "PathThreshold",
						},
						"single_threshold": {
							TargetField: "SingleThreshold",
						},
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				host, ok := d.Get("host").(string)
				if !ok {
					return errors.New("host is not string")
				}
				d.SetId(fmt.Sprintf("%s:%s", strconv.Itoa(int(id.(float64))), host))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafCcRuleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"CCType": {
				TargetField: "cc_type",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafCcRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCCRule",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"accurate_group_priority": {
					TargetField: "AccurateGroupPriority",
					ForceGet:    true,
				},
				"cron_confs": {
					ConvertType: bp.ConvertJsonObjectArray,
					ForceGet:    true,
					TargetField: "CronConfs",
					NextLevelConvert: map[string]bp.RequestConvert{
						"crontab": {
							ForceGet:    true,
							TargetField: "Crontab",
						},
						"path_threshold": {
							ForceGet:    true,
							TargetField: "PathThreshold",
						},
						"single_threshold": {
							ForceGet:    true,
							TargetField: "SingleThreshold",
						},
					},
				},
				"cron_enable": {
					TargetField: "CronEnable",
					ForceGet:    true,
				},
				"exemption_time": {
					TargetField: "ExemptionTime",
					ForceGet:    true,
				},
				"name": {
					TargetField: "Name",
					ForceGet:    true,
				},
				"url": {
					TargetField: "Url",
					ForceGet:    true,
				},
				"advanced_enable": {
					TargetField: "AdvancedEnable",
					ForceGet:    true,
				},
				"field": {
					TargetField: "Field",
					ForceGet:    true,
				},
				"single_threshold": {
					TargetField: "SingleThreshold",
					ForceGet:    true,
				},
				"path_threshold": {
					TargetField: "PathThreshold",
					ForceGet:    true,
				},
				"count_time": {
					TargetField: "CountTime",
					ForceGet:    true,
				},
				"cc_type": {
					TargetField: "CCType",
					ForceGet:    true,
				},
				"effect_time": {
					TargetField: "EffectTime",
					ForceGet:    true,
				},
				"rule_priority": {
					TargetField: "RulePriority",
					ForceGet:    true,
				},
				"enable": {
					TargetField: "Enable",
					ForceGet:    true,
				},
				"accurate_group": {
					ConvertType: bp.ConvertJsonObject,
					ForceGet:    true,
					TargetField: "AccurateGroup",
					NextLevelConvert: map[string]bp.RequestConvert{
						"accurate_rules": {
							ConvertType: bp.ConvertJsonObjectArray,
							ForceGet:    true,
							TargetField: "AccurateRules",
							NextLevelConvert: map[string]bp.RequestConvert{
								"http_obj": {
									TargetField: "HttpObj",
									ForceGet:    true,
								},
								"obj_type": {
									TargetField: "ObjType",
									ForceGet:    true,
								},
								"opretar": {
									TargetField: "Opretar",
									ForceGet:    true,
								},
								"property": {
									TargetField: "Property",
									ForceGet:    true,
								},
								"value_string": {
									TargetField: "ValueString",
									ForceGet:    true,
								},
							},
						},
						"logic": {
							TargetField: "Logic",
							ForceGet:    true,
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {

				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf acl rule resource id is invalid,%s", d.Id())
				}
				id := parts[0]
				host := parts[1]
				ccRuleId, err := strconv.Atoi(id)
				if err != nil {
					return false, fmt.Errorf(" ccRuleId cannot convert to int ")
				}
				(*call.SdkParam)["Id"] = ccRuleId
				(*call.SdkParam)["Host"] = host
				logic, ok := d.Get("accurate_group.0.logic").(int)
				if !ok {
					return false, fmt.Errorf("accurate_group.0.logic cannot convert to int ")
				}

				if logic == 0 {
					delete(*call.SdkParam, "AccurateGroup.Logic")
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

func (s *ByteplusWafCcRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCCRule",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf acl rule resource id is invalid,%s", d.Id())
				}
				id := parts[0]
				host := parts[1]
				ccRuleId, err := strconv.Atoi(id)
				if err != nil {
					return false, fmt.Errorf(" ccRuleId cannot convert to int ")
				}
				(*call.SdkParam)["ID"] = ccRuleId
				(*call.SdkParam)["Host"] = host
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

func (s *ByteplusWafCcRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"cc_type": {
				TargetField: "CCType",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		CollectField: "data",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"cc_type": {
				TargetField: "CCType",
			},
		},
	}
}

func (s *ByteplusWafCcRuleService) ReadResourceId(id string) string {
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
