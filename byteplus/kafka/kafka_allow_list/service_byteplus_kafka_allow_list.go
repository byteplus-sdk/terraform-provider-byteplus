package kafka_allow_list

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

type ByteplusKafkaAllowListService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewKafkaAllowListService(c *bp.SdkClient) *ByteplusKafkaAllowListService {
	return &ByteplusKafkaAllowListService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusKafkaAllowListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaAllowListService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(m, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeAllowLists"

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
		results, err = bp.ObtainSdkValue("Result.AllowLists", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.AllowLists is not Slice")
		}
		for _, allowList := range data {
			allowListMap, ok := allowList.(map[string]interface{})
			if !ok {
				continue
			}
			action = "DescribeAllowListDetail"
			req := map[string]interface{}{
				"AllowListId": allowListMap["AllowListId"],
			}
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			if err != nil {
				return data, err
			}
			al, err := bp.ObtainSdkValue("Result.AllowList", *resp)
			if err != nil {
				continue
			}
			alStr, ok := al.(string)
			if ok {
				allowListMap["AllowList"] = strings.Split(alStr, ",")
			}
			instances, err := bp.ObtainSdkValue("Result.AssociatedInstances", *resp)
			if err != nil {
				continue
			}
			allowListMap["AssociatedInstances"] = instances
		}
		return data, err
	})
}

func (s *ByteplusKafkaAllowListService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results  []interface{}
		tempData map[string]interface{}
		ok       bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"RegionId": s.Client.Region,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if tempData, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else {
			alId, ok := tempData["AllowListId"].(string)
			if ok && alId == id {
				data = tempData
				break
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kafka_allow_list %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusKafkaAllowListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
	}
}

func (s *ByteplusKafkaAllowListService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAllowList",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				allowlist, ok := d.GetOk("allow_list")
				if ok {
					alSet, ok := allowlist.(*schema.Set)
					if ok {
						alList := alSet.List()
						var alStrs []string
						for _, al := range alList {
							alStr, ok := al.(string)
							if ok {
								alStrs = append(alStrs, alStr)
							}
						}
						alStr := strings.Join(alStrs, ",")
						(*call.SdkParam)["AllowList"] = alStr
					}
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
				id, _ := bp.ObtainSdkValue("Result.AllowListId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusKafkaAllowListService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaAllowListService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyAllowList",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_list_name": {
					TargetField: "AllowListName",
					ForceGet:    true,
				},
				"allow_list_desc": {
					TargetField: "AllowListDesc",
					ForceGet:    true,
				},
				"allow_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllowListId"] = d.Id()
				resp, err := s.ReadResource(d, d.Id())
				if err != nil {
					return false, err
				}
				if d.HasChange("allow_list") {
					num := resp["AssociatedInstanceNum"].(float64)
					(*call.SdkParam)["ApplyInstanceNum"] = int(num)
					allowListsSet := d.Get("allow_list")
					alSet, ok := allowListsSet.(*schema.Set)
					if ok {
						alList := alSet.List()
						var alStrs []string
						for _, al := range alList {
							alStr, ok := al.(string)
							if ok {
								alStrs = append(alStrs, alStr)
							}
						}
						alStr := strings.Join(alStrs, ",")
						(*call.SdkParam)["AllowList"] = alStr
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaAllowListService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAllowList",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"AllowListId": resourceData.Id(),
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

func (s *ByteplusKafkaAllowListService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AllowListName",
		IdField:      "AllowListId",
		CollectField: "allow_lists",
	}
}

func (s *ByteplusKafkaAllowListService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Kafka",
		Version:     "2022-05-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
