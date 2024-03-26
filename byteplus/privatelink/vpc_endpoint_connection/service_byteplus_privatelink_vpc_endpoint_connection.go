package vpc_endpoint_connection

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

type ByteplusPrivateLinkVpcEndpointConnectionService struct {
	Client *bp.SdkClient
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpointConnections"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.EndpointConnections", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.EndpointConnections is not Slice")
		}
		return data, err
	})
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = v.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return nil, errors.New("vpc endpoint connection id err")
	}
	req := map[string]interface{}{
		"EndpointId": ids[0],
		"ServiceId":  ids[1],
	}
	results, err = v.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, r := range results {
		if data, ok = r.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("vpc endpoint connection %s not exist", id)
	}
	return data, nil
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			data, err = v.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("ConnectionStatus", data)
			if err != nil {
				return nil, "", err
			}
			for _, f := range failStates {
				if f == status.(string) {
					return nil, "", fmt.Errorf("Vpc endpoint connection status error, status: %s ", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return data, status.(string), err
		},
	}
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "EnableVpcEndpointConnection",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("endpoint_id"), d.Get("service_id"))
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Connected"},
				Timeout: data.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisableVpcEndpointConnection",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"EndpointId": resourceData.Get("endpoint_id"),
				"ServiceId":  resourceData.Get("service_id"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Rejected"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "connections",
		ResponseConverts: map[string]bp.ResponseConvert{
			"NetworkInterfaceIP": {
				TargetField: "network_interface_ip",
			},
		},
	}
}

func (v *ByteplusPrivateLinkVpcEndpointConnectionService) ReadResourceId(s string) string {
	return s
}

func NewVpcEndpointConnectionService(c *bp.SdkClient) *ByteplusPrivateLinkVpcEndpointConnectionService {
	return &ByteplusPrivateLinkVpcEndpointConnectionService{
		Client: c,
	}
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
