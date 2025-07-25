package kafka_instance

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

type ByteplusKafkaInstanceService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewKafkaInstanceService(c *bp.SdkClient) *ByteplusKafkaInstanceService {
	return &ByteplusKafkaInstanceService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusKafkaInstanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaInstanceService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	if v, ok := condition["Tags"]; ok {
		if len(v.(map[string]interface{})) == 0 {
			delete(condition, "Tags")
		}
	}
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeInstances"

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
		results, err = bp.ObtainSdkValue("Result.InstancesInfo", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}

		for _, element := range results.([]interface{}) {
			instance := element.(map[string]interface{})
			// 拆开 ChargeDetail
			if v, exist := instance["ChargeDetail"]; exist {
				if chargeInfo, ok := v.(map[string]interface{}); ok {
					for k, v := range chargeInfo {
						instance[k] = v
					}
				}
			}

			// update tags
			if v, ok := instance["Tags"]; ok {
				var tags []interface{}
				for k, v := range v.(map[string]interface{}) {
					tags = append(tags, map[string]interface{}{
						"Key":   k,
						"Value": v,
					})
				}
				instance["Tags"] = tags
			}

			// 获取 InstanceDetail 信息
			req := map[string]interface{}{
				"InstanceId": instance["InstanceId"],
			}
			logger.Debug(logger.ReqFormat, "DescribeInstanceDetail", req)
			detail, err := s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeInstanceDetail"), &req)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, "DescribeInstanceDetail", req, *detail)
			connection, err := bp.ObtainSdkValue("Result.ConnectionInfo", *detail)
			if err != nil {
				return data, err
			}
			instance["ConnectionInfo"] = connection
			params, err := bp.ObtainSdkValue("Result.Parameters", *detail)
			if err != nil {
				return data, err
			}
			paramsMap := make(map[string]interface{})
			if err = json.Unmarshal([]byte(params.(string)), &paramsMap); err != nil {
				return data, err
			}
			var paramsList []interface{}
			for k, v := range paramsMap {
				paramsList = append(paramsList, map[string]interface{}{
					"ParameterName":  k,
					"ParameterValue": v,
				})
			}
			instance["Parameters"] = paramsList
		}
		return results.([]interface{}), err
	})
}

func (s *ByteplusKafkaInstanceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"InstanceId": id,
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
		return data, fmt.Errorf("kafka_instance %s not exist ", id)
	}

	if zoneId, ok := data["ZoneId"]; ok {
		zoneIds := strings.Split(zoneId.(string), ",")
		data["ZoneIds"] = zoneIds
	}

	// parameters 会有默认参数，防止不一致产生
	delete(data, "Parameters")
	if parameterSet, ok := resourceData.GetOk("parameters"); ok {
		if set, ok := parameterSet.(*schema.Set); ok {
			data["Parameters"] = set.List()
		}
	}

	// 拆开 ChargeDetail
	if v, exist := data["ChargeDetail"]; exist {
		if chargeInfo, ok := v.(map[string]interface{}); ok {
			for k, v := range chargeInfo {
				data[k] = v
			}
		}
	} else {
		// 接口不返回 ChargeDetail，则回填
		data["ChargeType"] = resourceData.Get("charge_type")
		data["Period"] = resourceData.Get("period")
		data["AutoRenew"] = resourceData.Get("auto_renew")
	}

	return data, err
}

func (s *ByteplusKafkaInstanceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{},
		// 15s后才能查询 ChargeInfo
		Delay:      15 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "CreateFailed", "Error", "Fail", "Failed")

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				d, err = s.ReadResource(resourceData, id)
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

			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("InstanceStatus", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("kafka_instance status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusKafkaInstanceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateInstance",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"parameters": {
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"tags": {
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"zone_ids": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				subnetId := (*call.SdkParam)["SubnetId"].(string)
				action := "DescribeSubnetAttributes"
				req := map[string]interface{}{
					"SubnetId": subnetId,
				}
				resp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
				if err != nil {
					return false, err
				}
				logger.Debug(logger.RespFormat, action, req, *resp)
				vpcId, err := bp.ObtainSdkValue("Result.VpcId", *resp)
				if err != nil {
					return false, err
				}
				(*call.SdkParam)["VpcId"] = vpcId
				v, err := bp.ObtainSdkValue("Result.ZoneId", *resp)
				if err != nil {
					return false, err
				}
				zoneId, ok := v.(string)
				if !ok {
					return false, fmt.Errorf("Result.ZoneId is not string")
				}

				var zoneIdsStr string
				zoneIdsArr, ok := d.Get("zone_ids").([]interface{})
				if !ok {
					return false, fmt.Errorf("zone_ids is not slice")
				}
				if len(zoneIdsArr) > 0 {
					zoneIds := make([]string, 0)
					for _, id := range zoneIdsArr {
						zoneIds = append(zoneIds, id.(string))
					}
					zoneIdsStr = strings.Join(zoneIds, ",")
				} else {
					zoneIdsStr = zoneId
				}
				(*call.SdkParam)["ZoneId"] = zoneIdsStr

				// update charge info
				charge := make(map[string]interface{})
				if (*call.SdkParam)["ChargeType"] == "PrePaid" {
					if (*call.SdkParam)["Period"] == nil || (*call.SdkParam)["Period"].(int) < 1 {
						return false, fmt.Errorf("Instance Charge Type is PrePaid. Must set Period more than 1. ")
					}
					charge["PeriodUnit"] = "Month"
				}
				charge["ChargeType"] = (*call.SdkParam)["ChargeType"]
				delete(*call.SdkParam, "ChargeType")
				if v, ok := (*call.SdkParam)["AutoRenew"]; ok {
					charge["AutoRenew"] = v
					delete(*call.SdkParam, "AutoRenew")
				}
				if v, ok := (*call.SdkParam)["Period"]; ok {
					charge["Period"] = v
					delete(*call.SdkParam, "Period")
				}
				(*call.SdkParam)["ChargeInfo"] = charge
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// update tags
				if v, ok := (*call.SdkParam)["Tags"]; ok {
					tags := v.([]interface{})
					if len(tags) > 0 {
						temp := make(map[string]interface{})
						for _, ele := range tags {
							e := ele.(map[string]interface{})
							temp[e["Key"].(string)] = e["Value"]
						}
						(*call.SdkParam)["Tags"] = temp
					}
				}
				// update params
				if v, ok := (*call.SdkParam)["Parameters"]; ok {
					params := v.([]interface{})
					if len(params) > 0 {
						temp := make(map[string]interface{})
						for _, ele := range params {
							e := ele.(map[string]interface{})
							temp[e["ParameterName"].(string)] = e["ParameterValue"]
						}
						bytes, err := json.Marshal(&temp)
						if err != nil {
							return nil, err
						}
						(*call.SdkParam)["Parameters"] = string(bytes)
					}
				}

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.ReqFormat, call.Action, *resp)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.InstanceId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusKafkaInstanceService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaInstanceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var res []bp.Callback
	if resourceData.HasChange("instance_name") || resourceData.HasChange("instance_description") {
		res = append(res, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyInstanceAttributes",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"instance_name": {
						TargetField: "InstanceName",
					},
					"instance_description": {
						TargetField: "InstanceDescription",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.ReqFormat, call.Action, *resp)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		})
	}
	if resourceData.HasChanges("compute_spec", "storage_space", "partition_number") {
		res = append(res, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyInstanceSpec",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"compute_spec": {
						TargetField: "ComputeSpec",
					},
					"storage_space": {
						TargetField: "StorageSpace",
					},
					"partition_number": {
						TargetField: "PartitionNumber",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					if d.HasChange("compute_spec") { // 变更实例的计算规格时才需要选择是否再均衡
						if v, ok := d.GetOkExists("need_rebalance"); ok {
							(*call.SdkParam)["NeedRebalance"] = v
						}
						if v, ok := d.GetOkExists("rebalance_time"); ok {
							(*call.SdkParam)["RebalanceTime"] = v
						}
					}
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.ReqFormat, call.Action, *resp)
					return resp, err
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					time.Sleep(10 * time.Second)
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		})
	}

	if resourceData.HasChange("parameters") {
		parameterCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyInstanceParameters",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"parameters": {
						ConvertType: bp.ConvertJsonObjectArray,
						ForceGet:    true,
					},
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					if _, exist := (*call.SdkParam)["Parameters"]; !exist {
						return nil, nil
					}
					params := (*call.SdkParam)["Parameters"].([]interface{})
					if len(params) == 0 {
						return nil, nil
					}
					temp := make(map[string]interface{})
					for _, ele := range params {
						para := ele.(map[string]interface{})
						temp[para["ParameterName"].(string)] = para["ParameterValue"]
					}
					bytes, err := json.Marshal(&temp)
					if err != nil {
						return nil, err
					}
					(*call.SdkParam)["Parameters"] = string(bytes)
					(*call.SdkParam)["InstanceId"] = d.Id()

					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}

		res = append(res, parameterCallback)
	}
	if resourceData.HasChanges("charge_type") {
		res = append(res, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyInstanceChargeType",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					// 仅支持按量付费转包年包月
					if d.Get("charge_type") == "PostPaid" {
						return false, fmt.Errorf("onny support PostPaid to PrePaid")
					}

					if d.Get("charge_type") == "PrePaid" {
						if d.Get("period") == nil || d.Get("period").(int) < 1 {
							return false, fmt.Errorf("Instance Charge Type is PrePaid. Must set Period more than 1. ")
						}
					}

					(*call.SdkParam)["InstanceId"] = d.Id()
					charge := make(map[string]interface{})
					charge["PeriodUnit"] = "Month"
					charge["AutoRenew"] = d.Get("auto_renew")
					charge["Period"] = d.Get("period")
					charge["ChargeType"] = d.Get("charge_type")
					(*call.SdkParam)["ChargeInfo"] = charge
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		})
	}

	// 更新Tags
	res = s.setResourceTags(resourceData, res)
	return res
}

func (s *ByteplusKafkaInstanceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteInstance",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaInstanceService) setResourceTags(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveTagsFromResource",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["InstanceIds"] = []string{resourceData.Id()}
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
			Action:      "AddTagsToResource",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["InstanceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["Tags"] = make([]map[string]interface{}, 0)
					for _, tag := range addedTags.List() {
						t := tag.(map[string]interface{})
						temp := make(map[string]interface{})
						temp["Key"] = t["key"].(string)
						temp["Value"] = t["value"].(string)
						(*call.SdkParam)["Tags"] = append((*call.SdkParam)["Tags"].([]map[string]interface{}), temp)
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

func (s *ByteplusKafkaInstanceService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			// Map<String, Array of String> 类型
			"tags": {
				Convert: func(data *schema.ResourceData, i interface{}) interface{} {
					tags := i.(*schema.Set).List()
					res := make(map[string]interface{})
					for _, ele := range tags {
						tag := ele.(map[string]interface{})
						res[tag["key"].(string)] = []interface{}{tag["value"]}
					}
					return res
				},
			},
		},
		NameField:    "InstanceName",
		IdField:      "InstanceId",
		CollectField: "instances",
		ResponseConverts: map[string]bp.ResponseConvert{
			"InstanceId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusKafkaInstanceService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusKafkaInstanceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "Kafka",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusKafkaInstanceService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("charge_type") == "PrePaid" {
		info.Products = []string{"Message_Queue_for_Kafka"}
		info.NeedUnsubscribe = true
	}
	return &info, nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "kafka",
		Version:     "2022-05-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
