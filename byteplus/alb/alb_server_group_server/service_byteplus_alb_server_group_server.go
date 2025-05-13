package alb_server_group_server

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb_server_group"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusServerGroupServerService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewServerGroupServerService(c *bp.SdkClient) *ByteplusServerGroupServerService {
	return &ByteplusServerGroupServerService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusServerGroupServerService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusServerGroupServerService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	servers, err := bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		var (
			resp    *map[string]interface{}
			err     error
			results interface{}
		)
		action := "DescribeServerGroupAttributes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return nil, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return nil, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Servers", *resp)
		if err != nil {
			return []interface{}{}, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok := results.([]interface{}); !ok {
			return data, errors.New("Result.Servers is not Slice")
		} else {
			return data, err
		}
	})
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (s *ByteplusServerGroupServerService) ReadResource(resourceData *schema.ResourceData, serverGroupServerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if serverGroupServerId == "" {
		serverGroupServerId = resourceData.Id()
	}
	ids := strings.Split(serverGroupServerId, ":")
	req := map[string]interface{}{
		"ServerGroupId": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		// 找到对应的 server id
		if v.(map[string]interface{})["ServerId"] == ids[1] {
			return v.(map[string]interface{}), nil
		}
	}
	return data, fmt.Errorf("ServerGroup server %s not exist ", serverGroupServerId)
}

func (s *ByteplusServerGroupServerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (*ByteplusServerGroupServerService) WithResourceResponseHandlers(serverGroupServer map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusServerGroupServerService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddServerGroupBackendServers",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ServerGroupId"] = d.Get("server_group_id")
				(*call.SdkParam)["Servers.1.InstanceId"] = d.Get("instance_id")
				(*call.SdkParam)["Servers.1.Type"] = d.Get("type")
				(*call.SdkParam)["Servers.1.Weight"] = d.Get("weight")
				(*call.SdkParam)["Servers.1.Port"] = d.Get("port")
				(*call.SdkParam)["Servers.1.Description"] = d.Get("description")
				(*call.SdkParam)["Servers.1.Ip"] = d.Get("ip")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ServerIds.0", *resp)
				d.SetId(fmt.Sprintf("%s:%s", (*call.SdkParam)["ServerGroupId"], id.(string)))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("server_group_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb_server_group.NewAlbServerGroupService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("server_group_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusServerGroupServerService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyServerGroupBackendServers",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ServerGroupId"] = ids[0]
				(*call.SdkParam)["Servers.1.ServerId"] = ids[1]
				(*call.SdkParam)["Servers.1.Weight"] = d.Get("weight")
				(*call.SdkParam)["Servers.1.Port"] = d.Get("port")
				(*call.SdkParam)["Servers.1.Description"] = d.Get("description")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("server_group_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb_server_group.NewAlbServerGroupService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("server_group_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusServerGroupServerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveServerGroupBackendServers",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ServerGroupId": ids[0],
				"ServerIds.1":   ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除 Server Group
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("server_group_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb_server_group.NewAlbServerGroupService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("server_group_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusServerGroupServerService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "ServerId",
		IdField:      "ServerId",
		CollectField: "servers",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ServerId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusServerGroupServerService) ReadResourceId(id string) string {
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
