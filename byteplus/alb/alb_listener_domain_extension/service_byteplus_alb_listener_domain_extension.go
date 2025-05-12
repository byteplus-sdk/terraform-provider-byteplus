package alb_listener_domain_extension

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb_listener"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbListenerDomainExtensionService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbListenerDomainExtensionService(c *bp.SdkClient) *ByteplusAlbListenerDomainExtensionService {
	return &ByteplusAlbListenerDomainExtensionService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbListenerDomainExtensionService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbListenerDomainExtensionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp      *map[string]interface{}
		results   interface{}
		ok        bool
		listeners []interface{}
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
		if listeners, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Listeners is not Slice")
		}
		if len(listeners) != 1 {
			return data, nil
		}
		listenerMap := listeners[0].(map[string]interface{})
		if results, ok = listenerMap["DomainExtensions"]; !ok {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("DomainExtensions is not slice")
		}
		return data, err
	})
}

func (s *ByteplusAlbListenerDomainExtensionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
		temp    map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"ListenerIds.1": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else if temp["DomainExtensionId"].(string) == ids[1] {
			data = temp
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("alb_listener_domain_extension %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbListenerDomainExtensionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusAlbListenerDomainExtensionService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: bp.RequestConvertIgnore,
			LockId: func(d *schema.ResourceData) string {
				return loadBalancerId
			},
			Convert: map[string]bp.RequestConvert{},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ListenerId"] = d.Get("listener_id")
				(*call.SdkParam)["DomainExtensions.1.Domain"] = d.Get("domain")
				(*call.SdkParam)["DomainExtensions.1.CertificateId"] = d.Get("certificate_id")
				(*call.SdkParam)["DomainExtensions.1.Action"] = "create"
				if listener["Protocol"] == "HTTP" {
					return false, fmt.Errorf("Domain extensions only HTTPS protocol listener. ")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
			AfterRefresh: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
				// 先刷新 alb 到指定状态，再查询 DomainExtensionId
				extensionId, err := s.GetDomainExtensionId(d.Get("listener_id").(string),
					d.Get("domain").(string), d.Get("certificate_id").(string))
				if err != nil {
					return err
				}
				id := fmt.Sprint(d.Get("listener_id").(string) + ":" + extensionId)
				d.SetId(id)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbListenerDomainExtensionService) GetDomainExtensionId(listenerId, domain, certId string) (extensionId string, err error) {
	req := map[string]interface{}{
		"ListenerIds.1": listenerId,
	}
	results, err := s.ReadResources(req)
	if err != nil {
		return extensionId, err
	}
	for _, r := range results {
		if temp, ok := r.(map[string]interface{}); !ok {
			return "", errors.New("Value is not map ")
		} else if temp["Domain"].(string) == domain && temp["CertificateId"].(string) == certId {
			return temp["DomainExtensionId"].(string), nil
		}
	}
	return extensionId, errors.New("DomainExtension not exist")
}

func (ByteplusAlbListenerDomainExtensionService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbListenerDomainExtensionService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: bp.RequestConvertIgnore,
			LockId: func(d *schema.ResourceData) string {
				return loadBalancerId
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["ListenerId"] = ids[0]
				(*call.SdkParam)["DomainExtensions.1.DomainExtensionId"] = ids[1]
				(*call.SdkParam)["DomainExtensions.1.Domain"] = d.Get("domain")
				(*call.SdkParam)["DomainExtensions.1.CertificateId"] = d.Get("certificate_id")
				(*call.SdkParam)["DomainExtensions.1.Action"] = "modify"
				if listener["Protocol"] == "HTTP" {
					return false, fmt.Errorf("Domain extensions only HTTPS protocol listener. ")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbListenerDomainExtensionService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	listenerId := resourceData.Get("listener_id").(string)
	listener, _ := alb_listener.NewAlbListenerService(s.Client).ReadResource(resourceData, listenerId)
	loadBalancerId := listener["LoadBalancerId"].(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["ListenerId"] = ids[0]
				(*call.SdkParam)["DomainExtensions.1.DomainExtensionId"] = ids[1]
				(*call.SdkParam)["DomainExtensions.1.Domain"] = d.Get("domain")
				(*call.SdkParam)["DomainExtensions.1.CertificateId"] = d.Get("certificate_id")
				(*call.SdkParam)["DomainExtensions.1.Action"] = "delete"
				if listener["Protocol"] == "HTTP" {
					return false, fmt.Errorf("Domain extensions only HTTPS protocol listener. ")
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
				return loadBalancerId
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				alb.NewAlbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: loadBalancerId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbListenerDomainExtensionService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"listener_id": {
				TargetField: "ListenerIds.1",
			},
		},
		IdField:      "DomainExtensionId",
		CollectField: "domain_extensions",
		ResponseConverts: map[string]bp.ResponseConvert{
			"DomainExtensionId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusAlbListenerDomainExtensionService) ReadResourceId(id string) string {
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
