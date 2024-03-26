package vpc_endpoint_service

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusService struct {
	Client *bp.SdkClient
}

func NewService(c *bp.SdkClient) *ByteplusService {
	return &ByteplusService{
		Client: c,
	}
}

func (s *ByteplusService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpointServices"
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
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.VpcEndpointServices", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.VpcEndpointServices is not Slice")
		}

		for index, element := range data {
			resources, err := s.describeResources(map[string]interface{}{
				"ServiceId": element.(map[string]interface{})["ServiceId"],
			})
			if err != nil {
				return nil, err
			}
			data[index].(map[string]interface{})["Resources"] = resources
		}
		return data, err
	})
}

func (s *ByteplusService) describeResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpointServiceResources"
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
		results, err = bp.ObtainSdkValue("Result.Resources", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Resources is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ServiceIds.1": id,
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
		return data, fmt.Errorf("Vpc endpoint service %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusService) WithResourceResponseHandlers(nodePool map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func (s *ByteplusService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Error")
			data, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", data)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Vpc endpoint service status  error, status: %s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return data, status.(string), err
		},
	}
}

func (s *ByteplusService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVpcEndpointService",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"resources": {
					TargetField: "Resources",
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ServiceId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if resourceData.HasChange("description") || resourceData.HasChange("auto_accept_enabled") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyVpcEndpointServiceAttributes",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"auto_accept_enabled": {
						TargetField: "AutoAcceptEnabled",
					},
					"description": {
						TargetField: "Description",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["ServiceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("resources") {
		add, remove, _, _ := bp.GetSetDifference("resources", resourceData, resourceHash, false)
		for _, element := range add.List() {
			callbacks = append(callbacks, s.resourceActionCallback(resourceData, "AttachResourceToVpcEndpointService", element))
		}
		for _, element := range remove.List() {
			callbacks = append(callbacks, s.resourceActionCallback(resourceData, "DetachResourceFromVpcEndpointService", element))
		}
	}
	return callbacks
}

func (s *ByteplusService) resourceActionCallback(resourceData *schema.ResourceData, action string, element interface{}) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ServiceId"] = d.Id()
				(*call.SdkParam)["ResourceId"] = element.(map[string]interface{})["resource_id"]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
}

func (s *ByteplusService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVpcEndpointService",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ServiceId": resourceData.Id(),
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

func (ByteplusService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ServiceIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ServiceName",
		IdField:      "ServiceId",
		CollectField: "services",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ServiceId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "privatelink",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
		ContentType: bp.Default,
	}
}
