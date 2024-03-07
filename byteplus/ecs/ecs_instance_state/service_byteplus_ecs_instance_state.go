package ecs_instance_state

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

type ByteplusInstanceStateService struct {
	Client *bp.SdkClient
}

func NewInstanceStateService(c *bp.SdkClient) *ByteplusInstanceStateService {
	return &ByteplusInstanceStateService{
		Client: c,
	}
}

func (s *ByteplusInstanceStateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (ByteplusInstanceStateService) WithResourceResponseHandlers(subnet map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return subnet, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusInstanceStateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	targetStatus := []string{"RUNNING"}
	instanceAction := resourceData.Get("action").(string)
	if instanceAction == string(StartAction) {
		action = "StartInstance"
	} else {
		action = "StopInstance"
		targetStatus = []string{"STOPPED"}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if instanceAction == string(ForceStopAction) {
					(*call.SdkParam)["ForceStop"] = true
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				var (
					resp *map[string]interface{}
					err  error
				)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if instanceAction == string(StartAction) {
					resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				} else {
					resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				}
				logger.Debug(logger.RespFormat, call.Action, resp)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				logger.Debug(logger.RespFormat, call.Action, instanceId)
				d.SetId(fmt.Sprintf("state:%s", instanceId))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  targetStatus,
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *ByteplusInstanceStateService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) (data []interface{}, err error) {
		action := "DescribeInstances"
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
		logger.Debug(logger.RespFormat, action, resp)

		results, err = bp.ObtainSdkValue("Result.Instances", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Instances is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusInstanceStateService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		ok bool
	)
	if tmpId == "" {
		tmpId = resourceData.Id()
	}
	ids := strings.Split(tmpId, ":")
	if len(ids) != 2 {
		return nil, fmt.Errorf("invalid id format. id: %s", tmpId)
	}

	instanceId := ids[1]
	req := map[string]interface{}{
		"InstanceIds.1": instanceId,
	}

	var tempData []interface{}
	if tempData, err = s.ReadResources(req); err != nil {
		return nil, err
	}
	if len(tempData) == 0 {
		return nil, fmt.Errorf("instance %s not exist ", instanceId)
	}
	if data, ok = tempData[0].(map[string]interface{}); !ok {
		return nil, errors.New("Value is not map ")
	}

	if _, ok = resourceData.GetOk("action"); !ok {
		// check status
		status := data["Status"].(string)
		if status == "RUNNING" {
			_ = resourceData.Set("action", "Start")
		} else if status == "STOPPED" {
			_ = resourceData.Set("action", "Stop")
		} else {
			return nil, fmt.Errorf("instance %s status %s is not RUNNING or STOPPED", instanceId, status)
		}
	}
	return data, nil
}

func (s *ByteplusInstanceStateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				data       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "ERROR")
			data, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", data)
			logger.Debug(logger.ReqFormat, "DescribeInstances", data)
			logger.Debug(logger.ReqFormat, "DescribeInstances", status)
			logger.Debug(logger.ReqFormat, "DescribeInstances", target)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Ecs Instance  status  error, status:%s", status.(string))
				}
			}
			return data, status.(string), err
		},
	}
}

func (s *ByteplusInstanceStateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	targetStatus := []string{"RUNNING"}
	instanceAction := resourceData.Get("action").(string)
	if instanceAction == string(StartAction) {
		action = "StartInstance"
	} else {
		action = "StopInstance"
		targetStatus = []string{"STOPPED"}
	}

	strs := strings.Split(resourceData.Id(), ":")

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceId"] = strs[1]
				if instanceAction == string(StopAction) || instanceAction == string(ForceStopAction) {
					(*call.SdkParam)["StoppedMode"] = d.Get("stopped_mode")
				}
				if instanceAction == string(ForceStopAction) {
					(*call.SdkParam)["ForceStop"] = true
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if instanceAction == string(StartAction) {
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				} else {
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				}
			},
			Refresh: &bp.StateRefresh{
				Target:  targetStatus,
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusInstanceStateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusInstanceStateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusInstanceStateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
