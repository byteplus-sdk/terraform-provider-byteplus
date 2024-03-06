package rule

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/clb"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRuleService struct {
	Client *bp.SdkClient
}

func NewRuleService(c *bp.SdkClient) *ByteplusRuleService {
	return &ByteplusRuleService{
		Client: c,
	}
}

func (s *ByteplusRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRuleService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeRules"
		logger.Debug(logger.ReqFormat, action, condition)
		// 检查 RuleIds 是否存在
		idsMap := make(map[string]bool)
		if ids, ok := condition["RuleIds"]; ok {
			var values []interface{}
			switch _ids := ids.(type) {
			case *schema.Set:
				values = _ids.List() // from datasource
			default:
				values = _ids.([]interface{}) // from resource_read
			}
			for _, value := range values {
				if value == nil {
					continue
				}
				idsMap[strings.Trim(value.(string), " ")] = true
			}
			delete(condition, "RuleIds")
		}

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

		if len(idsMap) == 0 {
			return data, nil
		}
		// checkIds
		var res []interface{}
		for _, ele := range data {
			if _, ok := idsMap[ele.(map[string]interface{})["RuleId"].(string)]; ok {
				res = append(res, ele)
			}
		}
		return res, err
	})
}

func (s *ByteplusRuleService) ReadResource(resourceData *schema.ResourceData, ruleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	listenerId, ok := resourceData.GetOk("listener_id")
	if !ok {
		return nil, fmt.Errorf("non ListenerId")
	}
	req := map[string]interface{}{
		"ListenerId": listenerId.(string),
		"RuleIds":    []interface{}{ruleId},
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
		return data, fmt.Errorf("Rule %s not exist ", ruleId)
	}
	return data, err
}

func (s *ByteplusRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh:    nil,
	}
}

func (ByteplusRuleService) WithResourceResponseHandlers(rule map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rule, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	// 查询 LoadBalancerId
	clbId, err := s.queryLoadBalancerId(resourceData.Get("server_group_id").(string))
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRules",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"domain": {
					TargetField: "Rules.1.Domain",
				},
				"url": {
					TargetField: "Rules.1.Url",
				},
				"server_group_id": {
					TargetField: "Rules.1.ServerGroupId",
				},
				"description": {
					TargetField: "Rules.1.Description",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				ids, _ := bp.ObtainSdkValue("Result.RuleIds", *resp)
				d.SetId(ids.([]interface{})[0].(string))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	// 查询 LoadBalancerId
	clbId, err := s.queryLoadBalancerId(resourceData.Get("server_group_id").(string))
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyRules",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"server_group_id": {
					TargetField: "Rules.1.ServerGroupId",
				},
				"description": {
					TargetField: "Rules.1.Description",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Rules.1.RuleId"] = d.Id()
				(*call.SdkParam)["ListenerId"] = d.Get("listener_id")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	// 查询 LoadBalancerId
	clbId, err := s.queryLoadBalancerId(resourceData.Get("server_group_id").(string))
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRules",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"RuleIds.1":  resourceData.Id(),
				"ListenerId": resourceData.Get("listener_id"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading vpc on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "RuleIds",
			},
		},
		IdField:      "RuleId",
		CollectField: "rules",
		ResponseConverts: map[string]bp.ResponseConvert{
			"RuleId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusRuleService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusRuleService) queryLoadBalancerId(serverGroupId string) (string, error) {
	// 查询 LoadBalancerId
	action := "DescribeServerGroupAttributes"
	serverGroupResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"ServerGroupId": serverGroupId,
	})
	if err != nil {
		return "", err
	}
	clbId, err := bp.ObtainSdkValue("Result.LoadBalancerId", *serverGroupResp)
	if err != nil {
		return "", err
	}
	return clbId.(string), nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
