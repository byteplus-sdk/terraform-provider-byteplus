package cdn_domain

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
		action := "DescribeTemplateDomains"

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

		results, err = bp.ObtainSdkValue("Result.Domains", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Domains is not Slice")
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

	filter := map[string]interface{}{
		"Fuzzy": false,
		"Name":  "Domain",
		"Value": []string{id},
	}
	req := map[string]interface{}{
		"Filters": []interface{}{filter},
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
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddTemplateDomain",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"https_switch": {
					TargetField: "HTTPSSwitch",
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
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
	return []bp.Callback{callback}
}

func (ByteplusCdnDomainService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"HTTPSSwitch": {
				TargetField: "https_switch",
			},
			"WAFStatus": {
				TargetField: "waf_status",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnDomainService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateTemplateDomain",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"cert_id": {
					TargetField: "CertId",
					ForceGet:    true,
				},
				"cipher_template_id": {
					TargetField: "CipherTemplateId",
					ForceGet:    true,
				},
				"https_switch": {
					TargetField: "HTTPSSwitch",
					ForceGet:    true,
				},
				"service_region": {
					TargetField: "ServiceRegion",
					ForceGet:    true,
				},
				"service_template_id": {
					TargetField: "ServiceTemplateId",
					ForceGet:    true,
				},
				"project": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if d.HasChanges("service_template_id", "service_region",
					"https_switch", "cipher_template_id", "cert_id") {
					(*call.SdkParam)["Domains"] = []string{d.Id()}

					delete(*call.SdkParam, "Tags")
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
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

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, "domain", callbacks)

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
			"filters": {
				TargetField: "Filters",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Value",
						ConvertType: bp.ConvertJsonArray,
					},
				},
			},
		},
		NameField:    "Domain",
		IdField:      "Domain",
		CollectField: "domains",
		ResponseConverts: map[string]bp.ResponseConvert{
			"HTTPSSwitch": {
				TargetField: "https_switch",
			},
			"WAFStatus": {
				TargetField: "waf_status",
			},
		},
	}
}

func (s *ByteplusCdnDomainService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusCdnDomainService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["TagKeys"] = make([]string, 0)
					for _, tag := range removedTags.List() {
						(*call.SdkParam)["TagKeys"] = append((*call.SdkParam)["TagKeys"].([]string), tag.(map[string]interface{})["key"].(string))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "TagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["Tags"] = make([]map[string]interface{}, 0)
					for _, tag := range addedTags.List() {
						(*call.SdkParam)["Tags"] = append((*call.SdkParam)["Tags"].([]map[string]interface{}), tag.(map[string]interface{}))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
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

func (s *ByteplusCdnDomainService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "CDN",
		ResourceType:         "Domain",
		ProjectResponseField: "Project",
		ProjectSchemaField:   "project",
	}
}
