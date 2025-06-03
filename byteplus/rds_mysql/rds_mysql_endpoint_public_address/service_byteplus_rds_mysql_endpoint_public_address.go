package rds_mysql_endpoint_public_address

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlEndpointPublicAddressService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsMysqlEndpointPublicAddressService(c *bp.SdkClient) *ByteplusRdsMysqlEndpointPublicAddressService {
	return &ByteplusRdsMysqlEndpointPublicAddressService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstanceDetail"

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

func (s *ByteplusRdsMysqlEndpointPublicAddressService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results      []interface{}
		ok           bool
		temp         map[string]interface{}
		endpointData map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	// instanceId:endpointId:eipId
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"InstanceId": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else {
			if endpointId, ok := temp["EndpointId"]; ok {
				if ids[1] == endpointId.(string) {
					endpointData = temp
				}
			}
		}
	}
	if len(endpointData) == 0 {
		return data, fmt.Errorf("rds_mysql_endpoint_public_address %s not exist ", id)
	}
	logger.Debug(logger.ReqFormat, "Endpoint Data", endpointData)
	addresses := endpointData["Addresses"]
	if addresses != nil {
		for _, addr := range addresses.([]interface{}) {
			if eipId, ok := addr.(map[string]interface{})["EipId"]; ok {
				if eipId.(string) == ids[2] {
					data = addr.(map[string]interface{})
				}
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("rds_mysql_endpoint_public_address %s not exist ", id)
	}
	logger.Debug(logger.ReqFormat, "Address Data", data)
	return data, err
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				data   map[string]interface{}
				status interface{}
			)
			// 资源异步且无状态，加假状态防止读取不到
			if err = resource.Retry(10*time.Minute, func() *resource.RetryError {
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
			status = "Success"
			return data, status.(string), err
		},
	}
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBEndpointPublicAddress",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"domain": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				endpointId := d.Get("endpoint_id").(string)
				eipId := d.Get("eip_id").(string)
				d.SetId(fmt.Sprintf("%s:%s:%s", instanceId, endpointId, eipId))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Success"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					ResourceId: resourceData.Get("instance_id").(string),
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				},
				eip_address.NewEipAddressService(s.Client): {
					Target:     []string{"Attached"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("eip_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, callback)
	if domain, ok := resourceData.GetOk("domain"); ok {
		modifyCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBEndpointAddress",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Get("instance_id")
					(*call.SdkParam)["EndpointId"] = d.Get("endpoint_id")
					(*call.SdkParam)["NetworkType"] = "Public"
					arr := strings.Split(domain.(string), ".")
					if len(arr) < 2 {
						return false, fmt.Errorf("domain is not valid")
					}
					(*call.SdkParam)["DomainPrefix"] = arr[0]
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
						ResourceId: resourceData.Get("instance_id").(string),
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					},
				},
			},
		}
		callbacks = append(callbacks, modifyCallback)
	}
	return callbacks
}

func (ByteplusRdsMysqlEndpointPublicAddressService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	if resourceData.HasChange("domain") {
		modifyCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBEndpointAddress",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Get("instance_id")
					(*call.SdkParam)["EndpointId"] = d.Get("endpoint_id")
					(*call.SdkParam)["NetworkType"] = "Public"
					arr := strings.Split(d.Get("domain").(string), ".")
					if len(arr) < 2 {
						return false, fmt.Errorf("domain is not valid")
					}
					(*call.SdkParam)["DomainPrefix"] = arr[0]
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
						ResourceId: resourceData.Get("instance_id").(string),
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					},
				},
			},
		}
		callbacks = append(callbacks, modifyCallback)
	}
	return callbacks
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBEndpointPublicAddress",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": ids[0],
				"EndpointId": ids[1],
				"EipId":      ids[2],
				"Domain":     resourceData.Get("domain"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					ResourceId: resourceData.Get("instance_id").(string),
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				},
				eip_address.NewEipAddressService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("eip_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusRdsMysqlEndpointPublicAddressService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_mysql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
