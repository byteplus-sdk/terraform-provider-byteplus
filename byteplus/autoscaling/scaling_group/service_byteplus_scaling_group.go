package scaling_group

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusScalingGroupService struct {
	Client *bp.SdkClient
}

func NewScalingGroupService(c *bp.SdkClient) *ByteplusScalingGroupService {
	return &ByteplusScalingGroupService{
		Client: c,
	}
}

func (s *ByteplusScalingGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingGroups"
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
		results, err = bp.ObtainSdkValue("Result.ScalingGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingGroups is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusScalingGroupService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ScalingGroupIds.1": id,
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
		return data, fmt.Errorf("ScalingGroup %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusScalingGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("LifecycleState", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("ScalingGroup  LifecycleState  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusScalingGroupService) WithResourceResponseHandlers(scalingGroup map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return scalingGroup, map[string]bp.ResponseConvert{
			"MultiAZPolicy": {
				TargetField: "multi_az_policy",
			},
			"DBInstanceIds": {
				TargetField: "db_instance_ids",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusScalingGroupService) CreateResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateScalingGroup",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"subnet_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"server_group_attributes": {
					ConvertType: bp.ConvertListN,
				},
				"min_instance_number": {
					TargetField: "MinInstanceNumber",
					// 如果为0时，需要这样转一下，要不然不会出现在请求参数
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						return i
					},
				},
				"max_instance_number": {
					TargetField: "MaxInstanceNumber",
					// 如果为0时，需要这样转一下，要不然不会出现在请求参数
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						return i
					},
				},
				"desire_instance_number": {
					TargetField: "DesireInstanceNumber",
					// 如果为0时，需要这样转一下，要不然不会出现在请求参数
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						if _, ok := data.GetOkExists("desire_instance_number"); !ok {
							return -1
						}
						return i
					},
				},
				"multi_az_policy": {
					TargetField: "MultiAZPolicy",
					ConvertType: bp.ConvertDefault,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"db_instance_ids": {
					TargetField: "DBInstanceIds",
					ConvertType: bp.ConvertWithN,
				},
				"launch_template_overrides": {
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.ScalingGroupId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"InActive"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	return callbacks
}

func (s *ByteplusScalingGroupService) ModifyResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	// 修改伸缩组
	modifyGroupCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:         "ModifyScalingGroup",
			ConvertMode:    bp.RequestConvertInConvert,
			RequestIdField: "ScalingGroupId",
			Convert: map[string]bp.RequestConvert{
				"scaling_group_name": {
					ConvertType: bp.ConvertDefault,
				},
				"min_instance_number": {
					ConvertType: bp.ConvertDefault,
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						return i
					},
				},
				"max_instance_number": {
					ConvertType: bp.ConvertDefault,
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						return i
					},
				},
				"subnet_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"desire_instance_number": {
					ConvertType: bp.ConvertDefault,
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						return i
					},
				},
				"instance_terminate_policy": {
					ConvertType: bp.ConvertDefault,
				},
				"default_cooldown": {
					ConvertType: bp.ConvertDefault,
				},
				"multi_az_policy": {
					TargetField: "MultiAZPolicy",
				},
				"launch_template_id": {
					ConvertType: bp.ConvertDefault,
				},
				"launch_template_version": {
					ConvertType: bp.ConvertDefault,
				},
				"launch_template_overrides": {
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) < 2 {
					return false, nil
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, modifyGroupCallback)
	// serverGroup modify
	attrAdd, attrRemove, _, _ := bp.GetSetDifference("server_group_attributes", resourceData, serverGroupAttributeHash, false)
	removeAttrCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachServerGroups",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if attrRemove != nil && len(attrRemove.List()) > 0 {
					(*call.SdkParam)["ScalingGroupId"] = d.Id()
					for index, attr := range attrRemove.List() {
						(*call.SdkParam)["ServerGroupAttributes."+strconv.Itoa(index+1)+".ServerGroupId"] =
							attr.(map[string]interface{})["server_group_id"].(string)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						return resource.NonRetryableError(fmt.Errorf("error on reading scaling group %q: %w", d.Id(), callErr))
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
	attachAttrCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachServerGroups",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if attrAdd != nil && len(attrAdd.List()) > 0 {
					(*call.SdkParam)["ScalingGroupId"] = d.Id()
					for index, attr := range attrAdd.List() {
						(*call.SdkParam)["ServerGroupAttributes."+strconv.Itoa(index+1)+".Port"] =
							attr.(map[string]interface{})["port"].(int)
						(*call.SdkParam)["ServerGroupAttributes."+strconv.Itoa(index+1)+".ServerGroupId"] =
							attr.(map[string]interface{})["server_group_id"].(string)
						(*call.SdkParam)["ServerGroupAttributes."+strconv.Itoa(index+1)+".Weight"] =
							attr.(map[string]interface{})["weight"].(int)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						return resource.NonRetryableError(fmt.Errorf("error on reading scaling group %q, %w", d.Id(), callErr))
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
	callbacks = append(callbacks, removeAttrCallback, attachAttrCallback)
	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "scalinggroup", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)
	return callbacks
}

func (s *ByteplusScalingGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteScalingGroup",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ScalingGroupId": resourceData.Id(),
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
							return resource.NonRetryableError(fmt.Errorf("error on reading ScalingGroup on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusScalingGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ScalingGroupIds",
				ConvertType: bp.ConvertWithN,
			},
			"scaling_group_names": {
				TargetField: "ScalingGroupNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ScalingGroupName",
		IdField:      "ScalingGroupId",
		CollectField: "scaling_groups",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ScalingGroupId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"MultiAZPolicy": {
				TargetField: "multi_az_policy",
			},
			"DBInstanceIds": {
				TargetField: "db_instance_ids",
			},
		},
	}
}

func (s *ByteplusScalingGroupService) ReadResourceId(id string) string {
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

func (s *ByteplusScalingGroupService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "auto_scaling",
		ResourceType:         "scalinggroup",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
