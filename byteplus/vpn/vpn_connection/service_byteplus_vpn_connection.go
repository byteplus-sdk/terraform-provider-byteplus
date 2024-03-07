package vpn_connection

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

type ByteplusVpnConnectionService struct {
	Client *bp.SdkClient
}

func NewVpnConnectionService(c *bp.SdkClient) *ByteplusVpnConnectionService {
	return &ByteplusVpnConnectionService{
		Client: c,
	}
}

func (s *ByteplusVpnConnectionService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVpnConnectionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		nameSet = make(map[string]bool)
	)
	if _, ok = m["VpnConnectionNames.1"]; ok {
		i := 1
		for {
			filed := fmt.Sprintf("VpnConnectionNames.%d", i)
			tmpName, ok := m[filed]
			if !ok {
				break
			}
			nameSet[tmpName.(string)] = true
			i++
			delete(m, filed)
		}
	}
	connections, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeVpnConnections"
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
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.VpnConnections", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.VpnConnections is not Slice")
		}
		return data, err
	})
	if err != nil || len(nameSet) == 0 {
		return connections, err
	}

	res := make([]interface{}, 0)
	for _, connection := range connections {
		if !nameSet[connection.(map[string]interface{})["VpnConnectionName"].(string)] {
			continue
		}
		res = append(res, connection)
	}
	return res, nil
}

func (s *ByteplusVpnConnectionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"VpnConnectionIds.1": id,
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
		return data, fmt.Errorf("VpnConnection %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusVpnConnectionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				demo, err = s.ReadResource(resourceData, id)
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

			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("VpnConnection  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusVpnConnectionService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, map[string]bp.ResponseConvert{
			"IkeConfig.Psk": {
				TargetField: "ike_config_psk",
			},
			"IkeConfig.Version": {
				TargetField: "ike_config_version",
			},
			"IkeConfig.Mode": {
				TargetField: "ike_config_mode",
			},
			"IkeConfig.EncAlg": {
				TargetField: "ike_config_enc_alg",
			},
			"IkeConfig.AuthAlg": {
				TargetField: "ike_config_auth_alg",
			},
			"IkeConfig.DhGroup": {
				TargetField: "ike_config_dh_group",
			},
			"IkeConfig.Lifetime": {
				TargetField: "ike_config_lifetime",
			},
			"IkeConfig.LocalId": {
				TargetField: "ike_config_local_id",
			},
			"IkeConfig.RemoteId": {
				TargetField: "ike_config_remote_id",
			},
			"IpsecConfig.EncAlg": {
				TargetField: "ipsec_config_enc_alg",
			},
			"IpsecConfig.AuthAlg": {
				TargetField: "ipsec_config_auth_alg",
			},
			"IpsecConfig.DhGroup": {
				TargetField: "ipsec_config_dh_group",
			},
			"IpsecConfig.Lifetime": {
				TargetField: "ipsec_config_lifetime",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusVpnConnectionService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 创建vpnConnection
	createVpnConnection := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVpnConnection",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"local_subnet": {
					TargetField: "LocalSubnet",
					ConvertType: bp.ConvertWithN,
				},
				"remote_subnet": {
					TargetField: "RemoteSubnet",
					ConvertType: bp.ConvertWithN,
				},
				"ike_config_psk": {
					TargetField: "IkeConfig.Psk",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_version": {
					TargetField: "IkeConfig.Version",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_mode": {
					TargetField: "IkeConfig.Mode",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_enc_alg": {
					TargetField: "IkeConfig.EncAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_auth_alg": {
					TargetField: "IkeConfig.AuthAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_dh_group": {
					TargetField: "IkeConfig.DhGroup",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_lifetime": {
					TargetField: "IkeConfig.Lifetime",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_local_id": {
					TargetField: "IkeConfig.LocalId",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_remote_id": {
					TargetField: "IkeConfig.RemoteId",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_enc_alg": {
					TargetField: "IpsecConfig.EncAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_auth_alg": {
					TargetField: "IpsecConfig.AuthAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_dh_group": {
					TargetField: "IpsecConfig.DhGroup",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_lifetime": {
					TargetField: "IpsecConfig.Lifetime",
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.VpnConnectionId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("customer_gateway_id").(string)
			},
		},
	}
	callbacks = append(callbacks, createVpnConnection)

	return callbacks

}

func (s *ByteplusVpnConnectionService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 修改vpnConnection
	modifyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyVpnConnectionAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"vpn_connection_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
				"local_subnet": {
					TargetField: "LocalSubnet",
					ConvertType: bp.ConvertWithN,
				},
				"remote_subnet": {
					TargetField: "RemoteSubnet",
					ConvertType: bp.ConvertWithN,
				},
				"nat_traversal": {
					ConvertType: bp.ConvertDefault,
				},
				"dpd_action": {
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_psk": {
					TargetField: "IkeConfig.Psk",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_version": {
					TargetField: "IkeConfig.Version",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_mode": {
					TargetField: "IkeConfig.Mode",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_enc_alg": {
					TargetField: "IkeConfig.EncAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_auth_alg": {
					TargetField: "IkeConfig.AuthAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_dh_group": {
					TargetField: "IkeConfig.DhGroup",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_lifetime": {
					TargetField: "IkeConfig.Lifetime",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_local_id": {
					TargetField: "IkeConfig.LocalId",
					ConvertType: bp.ConvertDefault,
				},
				"ike_config_remote_id": {
					TargetField: "IkeConfig.RemoteId",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_enc_alg": {
					TargetField: "IpsecConfig.EncAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_auth_alg": {
					TargetField: "IpsecConfig.AuthAlg",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_dh_group": {
					TargetField: "IpsecConfig.DhGroup",
					ConvertType: bp.ConvertDefault,
				},
				"ipsec_config_lifetime": {
					TargetField: "IpsecConfig.Lifetime",
					ConvertType: bp.ConvertDefault,
				},
				"log_enabled": {
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) < 1 {
					return false, nil
				}
				(*call.SdkParam)["VpnConnectionId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("customer_gateway_id").(string)
			},
		},
	}
	callbacks = append(callbacks, modifyCallback)

	return callbacks
}

func (s *ByteplusVpnConnectionService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVpnConnection",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"VpnConnectionId": resourceData.Id(),
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
							return resource.NonRetryableError(fmt.Errorf("error on reading VpnConnection on delete %q, %w", d.Id(), callErr))
						}
					}
					resp, callErr := call.ExecuteCall(d, client, call)
					logger.Debug(logger.AllFormat, call.Action, call.SdkParam, resp, callErr)
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
				return d.Get("customer_gateway_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusVpnConnectionService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "VpnConnectionIds",
				ConvertType: bp.ConvertWithN,
			},
			"vpn_connection_names": {
				TargetField: "VpnConnectionNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "VpnConnectionName",
		IdField:      "VpnConnectionId",
		CollectField: "vpn_connections",
		ResponseConverts: map[string]bp.ResponseConvert{
			"VpnConnectionId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"IkeConfig.Psk": {
				TargetField: "ike_config_psk",
			},
			"IkeConfig.Version": {
				TargetField: "ike_config_version",
			},
			"IkeConfig.Mode": {
				TargetField: "ike_config_mode",
			},
			"IkeConfig.EncAlg": {
				TargetField: "ike_config_enc_alg",
			},
			"IkeConfig.AuthAlg": {
				TargetField: "ike_config_auth_alg",
			},
			"IkeConfig.DhGroup": {
				TargetField: "ike_config_dh_group",
			},
			"IkeConfig.Lifetime": {
				TargetField: "ike_config_lifetime",
			},
			"IkeConfig.LocalId": {
				TargetField: "ike_config_local_id",
			},
			"IkeConfig.RemoteId": {
				TargetField: "ike_config_remote_id",
			},
			"IpsecConfig.EncAlg": {
				TargetField: "ipsec_config_enc_alg",
			},
			"IpsecConfig.AuthAlg": {
				TargetField: "ipsec_config_auth_alg",
			},
			"IpsecConfig.DhGroup": {
				TargetField: "ipsec_config_dh_group",
			},
			"IpsecConfig.Lifetime": {
				TargetField: "ipsec_config_lifetime",
			},
		},
	}
}

func (s *ByteplusVpnConnectionService) ReadResourceId(id string) string {
	return id
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

func (s *ByteplusVpnConnectionService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpn",
		ResourceType:         "vpnconnection",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
