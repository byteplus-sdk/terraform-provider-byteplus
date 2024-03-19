package scaling_instance_attachment

import (
	"errors"
	"fmt"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_group"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusScalingInstanceAttachmentService struct {
	Client *bp.SdkClient
}

func NewScalingInstanceAttachmentService(c *bp.SdkClient) *ByteplusScalingInstanceAttachmentService {
	return &ByteplusScalingInstanceAttachmentService{
		Client: c,
	}
}

func (s *ByteplusScalingInstanceAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingInstanceAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingInstances"
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
		logger.Debug(logger.RespFormat, action, m, resp, condition)
		results, err = bp.ObtainSdkValue("Result.ScalingInstances", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingInstances is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusScalingInstanceAttachmentService) ReadResource(resourceData *schema.ResourceData, id string) (res map[string]interface{}, err error) {
	var (
		results    []interface{}
		data       = make(map[string]interface{})
		instanceId string
		status     string
	)
	if len(id) == 0 {
		id = resourceData.Id()
	}
	ids := strings.Split(id, ":")
	// 查询伸缩组下所有实例id
	results, err = s.ReadResources(map[string]interface{}{
		"ScalingGroupId": ids[0],
		"InstanceIds.1":  ids[1],
	})
	if err != nil {
		return data, err
	}
	if len(results) == 0 {
		return data, errors.New("instance not found")
	}
	tempData, ok := results[0].(map[string]interface{})
	if !ok {
		return data, errors.New("value is not map")
	}
	instanceId = tempData["InstanceId"].(string)
	status = tempData["Status"].(string)
	data["InstanceId"] = instanceId
	data["Status"] = status
	data["Entrusted"] = tempData["Entrusted"]

	return data, nil
}

func (s *ByteplusScalingInstanceAttachmentService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("instance status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusScalingInstanceAttachmentService) WithResourceResponseHandlers(scalingGroup map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return scalingGroup, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusScalingInstanceAttachmentService) CreateResource(d *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := d.Get("instance_id").(string)
	return s.attachInstances(d, d.Get("scaling_group_id").(string), instanceId, d.Timeout(schema.TimeoutUpdate))
}

func (s *ByteplusScalingInstanceAttachmentService) ModifyResource(d *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingInstanceAttachmentService) RemoveResource(d *schema.ResourceData, r *schema.Resource) []bp.Callback {
	instanceId := d.Get("instance_id").(string)
	deleteType := d.Get("delete_type").(string)
	detachOption := d.Get("detach_option").(string)
	return s.removeInstances(d, d.Get("scaling_group_id").(string), instanceId, deleteType, detachOption)
}

func (s *ByteplusScalingInstanceAttachmentService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusScalingInstanceAttachmentService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusScalingInstanceAttachmentService) attachInstances(d *schema.ResourceData, groupId string, instanceId string, timeout time.Duration) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	attachCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachInstances",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				logger.Debug(logger.RespFormat, call.Action, instanceId)
				param := formatInstanceIdsRequest(instanceId)
				param["ScalingGroupId"] = groupId
				if entrusted, ok := d.GetOk("entrusted"); ok {
					param["Entrusted"] = entrusted
				}
				*call.SdkParam = param
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				common, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				if err != nil {
					return common, err
				}
				time.Sleep(10 * time.Second) // attach以后需要等一下
				return common, nil
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["ScalingGroupId"], ":", (*call.SdkParam)["InstanceIds.1"]))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("scaling_group_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"InService", "Protected"},
				Timeout: timeout,
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				scaling_group.NewScalingGroupService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    d.Timeout(schema.TimeoutCreate),
					ResourceId: d.Get("scaling_group_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, attachCallback)
	return callbacks
}

func (s *ByteplusScalingInstanceAttachmentService) removeInstances(d *schema.ResourceData, groupId string, instanceId string, deleteType, detachOption string) []bp.Callback {
	var action string
	if deleteType == "Detach" {
		action = "DetachInstances"
	} else {
		// 默认remove
		action = "RemoveInstances"
	}
	if detachOption != "none" {
		detachOption = "both"
	}
	callbacks := make([]bp.Callback, 0)
	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				param := formatInstanceIdsRequest(instanceId)
				param["ScalingGroupId"] = groupId
				if action == "DetachInstances" {
					param["DetachOption"] = detachOption
				}
				*call.SdkParam = param
				return true, nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("scaling_group_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				common, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				if err != nil {
					return common, err
				}
				time.Sleep(10 * time.Second) // remove以后需要等一下
				return common, nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				scaling_group.NewScalingGroupService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    d.Timeout(schema.TimeoutCreate),
					ResourceId: d.Get("scaling_group_id").(string),
				},
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, d.Id())
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading scaling instance on delete #{d.Id()}, #{callErr}"))
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
	callbacks = append(callbacks, removeCallback)
	return callbacks
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
