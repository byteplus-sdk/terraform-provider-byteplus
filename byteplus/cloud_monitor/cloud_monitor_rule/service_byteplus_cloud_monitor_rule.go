package cloud_monitor_rule

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

type ByteplusCloudMonitorRuleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCloudMonitorRuleService(c *bp.SdkClient) *ByteplusCloudMonitorRuleService {
	return &ByteplusCloudMonitorRuleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCloudMonitorRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCloudMonitorRuleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRules"
		if condition != nil {
			if ids, exist := condition["Ids"]; exist && len(ids.([]interface{})) != 0 {
				action = "ListRulesByIds"
			}
		}

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

		for _, v := range data {
			ruleMap, ok := v.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("Result.Data Rule is not map")
			}
			dimensionArr := make([]interface{}, 0)
			if dimensions, exist := ruleMap["OriginalDimensions"]; exist {
				dimensionMap, ok := dimensions.(map[string]interface{})
				if !ok {
					return data, fmt.Errorf("OriginalDimensions is not map")
				}
				for key, value := range dimensionMap {
					dimensionArr = append(dimensionArr, map[string]interface{}{
						"Key":   key,
						"Value": value,
					})
				}
			}
			ruleMap["OriginalDimensions"] = dimensionArr
		}

		return data, err
	})
}

func (s *ByteplusCloudMonitorRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Ids": []interface{}{id},
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
		return data, fmt.Errorf("cloud_monitor_rule %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCloudMonitorRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("cloud_monitor_rule status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (ByteplusCloudMonitorRuleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCloudMonitorRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRule",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"recovery_notify": {
					TargetField: "RecoveryNotify",
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enable": {
							TargetField: "Enable",
						},
					},
				},
				"alert_methods": {
					TargetField: "AlertMethods",
					ConvertType: bp.ConvertJsonArray,
				},
				"contact_group_ids": {
					TargetField: "ContactGroupIds",
					ConvertType: bp.ConvertJsonArray,
				},
				"regions": {
					TargetField: "Regions",
					ConvertType: bp.ConvertJsonArray,
				},
				"conditions": {
					TargetField: "Conditions",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"metric_name": {
							TargetField: "MetricName",
						},
						"metric_unit": {
							TargetField: "MetricUnit",
						},
						"statistics": {
							TargetField: "Statistics",
						},
						"comparison_operator": {
							TargetField: "ComparisonOperator",
						},
						"threshold": {
							TargetField: "Threshold",
						},
					},
				},
				"original_dimensions": {
					Ignore: true,
				},
				"webhook_ids": {
					TargetField: "WebhookIds",
					ConvertType: bp.ConvertJsonArray,
				},
				"no_data": {
					TargetField: "NoData",
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enable": {
							TargetField: "Enable",
						},
						"evaluation_count": {
							TargetField: "EvaluationCount",
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RuleType"] = "static"

				dimensions := d.Get("original_dimensions").(*schema.Set).List()
				dimensionMap := make(map[string]interface{})
				for _, v := range dimensions {
					dimension, ok := v.(map[string]interface{})
					if !ok {
						return false, fmt.Errorf("dimension is not map")
					}
					value := dimension["value"].(*schema.Set).List()
					dimensionMap[dimension["key"].(string)] = value
				}
				(*call.SdkParam)["OriginalDimensions"] = dimensionMap

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				data, err := bp.ObtainSdkValue("Result.Data", *resp)
				if err != nil {
					return err
				}
				dataArr, ok := data.([]interface{})
				if !ok || len(dataArr) == 0 {
					return fmt.Errorf("create cloud monitor rule failed")
				}
				d.SetId(dataArr[0].(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateRule",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"rule_name": {
					TargetField: "RuleName",
					ForceGet:    true,
				},
				"description": {
					TargetField: "Description",
					ForceGet:    true,
				},
				"namespace": {
					TargetField: "Namespace",
					ForceGet:    true,
				},
				"sub_namespace": {
					TargetField: "SubNamespace",
					ForceGet:    true,
				},
				"level": {
					TargetField: "Level",
					ForceGet:    true,
				},
				"enable_state": {
					TargetField: "EnableState",
					ForceGet:    true,
				},
				"evaluation_count": {
					TargetField: "EvaluationCount",
					ForceGet:    true,
				},
				"effect_start_at": {
					TargetField: "EffectStartAt",
					ForceGet:    true,
				},
				"effect_end_at": {
					TargetField: "EffectEndAt",
					ForceGet:    true,
				},
				"silence_time": {
					TargetField: "SilenceTime",
					ForceGet:    true,
				},
				"web_hook": {
					//TargetField: "WebHook",
					Ignore: true,
				},
				"multiple_conditions": {
					TargetField: "MultipleConditions",
					ForceGet:    true,
				},
				"condition_operator": {
					TargetField: "ConditionOperator",
					ForceGet:    true,
				},
				"recovery_notify": {
					TargetField: "RecoveryNotify",
					ConvertType: bp.ConvertJsonObject,
					ForceGet:    true,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enable": {
							TargetField: "Enable",
						},
					},
				},
				"alert_methods": {
					TargetField: "AlertMethods",
					ForceGet:    true,
					ConvertType: bp.ConvertJsonArray,
				},
				"contact_group_ids": {
					//TargetField: "ContactGroupIds",
					//ConvertType: bp.ConvertJsonArray,
					Ignore: true,
				},
				"regions": {
					TargetField: "Regions",
					ForceGet:    true,
					ConvertType: bp.ConvertJsonArray,
				},
				"conditions": {
					TargetField: "Conditions",
					ForceGet:    true,
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"metric_name": {
							TargetField: "MetricName",
						},
						"metric_unit": {
							TargetField: "MetricUnit",
						},
						"statistics": {
							TargetField: "Statistics",
						},
						"comparison_operator": {
							TargetField: "ComparisonOperator",
						},
						"threshold": {
							TargetField: "Threshold",
						},
					},
				},
				"original_dimensions": {
					Ignore: true,
				},
				"webhook_ids": {
					TargetField: "WebhookIds",
					ConvertType: bp.ConvertJsonArray,
					ForceGet:    true,
				},
				"notify_mode": {
					TargetField: "NotifyMode",
					ForceGet:    true,
				},
				"no_data": {
					TargetField: "NoData",
					ConvertType: bp.ConvertJsonObject,
					ForceGet:    true,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enable": {
							TargetField: "Enable",
						},
						"evaluation_count": {
							TargetField: "EvaluationCount",
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// 将 original_dimensions 转为 map 形式
				dimensions := d.Get("original_dimensions").(*schema.Set).List()
				dimensionMap := make(map[string]interface{})
				for _, v := range dimensions {
					dimension, ok := v.(map[string]interface{})
					if !ok {
						return false, fmt.Errorf("dimension is not map")
					}
					value := dimension["value"].(*schema.Set).List()
					dimensionMap[dimension["key"].(string)] = value
				}
				(*call.SdkParam)["OriginalDimensions"] = dimensionMap

				methods := d.Get("alert_methods").(*schema.Set).List()
				// alert_methods 包含 Webhook
				if contains("Webhook", methods) {
					if webhook, ok := d.GetOk("web_hook"); ok {
						(*call.SdkParam)["WebHook"] = webhook
					}
				}
				// alert_methods 包含 Email
				if contains("Email", methods) {
					if groupIds, ok := d.GetOk("contact_group_ids"); ok {
						(*call.SdkParam)["ContactGroupIds"] = groupIds.(*schema.Set).List()
					}
				}

				(*call.SdkParam)["Id"] = d.Id()
				(*call.SdkParam)["RuleType"] = "static"
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

func (s *ByteplusCloudMonitorRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRulesByIds",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Ids": []string{resourceData.Id()},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading cloud monitor rule on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Ids",
				ConvertType: bp.ConvertJsonArray,
			},
			"alert_state": {
				TargetField: "AlertState",
				ConvertType: bp.ConvertJsonArray,
			},
			"namespace": {
				TargetField: "Namespace",
				ConvertType: bp.ConvertJsonArray,
			},
			"level": {
				TargetField: "Level",
				ConvertType: bp.ConvertJsonArray,
			},
			"enable_state": {
				TargetField: "EnableState",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		NameField:    "RuleName",
		IdField:      "Id",
		CollectField: "rules",
		ContentType:  bp.ContentTypeJson,
	}
}

func (s *ByteplusCloudMonitorRuleService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Volc_Observe",
		Version:     "2018-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
