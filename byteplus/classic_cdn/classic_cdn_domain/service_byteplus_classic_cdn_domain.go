package classic_cdn_domain

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

type ByteplusCdnDomainService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnDomainService(c *bp.SdkClient) *ByteplusCdnDomainService {
	return &ByteplusCdnDomainService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnDomainService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnDomainService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNum", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListCdnDomains"

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
		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Data is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusCdnDomainService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Domain": id,
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
		return data, fmt.Errorf("cdn_domain %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnDomainService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
					return nil, "", fmt.Errorf("cdn_domain status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusCdnDomainService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddCdnDomain",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "ResourceTags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"shared_cname": {
					TargetField: "SharedCname",
					ConvertType: bp.ConvertJsonObject,
				},
				"domain_config": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var (
					config map[string]interface{}
				)
				domainConfig, ok := d.Get("domain_config").(string)
				if !ok {
					return false, errors.New("domain config is not a map")
				}
				err := json.Unmarshal([]byte(domainConfig), &config)
				if err != nil || len(config) == 0 {
					return false, errors.New("domain config err or is empty")
				}
				for k, v := range config {
					(*call.SdkParam)[k] = v
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := d.Get("domain").(string)
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"online", "offline"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	return callbacks
}

func (ByteplusCdnDomainService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"ResourceTags": {
				TargetField: "tags",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnDomainService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCdnConfig",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var (
					config map[string]interface{}
				)
				(*call.SdkParam)["Domain"] = d.Id()
				domainConfig, ok := d.Get("domain_config").(string)
				if !ok {
					return false, errors.New("domain config is not a map")
				}
				err := json.Unmarshal([]byte(domainConfig), &config)
				if err != nil || len(config) == 0 {
					return false, errors.New("domain config err or is empty")
				}
				for k, v := range config {
					(*call.SdkParam)[k] = v
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"online", "offline"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	// 更新tag
	addTags, removeTags, _, _ := bp.GetSetDifference("tags", resourceData, TagsHash, false)
	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteResourceTags",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removeTags != nil && removeTags.Len() > 0 {
					(*call.SdkParam)["Resources"] = []string{d.Id()}
					tags := make([]interface{}, 0)
					for _, tag := range removeTags.List() {
						tagMap := tag.(map[string]interface{})
						tags = append(tags, map[string]interface{}{
							"Key":   tagMap["key"].(string),
							"Value": tagMap["value"].(string),
						})
					}
					(*call.SdkParam)["ResourceTags"] = tags
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	callbacks = append(callbacks, removeCallback)
	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddResourceTags",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addTags != nil && addTags.Len() > 0 {
					(*call.SdkParam)["Resources"] = []string{d.Id()}
					tags := make([]interface{}, 0)
					for _, tag := range addTags.List() {
						tagMap := tag.(map[string]interface{})
						tags = append(tags, map[string]interface{}{
							"Key":   tagMap["key"].(string),
							"Value": tagMap["value"].(string),
						})
					}
					(*call.SdkParam)["ResourceTags"] = tags
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	callbacks = append(callbacks, addCallback)
	return callbacks
}

func (s *ByteplusCdnDomainService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	status := resourceData.Get("status").(string)
	if status == "online" {
		stopCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "StopCdnDomain",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				SdkParam: &map[string]interface{}{
					"Domain": resourceData.Id(),
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"offline"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, stopCallback)
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCdnDomain",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Domain": resourceData.Id(),
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
	callbacks = append(callbacks, callback)
	return callbacks
}

func (s *ByteplusCdnDomainService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"tags": {
				TargetField: "ResourceTags",
				ConvertType: bp.ConvertJsonArray,
			},
			"ipv6": {
				TargetField: "IPv6",
			},
			"https": {
				TargetField: "HTTPS",
			},
		},
		IdField:      "Domain",
		CollectField: "domains",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"HTTPS": {
				TargetField: "https",
			},
			"IPv6": {
				TargetField: "ipv6",
			},
			"ResourceTags": {
				TargetField: "tags",
			},
		},
	}
}

func (s *ByteplusCdnDomainService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
