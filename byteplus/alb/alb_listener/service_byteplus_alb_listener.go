package alb_listener

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbListenerService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbListenerService(c *bp.SdkClient) *ByteplusAlbListenerService {
	return &ByteplusAlbListenerService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbListenerService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbListenerService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeListeners"

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
		results, err = bp.ObtainSdkValue("Result.Listeners", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Listeners is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusAlbListenerService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ListenerIds.1": id,
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
		return data, fmt.Errorf("alb_listener %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbListenerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Failed")
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
					return nil, "", fmt.Errorf("alb_listener status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusAlbListenerService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateListener",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ca_certificate_id": {
					TargetField: "CACertificateId",
				},
				"acl_ids": {
					ConvertType: bp.ConvertWithN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ListenerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("load_balancer_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, callback)
	if customizedCfgId, ok := resourceData.GetOk("customized_cfg_id"); ok {
		customCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyListenerAttributes",
				ConvertMode: bp.RequestConvertIgnore,
				SdkParam: &map[string]interface{}{
					"CustomizedCfgId": customizedCfgId,
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					listenerId := d.State().ID
					logger.Debug(logger.ReqFormat, "Update Customized Cfg Id", listenerId)
					(*call.SdkParam)["ListenerId"] = listenerId
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active", "Disabled"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("load_balancer_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					alb.NewAlbService(s.Client): {
						Target:     []string{"Active", "Inactive"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
						ResourceId: resourceData.Get("load_balancer_id").(string),
					},
				},
			},
		}
		callbacks = append(callbacks, customCallback)
	}
	return callbacks
}

func (ByteplusAlbListenerService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbListenerService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"listener_name": {
					ConvertType: bp.ConvertDefault,
				},
				"enabled": {
					ConvertType: bp.ConvertDefault,
				},
				"certificate_source": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"cert_center_certificate_id": {
					ConvertType: bp.ConvertDefault,
				},
				"certificated_id": {
					ConvertType: bp.ConvertDefault,
				},
				"ca_certificate_id": {
					TargetField: "CACertificateId",
				},
				"acl_ids": {
					ConvertType: bp.ConvertWithN,
					ForceGet:    true,
				},
				"server_group_id": {
					ConvertType: bp.ConvertDefault,
				},
				"enable_http2": {
					ConvertType: bp.ConvertDefault,
				},
				"enable_quic": {
					ConvertType: bp.ConvertDefault,
				},
				"customized_cfg_id": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"acl_status": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"acl_type": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ListenerId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("load_balancer_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbListenerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteListener",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ListenerId": resourceData.Id(),
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				enabled := resourceData.Get("enabled")
				if enabled.(string) == "on" {
					return false, fmt.Errorf("The listener can only be deleted when it is stopped. " +
						"Please modify the enable field to off before performing the deletion operation. ")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("load_balancer_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, callback)
	return callbacks
}

func (s *ByteplusAlbListenerService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ListenerIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ListenerName",
		IdField:      "ListenerId",
		CollectField: "listeners",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ListenerId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"CACertificateId": {
				TargetField: "ca_certificate_id",
			},
		},
	}
}

func (s *ByteplusAlbListenerService) ReadResourceId(id string) string {
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
