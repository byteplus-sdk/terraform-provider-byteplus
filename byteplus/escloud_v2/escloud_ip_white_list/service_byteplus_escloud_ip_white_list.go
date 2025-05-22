package escloud_ip_white_list

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/escloud_v2/escloud_instance_v2"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusEscloudIpWhiteListService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEscloudIpWhiteListService(c *bp.SdkClient) *ByteplusEscloudIpWhiteListService {
	return &ByteplusEscloudIpWhiteListService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusEscloudIpWhiteListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusEscloudIpWhiteListService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeInstances"

		// 重新组织 Filter 的格式
		if filter, filterExist := condition["Filters"]; filterExist {
			newFilter := make([]interface{}, 0)
			for k, v := range filter.(map[string]interface{}) {
				newFilter = append(newFilter, map[string]interface{}{
					"Name":   k,
					"Values": v,
				})
			}
			condition["Filters"] = newFilter
		}

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

		results, err = bp.ObtainSdkValue("Result.Instances", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Instances is not Slice")
		}

		return data, err
	})
}

func (s *ByteplusEscloudIpWhiteListService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")
	if len(ids) != 3 {
		return data, fmt.Errorf("Invalid ip white list id: %s ", id)
	}

	req := map[string]interface{}{
		"Filters": map[string]interface{}{
			"InstanceId": []string{ids[0]},
		},
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
		return data, fmt.Errorf("escloud_instance_v2 %s not exist ", id)
	}

	if ids[1] == "public" && ids[2] == "es" {
		data["IpList"] = strings.Split(data["ESPublicIpWhitelist"].(string), ",")
	} else if ids[1] == "public" && ids[2] == "kibana" {
		data["IpList"] = strings.Split(data["KibanaPublicIpWhitelist"].(string), ",")
	} else if ids[1] == "private" && ids[2] == "es" {
		data["IpList"] = strings.Split(data["ESPrivateIpWhitelist"].(string), ",")
	} else if ids[1] == "private" && ids[2] == "kibana" {
		data["IpList"] = strings.Split(data["KibanaPrivateIpWhitelist"].(string), ",")
	}

	return data, err
}

func (s *ByteplusEscloudIpWhiteListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (ByteplusEscloudIpWhiteListService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusEscloudIpWhiteListService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyIpWhitelist",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"ip_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipArr := d.Get("ip_list").(*schema.Set).List()
				ipList := make([]string, 0)
				for _, id := range ipArr {
					ipList = append(ipList, id.(string))
				}
				(*call.SdkParam)["IpList"] = strings.Join(ipList, ",")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				ipType := d.Get("type").(string)
				component := d.Get("component").(string)
				d.SetId(instanceId + ":" + ipType + ":" + component)
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				escloud_instance_v2.NewEscloudInstanceV2Service(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusEscloudIpWhiteListService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyIpWhitelist",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"instance_id": {
					TargetField: "InstanceId",
					ForceGet:    true,
				},
				"type": {
					TargetField: "Type",
					ForceGet:    true,
				},
				"component": {
					TargetField: "Component",
					ForceGet:    true,
				},
				"ip_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipArr := d.Get("ip_list").(*schema.Set).List()
				ipList := make([]string, 0)
				for _, id := range ipArr {
					ipList = append(ipList, id.(string))
				}
				(*call.SdkParam)["IpList"] = strings.Join(ipList, ",")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				escloud_instance_v2.NewEscloudInstanceV2Service(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusEscloudIpWhiteListService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyIpWhitelist",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"instance_id": {
					TargetField: "InstanceId",
					ForceGet:    true,
				},
				"type": {
					TargetField: "Type",
					ForceGet:    true,
				},
				"component": {
					TargetField: "Component",
					ForceGet:    true,
				},
				"ip_list": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["IpList"] = ""
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				escloud_instance_v2.NewEscloudInstanceV2Service(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusEscloudIpWhiteListService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusEscloudIpWhiteListService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ESCloud",
		Version:     "2023-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
