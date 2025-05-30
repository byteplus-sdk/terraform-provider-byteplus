package alb_rule

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb_listener"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbRuleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbRuleService(c *bp.SdkClient) *ByteplusAlbRuleService {
	return &ByteplusAlbRuleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbRuleService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		var (
			resp    *map[string]interface{}
			results interface{}
			ok      bool
		)
		action := "DescribeRules"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		}
		if err != nil {
			return nil, err
		}
		logger.Debug(logger.RespFormat, action, *resp)
		results, err = bp.ObtainSdkValue("Result.Rules", *resp)
		if err != nil {
			return []interface{}{}, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Rules is not Slice")
		} else {
			return data, err
		}
	})
}

func (s *ByteplusAlbRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		temp    map[string]interface{}
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"ListenerId": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else if temp["RuleId"].(string) == ids[1] {
			data = temp
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("alb_rule %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusAlbRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRules",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"listener_id": {
					TargetField: "ListenerId",
				},
				"domain": {
					TargetField: "Rules.1.Domain",
				},
				"url": {
					TargetField: "Rules.1.Url",
				},
				"rule_action": {
					TargetField: "Rules.1.RuleAction",
				},
				"server_group_id": {
					TargetField: "Rules.1.ServerGroupId",
				},
				"description": {
					TargetField: "Rules.1.Description",
				},
				"traffic_limit_enabled": {
					TargetField: "Rules.1.TrafficLimitEnabled",
				},
				"traffic_limit_qps": {
					TargetField: "Rules.1.TrafficLimitQPS",
				},
				"rewrite_enabled": {
					TargetField: "Rules.1.RewriteEnabled",
				},
				"rewrite_config": {
					TargetField: "Rules.1.RewriteConfig",
					ConvertType: bp.ConvertListUnique,
				},
				"redirect_config": {
					TargetField: "Rules.1.RedirectConfig",
					ConvertType: bp.ConvertListUnique,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return loadBalancerId
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				ids, _ := bp.ObtainSdkValue("Result.RuleIds", *resp)
				if len(ids.([]interface{})) < 1 {
					return fmt.Errorf("rule id not found")
				}
				ruleId := ids.([]interface{})[0].(string)
				d.SetId(fmt.Sprintf("%v:%v", listenerId, ruleId))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusAlbRuleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyRules",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"server_group_id": {
					TargetField: "Rules.1.ServerGroupId",
				},
				"description": {
					TargetField: "Rules.1.Description",
				},
				"traffic_limit_enabled": {
					TargetField: "Rules.1.TrafficLimitEnabled",
				},
				"traffic_limit_qps": {
					TargetField: "Rules.1.TrafficLimitQPS",
				},
				"rewrite_enabled": {
					TargetField: "Rules.1.RewriteEnabled",
				},
				"rewrite_config": {
					TargetField: "Rules.1.RewriteConfig",
					ConvertType: bp.ConvertListUnique,
					ForceGet:    true,
				},
				"redirect_config": {
					TargetField: "Rules.1.RedirectConfig",
					ConvertType: bp.ConvertListUnique,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["ListenerId"] = ids[0]
				(*call.SdkParam)["Rules.1.RuleId"] = ids[1]
				ruleAction, ok := d.GetOk("rule_action")
				/*
					1. ruleAction = Redirect，则redirect_config必传
					2. 若ruleAction没写，则serverGroupId必传
				*/
				if ok {
					(*call.SdkParam)["Rules.1.RuleAction"] = ruleAction
					_, ok = d.GetOk("redirect_config")
					if ruleAction.(string) == "Redirect" && !ok {
						return false, fmt.Errorf("redirect_config is required when rule_action is Redirect")
					}
				}
				return true, nil
			},
			LockId: func(d *schema.ResourceData) string {
				return loadBalancerId
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRules",
			ConvertMode: bp.RequestConvertIgnore,
			LockId: func(d *schema.ResourceData) string {
				return loadBalancerId
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["ListenerId"] = ids[0]
				(*call.SdkParam)["RuleIds.1"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		IdField:      "RuleId",
		CollectField: "rules",
		ResponseConverts: map[string]bp.ResponseConvert{
			"RuleId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"TrafficLimitQPS": {
				TargetField: "traffic_limit_qps",
			},
		},
	}
}

func (s *ByteplusAlbRuleService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "alb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
