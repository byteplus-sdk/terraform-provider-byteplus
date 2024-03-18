package scaling_lifecycle_hook

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

type ByteplusScalingLifecycleHookService struct {
	Client *bp.SdkClient
}

func NewScalingLifecycleHookService(c *bp.SdkClient) *ByteplusScalingLifecycleHookService {
	return &ByteplusScalingLifecycleHookService{
		Client: c,
	}
}

func (s *ByteplusScalingLifecycleHookService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingLifecycleHookService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		nameSet = make(map[string]bool)
	)
	if _, ok = m["LifecycleHookNames.1"]; ok {
		i := 1
		for {
			field := fmt.Sprintf("LifecycleHookNames.%d", i)
			if name, ok := m[field]; ok {
				nameSet[name.(string)] = true
				i = i + 1
				delete(m, field)
			} else {
				break
			}
		}
	}

	hooks, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeLifecycleHooks"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, action, *resp)
		results, err = bp.ObtainSdkValue("Result.LifecycleHooks", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.LifecycleHooks is not Slice")
		}
		return data, err
	})
	if err != nil {
		return hooks, err
	}
	res := make([]interface{}, 0)
	for _, ele := range data {
		e, ok := ele.(map[string]interface{})
		if !ok {
			continue
		}
		name := e["LifecycleHookName"].(string)
		if len(nameSet) == 0 || nameSet[name] {
			res = append(res, ele)
		}
	}
	return res, nil
}

func (s *ByteplusScalingLifecycleHookService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"LifecycleHookIds.1": ids[1],
		"ScalingGroupId":     ids[0],
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
		return data, fmt.Errorf("ScalingLifecycleHook %s not exist ", ids[1])
	}
	return data, err
}

func (s *ByteplusScalingLifecycleHookService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (ByteplusScalingLifecycleHookService) WithResourceResponseHandlers(scalingGroup map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return scalingGroup, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusScalingLifecycleHookService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateLifecycleHook",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"lifecycle_command": {
					TargetField: "LifecycleCommand",
					ConvertType: bp.ConvertListUnique,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.LifecycleHookId", *resp)
				logger.Debug(logger.RespFormat, call.Action, resourceData.Get("scaling_group_id"))
				d.SetId(fmt.Sprintf("%v:%v", resourceData.Get("scaling_group_id"), id.(string)))
				return nil
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusScalingLifecycleHookService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyLifecycleHook",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"lifecycle_command": {
					Ignore: true,
				},
			},
			RequestIdField: "LifecycleHookId",
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) < 1 {
					return false, nil
				}
				(*call.SdkParam)["LifecycleHookId"] = ids[1]
				if d.HasChange("lifecycle_command") {
					commandId, ok := d.GetOk("lifecycle_command.0.command_id")
					if ok {
						(*call.SdkParam)["LifecycleCommand.CommandId"] = commandId
					} else {
						(*call.SdkParam)["LifecycleCommand.CommandId"] = ""
					}
					params, ok := d.GetOk("lifecycle_command.0.parameters")
					if ok {
						(*call.SdkParam)["LifecycleCommand.Parameters"] = params
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusScalingLifecycleHookService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteLifecycleHook",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"LifecycleHookId": ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading LifecycleHook on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusScalingLifecycleHookService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "LifecycleHookIds",
				ConvertType: bp.ConvertWithN,
			},
			"lifecycle_hook_names": {
				TargetField: "LifecycleHookNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "LifecycleHookName",
		IdField:      "LifecycleHookId",
		CollectField: "lifecycle_hooks",
		ResponseConverts: map[string]bp.ResponseConvert{
			"LifecycleHookId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusScalingLifecycleHookService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "auto_scaling",
		Action:      actionName,
		Version:     "2020-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
