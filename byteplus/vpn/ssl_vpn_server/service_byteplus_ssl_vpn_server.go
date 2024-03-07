package ssl_vpn_server

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusSslVpnServerService struct {
	Client *bp.SdkClient
}

func (v *ByteplusSslVpnServerService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *ByteplusSslVpnServerService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeSslVpnServers"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		} else {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		}
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.SslVpnServers", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.SslVpnServers is not Slice")
		}
		return data, err
	})
}

func (v *ByteplusSslVpnServerService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = v.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"SslVpnServerIds.1": id,
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
		return data, fmt.Errorf("SSL Vpn Server %s not exist", id)
	}
	return data, err
}

func (v *ByteplusSslVpnServerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			demo, err = v.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, f := range failStates {
				if f == status.(string) {
					return nil, "", fmt.Errorf("VpnGateway  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (v *ByteplusSslVpnServerService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (v *ByteplusSslVpnServerService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateSslVpnServer",
			ConvertMode: bp.RequestConvertAll,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpn_gateway_id").(string)
			},
			Convert: map[string]bp.RequestConvert{
				"local_subnets": {
					ConvertType: bp.ConvertWithN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.SslVpnServerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: data.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusSslVpnServerService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifySslVpnServer",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"ssl_vpn_server_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["SslVpnServerId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: data.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusSslVpnServerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSslVpnServer",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"SslVpnServerId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := v.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading ssl vpn server on delete %q, %w", d.Id(), callErr))
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
				return bp.CheckResourceUtilRemoved(d, v.ReadResource, 3*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusSslVpnServerService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "SslVpnServerIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "SslVpnServerName",
		IdField:      "SslVpnServerId",
		CollectField: "ssl_vpn_servers",
		ResponseConverts: map[string]bp.ResponseConvert{
			"SslVpnServerId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (v *ByteplusSslVpnServerService) ReadResourceId(s string) string {
	return s
}

func NewSslVpnServerService(client *bp.SdkClient) *ByteplusSslVpnServerService {
	return &ByteplusSslVpnServerService{
		Client: client,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpn",
		Action:      actionName,
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
