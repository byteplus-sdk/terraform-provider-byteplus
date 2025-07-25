package kafka_topic

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/kafka/kafka_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusKafkaTopicService struct {
	Client *bp.SdkClient
}

func NewKafkaTopicService(c *bp.SdkClient) *ByteplusKafkaTopicService {
	return &ByteplusKafkaTopicService{
		Client: c,
	}
}

func (s *ByteplusKafkaTopicService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaTopicService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeTopics"

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
		results, err = bp.ObtainSdkValue("Result.TopicsInfo", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.TopicsInfo is not Slice")
		}

		for _, ele := range data {
			topic, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" Topic is not Map ")
			}
			// 查询参数信息
			action := "DescribeTopicParameters"
			req := map[string]interface{}{
				"InstanceId": m["InstanceId"],
				"TopicName":  topic["TopicName"],
			}
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
			paramStr, err := bp.ObtainSdkValue("Result.Parameters", *resp)
			if err != nil {
				return data, err
			}
			param := make(map[string]interface{})
			err = json.Unmarshal([]byte(paramStr.(string)), &param)
			if err != nil {
				return data, fmt.Errorf(" json Unmarshal Parameters error: %v", err)
			}
			logger.DebugInfo(" Unmarshal Parameters", param)
			param["MinInsyncReplicaNumber"], _ = strconv.Atoi(param["MinInsyncReplicaNumber"].(string))
			param["MessageMaxByte"], _ = strconv.Atoi(param["MessageMaxByte"].(string))
			param["LogRetentionHours"], _ = strconv.Atoi(param["LogRetentionHours"].(string))
			topic["Parameters"] = param

			// 查询权限信息
			action = "DescribeTopicAccessPolicies"
			con := map[string]interface{}{
				"InstanceId": m["InstanceId"],
				"TopicName":  topic["TopicName"],
			}
			if userName, exist := m["UserName"]; exist && (len(userName.(string)) > 0) {
				con["UserName"] = userName
			}
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &con)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, req, *resp)

			accessPolicies, err := bp.ObtainSdkValue("Result", *resp)
			if err != nil {
				return data, err
			}
			apMap, ok := accessPolicies.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" Result is not Map ")
			}
			for k, v := range apMap {
				topic[k] = v
			}
		}

		return data, err
	})
}

func (s *ByteplusKafkaTopicService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf(" the id format must be 'instance_id:topic_name'")
	}
	req := map[string]interface{}{
		"InstanceId": ids[0],
		"TopicName":  ids[1],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		topicMap := make(map[string]interface{})
		if topicMap, ok = v.(map[string]interface{}); !ok {
			return nil, errors.New("Value is not map ")
		}
		if topicMap["TopicName"].(string) == ids[1] {
			data = topicMap
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kafka_topic %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusKafkaTopicService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Fault")
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
					return nil, "", fmt.Errorf("kafka_topic status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (ByteplusKafkaTopicService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaTopicService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateTopic",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"all_authority": {
					TargetField: "AllAuthority",
					ForceGet:    true,
				},
				"parameters": {
					TargetField: "Parameters",
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"min_insync_replica_number": {
							TargetField: "MinInsyncReplicaNumber",
						},
						"message_max_byte": {
							TargetField: "MessageMaxByte",
						},
						"log_retention_hours": {
							TargetField: "LogRetentionHours",
						},
					},
				},
				"access_policies": {
					TargetField: "AccessPolicies",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"user_name": {
							TargetField: "UserName",
						},
						"access_policy": {
							TargetField: "AccessPolicy",
						},
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// 转换 Parameters
				if param, exist := (*call.SdkParam)["Parameters"]; exist {
					paramMap, ok := param.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf(" Parameters is not map ")
					}
					for key, value := range paramMap {
						paramMap[key] = strconv.Itoa(value.(int))
					}
					paramBytes, err := json.Marshal(paramMap)
					if err != nil {
						return nil, fmt.Errorf(" Marshal Parameters error: %v", err)
					}
					logger.DebugInfo("Marshal Parameters", string(paramBytes))
					(*call.SdkParam)["Parameters"] = string(paramBytes)
				}
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%v:%v", d.Get("instance_id"), d.Get("topic_name")))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaTopicService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	ids := strings.Split(resourceData.Id(), ":")

	if resourceData.HasChange("description") {
		topicCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyTopicAttributes",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"description": {
						TargetField: "Description",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["InstanceId"] = ids[0]
						(*call.SdkParam)["TopicName"] = ids[1]
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
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					kafka_instance.NewKafkaInstanceService(s.Client): {
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
						ResourceId: resourceData.Get("instance_id").(string),
					},
				},
			},
		}
		callbacks = append(callbacks, topicCallback)
	}

	if resourceData.HasChanges("partition_number", "parameters", "replica_number") {
		paramCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyTopicParameters",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"replica_number": {
						TargetField: "ReplicaNumber",
					},
					"partition_number": {
						TargetField: "PartitionNumber",
					},
					"parameters": {
						TargetField: "Parameters",
						ConvertType: bp.ConvertJsonObject,
						ForceGet:    true,
						NextLevelConvert: map[string]bp.RequestConvert{
							"min_insync_replica_number": {
								TargetField: "MinInsyncReplicaNumber",
							},
							"message_max_byte": {
								TargetField: "MessageMaxByte",
							},
							"log_retention_hours": {
								TargetField: "LogRetentionHours",
							},
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["InstanceId"] = ids[0]
						(*call.SdkParam)["TopicName"] = ids[1]
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					// 转换 Parameters
					if param, exist := (*call.SdkParam)["Parameters"]; exist {
						paramMap, ok := param.(map[string]interface{})
						if !ok {
							return nil, fmt.Errorf(" Parameters is not map ")
						}
						for key, value := range paramMap {
							paramMap[key] = strconv.Itoa(value.(int))
						}
						paramBytes, err := json.Marshal(paramMap)
						if err != nil {
							return nil, fmt.Errorf(" Marshal Parameters error: %v", err)
						}
						logger.DebugInfo("Marshal Parameters", string(paramBytes))
						(*call.SdkParam)["Parameters"] = string(paramBytes)
					}
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					kafka_instance.NewKafkaInstanceService(s.Client): {
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
						ResourceId: resourceData.Get("instance_id").(string),
					},
				},
			},
		}
		callbacks = append(callbacks, paramCallback)
	}

	if resourceData.HasChanges("all_authority", "access_policies") {
		added, removed, _, _ := bp.GetSetDifference("access_policies", resourceData, kafkaAccessPolicyHash, false)

		callbacks = append(callbacks, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyTopicAccessPolicies",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"all_authority": {
						TargetField: "AllAuthority",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["TopicName"] = ids[1]
					if (*call.SdkParam)["AllAuthority"].(bool) {
						return true, nil
					}

					(*call.SdkParam)["AccessPolicies"] = make([]interface{}, 0)
					(*call.SdkParam)["DeletePolicies"] = make([]string, 0)
					userNames := make(map[string]bool)
					if added != nil && len(added.List()) > 0 {
						for _, ele := range added.List() {
							(*call.SdkParam)["AccessPolicies"] = append((*call.SdkParam)["AccessPolicies"].([]interface{}),
								map[string]interface{}{
									"UserName":     ele.(map[string]interface{})["user_name"],
									"AccessPolicy": ele.(map[string]interface{})["access_policy"],
								})
							userNames[ele.(map[string]interface{})["user_name"].(string)] = true
						}
					}
					if removed != nil && len(removed.List()) > 0 {
						for _, ele := range removed.List() {
							if _, exist := userNames[ele.(map[string]interface{})["user_name"].(string)]; exist {
								continue
							}
							(*call.SdkParam)["DeletePolicies"] = append((*call.SdkParam)["DeletePolicies"].([]string), ele.(map[string]interface{})["user_name"].(string))
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
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					kafka_instance.NewKafkaInstanceService(s.Client): {
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
						ResourceId: resourceData.Get("instance_id").(string),
					},
				},
			},
		})
	}
	return callbacks
}

func (s *ByteplusKafkaTopicService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteTopic",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": ids[0],
				"TopicName":  ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading kafka topic on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaTopicService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		IdField:      "TopicId",
		NameField:    "TopicName",
		CollectField: "topics",
		ExtraData: func(i []interface{}) ([]interface{}, error) {
			for index, ele := range i {
				element := ele.(map[string]interface{})
				i[index].(map[string]interface{})["TopicId"] = fmt.Sprintf("%v-%v", element["InstanceId"], element["TopicName"])
			}
			return i, nil
		},
	}
}

func (s *ByteplusKafkaTopicService) ReadResourceId(id string) string {
	return id
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
