package cr_registry

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCrRegistryService struct {
	Client *bp.SdkClient
}

func NewCrRegistryService(c *bp.SdkClient) *ByteplusCrRegistryService {
	return &ByteplusCrRegistryService{
		Client: c,
	}
}

func (s *ByteplusCrRegistryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCrRegistryService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	pageCall := func(condition map[string]interface{}) ([]interface{}, error) {
		// Get registry
		action := "ListRegistries"
		logger.Debug(logger.ReqFormat, action, condition)
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

		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}

		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Results.Items is not slice")
		}

		for i, v := range data {
			ins := v.(map[string]interface{})
			condition := &map[string]interface{}{
				"Registry": ins["Name"],
			}

			status, err := bp.ObtainSdkValue("Status.Phase", ins)
			if err != nil {
				return data, err
			}
			if status.(string) == "Creating" || status.(string) == "Deleting" || status.(string) == "Failed" {
				logger.DebugInfo("registry status is Creating/Deleting/Failed,skip GetUser and ListDomains%s", "")
				continue
			}

			//get user
			action = "GetUser"
			logger.Debug(logger.ReqFormat, action, condition)
			resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), condition)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, condition, *resp)
			username, err := bp.ObtainSdkValue("Result.Username", *resp)
			if err != nil {
				return data, err
			}
			userStatus, err := bp.ObtainSdkValue("Result.Status", *resp)
			if err != nil {
				return data, err
			}

			data[i].(map[string]interface{})["Username"] = username
			data[i].(map[string]interface{})["UserStatus"] = userStatus

			//get domains
			action = "ListDomains"
			logger.Debug(logger.ReqFormat, action, condition)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), condition)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, condition, *resp)
			results, err = bp.ObtainSdkValue("Result.Items", *resp)
			if err != nil {
				return data, err
			}
			if results == nil {
				results = []interface{}{}
			}
			data[i].(map[string]interface{})["Domains"] = results
		}

		return data, err
	}

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, pageCall)
}

func (s *ByteplusCrRegistryService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	req := map[string]interface{}{
		"Filter": map[string]interface{}{
			"Names": []string{id},
		},
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}

	for _, v := range results {
		data, ok = v.(map[string]interface{})
		if !ok {
			return data, errors.New("value is not a map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("CrRegistry %s is not exist", id)
	}
	return data, err
}

func (s *ByteplusCrRegistryService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, name string) *resource.StateChangeConf {
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
				userStatus interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			demo, err = s.ReadResource(resourceData, name)
			if err != nil {
				return nil, "", err
			}
			logger.Debug("Refresh CrRegistry status resp:%v", "ReadResource", demo)

			status, err = bp.ObtainSdkValue("Status.Phase", demo)
			if err != nil {
				return nil, "", err
			}

			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("CrRegistry status error,status %s", status.(string))
				}
			}

			//must wait user status
			if len(target) > 0 && target[0] == "Active" {
				userStatus, err = bp.ObtainSdkValue("UserStatus", demo)
				if userStatus != "Active" {
					status = "InActive"
				} else {
					if status == "Running" {
						status = userStatus
					}
				}
			}
			return demo, status.(string), err
		},
	}
}

func (s *ByteplusCrRegistryService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, map[string]bp.ResponseConvert{
			"SkipSSLVerify": {
				TargetField: "skip_ssl_verify",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCrRegistryService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRegistry",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"project": {
					TargetField: "Project",
				},
				"resource_tags": {
					TargetField: "ResourceTags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"type": {
					TargetField: "Type",
				},
				"proxy_cache_enabled": {
					TargetField: "ProxyCacheEnabled",
				},
				"proxy_cache": {
					TargetField: "ProxyCache",
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"skip_ssl_verify": {
							TargetField: "SkipSSLVerify",
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Name"] = resourceData.Get("name")
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := d.Get("name").(string)
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	if password, ok := resourceData.GetOkExists("password"); ok {
		action := "SetUser"
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      action,
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["Registry"] = resourceData.Get("name")
					(*call.SdkParam)["Password"] = base64.StdEncoding.EncodeToString([]byte(password.(string)))
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}
	return callbacks
}

func (s *ByteplusCrRegistryService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	if resourceData.HasChange("password") {
		action := "SetUser"
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      action,
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					bytes := []byte(resourceData.Get("password").(string))
					(*call.SdkParam)["Registry"] = resourceData.Get("name")
					(*call.SdkParam)["Password"] = base64.StdEncoding.EncodeToString(bytes)
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}
	return callbacks
}

func (s *ByteplusCrRegistryService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRegistry",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Name"] = d.Id()
				(*call.SdkParam)["DeleteImmediately"] = d.Get("delete_immediately")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				deleteImmediately := d.Get("delete_immediately").(bool)
				// 如选择立即销毁，则进行removed检查
				if deleteImmediately {
					return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
				}
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCrRegistryService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		IdField:      "Name",
		CollectField: "registries",
		RequestConverts: map[string]bp.RequestConvert{
			"names": {
				TargetField: "Filter.Names",
				ConvertType: bp.ConvertJsonArray,
			},
			"types": {
				TargetField: "Filter.Types",
				ConvertType: bp.ConvertJsonArray,
			},
			"projects": {
				TargetField: "Filter.Projects",
				ConvertType: bp.ConvertJsonArray,
			},
			"statuses": {
				TargetField: "Filter.Statuses",
				ConvertType: bp.ConvertJsonObjectArray,
			},
			"resource_tags": {
				TargetField: "ResourceTagFilters",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"key": {
						TargetField: "Key",
					},
					"values": {
						TargetField: "Values",
						ConvertType: bp.ConvertJsonArray,
					},
				},
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"SkipSSLVerify": {
				TargetField: "skip_ssl_verify",
			},
		},
	}
}

func (s *ByteplusCrRegistryService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusCrRegistryService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "cr",
		ResourceType:         "instance",
		ProjectResponseField: "Project",
		ProjectSchemaField:   "project",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cr",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
