package direct_connect_connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusDirectConnectConnectionService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewDirectConnectConnectionService(c *bp.SdkClient) *ByteplusDirectConnectConnectionService {
	return &ByteplusDirectConnectConnectionService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusDirectConnectConnectionService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusDirectConnectConnectionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDirectConnectConnections"

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
		results, err = bp.ObtainSdkValue("Result.DirectConnectConnections", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.DirectConnectConnections is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusDirectConnectConnectionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"DirectConnectConnectionIds.1": id,
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
		return data, fmt.Errorf("direct_connect_connection %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusDirectConnectConnectionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("direct_connect_connection status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusDirectConnectConnectionService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDirectConnectConnection",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "Tags",
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
				id, _ := bp.ObtainSdkValue("Result.DirectConnectConnectionId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Initial"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusDirectConnectConnectionService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusDirectConnectConnectionService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDirectConnectConnectionAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert:     map[string]bp.RequestConvert{},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["DirectConnectConnectionId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available", "Initial"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusDirectConnectConnectionService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDirectConnectConnection",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"DirectConnectConnectionId": resourceData.Id(),
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

func (s *ByteplusDirectConnectConnectionService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "DirectConnectConnectionIds",
				ConvertType: bp.ConvertWithN,
			},
			"tag_filters": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Values.1",
					},
				},
			},
		},
		NameField:    "DirectConnectConnectionName",
		IdField:      "DirectConnectConnectionId",
		CollectField: "direct_connect_connections",
		ResponseConverts: map[string]bp.ResponseConvert{
			"DirectConnectConnectionId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusDirectConnectConnectionService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "directconnect",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
