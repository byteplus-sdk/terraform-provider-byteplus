package cloud_monitor_object_group

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCloudMonitorObjectGroupService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCloudMonitorObjectGroupService(c *bp.SdkClient) *ByteplusCloudMonitorObjectGroupService {
	return &ByteplusCloudMonitorObjectGroupService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCloudMonitorObjectGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCloudMonitorObjectGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListObjectGroups"

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

		for _, v := range data {
			groupMap, ok := v.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("Result.Data Rule is not map")
			}
			if objects, exist := groupMap["Objects"]; exist {
				objectArr, ok := objects.([]interface{})
				if !ok {
					return data, fmt.Errorf("Objects is not slice ")
				}
				for _, object := range objectArr {
					objectMap, ok := object.(map[string]interface{})
					if !ok {
						return data, fmt.Errorf("Result.Data Object is not map")
					}
					region := objectMap["Region"]
					regions := strings.Split(region.(string), ",")
					objectMap["Region"] = regions

					dimensionArr := make([]interface{}, 0)
					if dimensions, exist := objectMap["Dimensions"]; exist {
						dimensionMap, ok := dimensions.(map[string]interface{})
						if !ok {
							return data, fmt.Errorf("Dimensions is not map ")
						}
						for key, value := range dimensionMap {
							dimensionArr = append(dimensionArr, map[string]interface{}{
								"Key":   key,
								"Value": value,
							})
						}
					}
					objectMap["Dimensions"] = dimensionArr
				}
			}
		}

		return data, err
	})
}

func (s *ByteplusCloudMonitorObjectGroupService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Ids": []interface{}{id},
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
		return data, fmt.Errorf("cloud_monitor_object_group %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCloudMonitorObjectGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (ByteplusCloudMonitorObjectGroupService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCloudMonitorObjectGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateObjectGroup",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"objects": {
					TargetField: "Objects",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"region": {
							TargetField: "Region",
							ConvertType: bp.ConvertJsonArray,
						},
						"dimensions": {
							TargetField: "Dimensions",
							ForceGet:    true,
							ConvertType: bp.ConvertJsonObjectArray,
						},
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				if objects, exist := (*call.SdkParam)["Objects"]; exist {
					objectArr, ok := objects.([]interface{})
					if !ok {
						return nil, fmt.Errorf("Objects is not slice ")
					}
					for _, object := range objectArr {
						if objectMap, ok := object.(map[string]interface{}); ok {
							if region, ok := objectMap["Region"].([]interface{}); ok {
								regions := make([]string, 0)
								for _, v := range region {
									regions = append(regions, v.(string))
								}
								objectMap["Region"] = strings.Join(regions, ",")
							}
							if dimensions, ok := objectMap["Dimensions"].([]interface{}); ok {
								dimensionMap := make(map[string]interface{})
								for _, v := range dimensions {
									dimension, ok := v.(map[string]interface{})
									if !ok {
										return nil, fmt.Errorf("dimension is not map")
									}
									value := dimension["Value"]
									dimensionMap[dimension["Key"].(string)] = value
								}
								objectMap["Dimensions"] = dimensionMap
							}
							//objectMap["Type"] = "enum"
						}
					}
				}

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Data", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorObjectGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateObjectGroup",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"name": {
					TargetField: "Name",
					ForceGet:    true,
				},
				"objects": {
					TargetField: "Objects",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"namespace": {
							TargetField: "Namespace",
							ForceGet:    true,
						},
						"dimensions": {
							TargetField: "Dimensions",
							ConvertType: bp.ConvertJsonObjectArray,
							ForceGet:    true,
						},
						"region": {
							TargetField: "Region",
							ConvertType: bp.ConvertJsonArray,
							ForceGet:    true,
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Id"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				if objects, exist := (*call.SdkParam)["Objects"]; exist {
					objectArr, ok := objects.([]interface{})
					if !ok {
						return nil, fmt.Errorf("Objects is not slice ")
					}
					for _, object := range objectArr {
						if objectMap, ok := object.(map[string]interface{}); ok {
							if region, ok := objectMap["Region"].([]interface{}); ok {
								regions := make([]string, 0)
								for _, v := range region {
									regions = append(regions, v.(string))
								}
								objectMap["Region"] = strings.Join(regions, ",")
							}
							if dimensions, ok := objectMap["Dimensions"].([]interface{}); ok {
								dimensionMap := make(map[string]interface{})
								for _, v := range dimensions {
									dimension, ok := v.(map[string]interface{})
									if !ok {
										return nil, fmt.Errorf("dimension is not map")
									}
									value := dimension["Value"]
									dimensionMap[dimension["Key"].(string)] = value
								}
								objectMap["Dimensions"] = dimensionMap
							}
							//objectMap["Type"] = "enum"
						}
					}
				}

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorObjectGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteObjectGroup",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Id": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading cloud monitor resource group on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorObjectGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Ids",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		NameField:    "Name",
		IdField:      "Id",
		CollectField: "object_groups",
		ContentType:  bp.ContentTypeJson,
	}
}

func (s *ByteplusCloudMonitorObjectGroupService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Volc_Observe",
		Version:     "2018-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
