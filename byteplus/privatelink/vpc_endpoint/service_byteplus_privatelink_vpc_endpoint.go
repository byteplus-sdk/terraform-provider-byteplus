package vpc_endpoint

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVpcEndpointService struct {
	Client *bp.SdkClient
}

func NewVpcEndpointService(c *bp.SdkClient) *ByteplusVpcEndpointService {
	return &ByteplusVpcEndpointService{
		Client: c,
	}
}

func (s *ByteplusVpcEndpointService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVpcEndpointService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpoints"
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
		results, err = bp.ObtainSdkValue("Result.Endpoints", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Endpoints is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusVpcEndpointService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"EndpointIds.1": id,
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

	// 查询 security_group
	action := "DescribeVpcEndpointSecurityGroups"
	condition := &map[string]interface{}{
		"EndpointId": id,
	}
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), condition)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	securityGroupIds, err := bp.ObtainSdkValue("Result.SecurityGroupIds", *resp)
	if err != nil {
		return data, err
	}
	if securityGroupIds == nil {
		securityGroupIds = []interface{}{}
	}
	data["SecurityGroupIds"] = securityGroupIds

	return data, err
}

func (s *ByteplusVpcEndpointService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				data   map[string]interface{}
				status interface{}
			)
			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				data, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", data)
			if err != nil {
				return nil, "", err
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return data, status.(string), err
		},
	}
}

func (s *ByteplusVpcEndpointService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return data, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusVpcEndpointService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVpcEndpoint",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"security_group_ids": {
					TargetField: "SecurityGroupIds",
					ConvertType: bp.ConvertWithN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// 通过安全组查询 vpc
				securityGroupId := d.Get("security_group_ids").(*schema.Set).List()[0].(string)
				action := "DescribeSecurityGroups"
				req := map[string]interface{}{
					"SecurityGroupIds.1": securityGroupId,
				}
				resp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
				if err != nil {
					return false, err
				}
				logger.Debug(logger.RespFormat, action, req, *resp)
				results, err := bp.ObtainSdkValue("Result.SecurityGroups", *resp)
				if err != nil {
					return false, err
				}
				if results == nil {
					results = []interface{}{}
				}
				securityGroups, ok := results.([]interface{})
				if !ok {
					return false, errors.New("Result.SecurityGroups is not Slice")
				}
				if len(securityGroups) == 0 {
					return false, fmt.Errorf("securityGroup %s not exist", securityGroupId)
				}
				vpcId := securityGroups[0].(map[string]interface{})["VpcId"].(string)

				(*call.SdkParam)["VpcId"] = vpcId
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.EndpointId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("service_id").(string)
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusVpcEndpointService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if resourceData.HasChange("endpoint_name") || resourceData.HasChange("description") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyVpcEndpointAttributes",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"endpoint_name": {
						TargetField: "EndpointName",
					},
					"description": {
						TargetField: "Description",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["EndpointId"] = d.Id()
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

	if resourceData.HasChange("security_group_ids") {
		add, remove, _, _ := bp.GetSetDifference("security_group_ids", resourceData, schema.HashString, false)
		for _, element := range add.List() {
			callbacks = append(callbacks, s.securityGroupActionCallback(resourceData, "AttachSecurityGroupToVpcEndpoint", element))
		}
		for _, element := range remove.List() {
			callbacks = append(callbacks, s.securityGroupActionCallback(resourceData, "DetachSecurityGroupFromVpcEndpoint", element))
		}
	}

	return callbacks
}

func (s *ByteplusVpcEndpointService) securityGroupActionCallback(resourceData *schema.ResourceData, action string, element interface{}) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["EndpointId"] = d.Id()
				(*call.SdkParam)["SecurityGroupId"] = element.(string)
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

func (s *ByteplusVpcEndpointService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVpcEndpoint",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"EndpointId": resourceData.Id(),
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
							return resource.NonRetryableError(fmt.Errorf("error on reading vpc endpoint on delete %q, %w", d.Id(), callErr))
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
			LockId: func(d *schema.ResourceData) string {
				return d.Get("service_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusVpcEndpointService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "EndpointIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "EndpointName",
		IdField:      "EndpointId",
		CollectField: "vpc_endpoints",
		ResponseConverts: map[string]bp.ResponseConvert{
			"EndpointId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusVpcEndpointService) ReadResourceId(id string) string {
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

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
