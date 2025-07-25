package kafka_sasl_user

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/kafka/kafka_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusKafkaSaslUserService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewKafkaSaslUserService(c *bp.SdkClient) *ByteplusKafkaSaslUserService {
	return &ByteplusKafkaSaslUserService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusKafkaSaslUserService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaSaslUserService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeUsers"

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
		results, err = bp.ObtainSdkValue("Result.UsersInfo", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.UsersInfo is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusKafkaSaslUserService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"InstanceId": ids[0],
		"UserName":   ids[1],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if _, ok = v.(map[string]interface{}); !ok {
			return nil, errors.New("Value is not map ")
		}
		if v.(map[string]interface{})["UserName"] == ids[1] { // 通过名称匹配
			data = v.(map[string]interface{})
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kafka_sasl_user %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusKafkaSaslUserService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusKafkaSaslUserService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	// mlp 定制化
	instanceId := resourceData.Get("instance_id")
	userName := resourceData.Get("user_name")
	_, err := s.ReadResource(resourceData, fmt.Sprintf("%v:%v", instanceId, userName))
	if err == nil {
		// 事先有用户，先删除
		deleteCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "DeleteUser",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				SdkParam: &map[string]interface{}{
					"InstanceId": instanceId,
					"UserName":   userName,
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
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
		callbacks = append(callbacks, deleteCallback)
	}
	// mpl 定制化

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateUser",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"all_authority": {
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						v, ok := data.GetOkExists("all_authority")
						if !ok {
							return false
						}
						return v
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%v:%v", d.Get("instance_id"), d.Get("user_name")))
				return nil
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
	callbacks = append(callbacks, callback)
	return callbacks
}

func (ByteplusKafkaSaslUserService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaSaslUserService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	res := make([]bp.Callback, 0)
	ids := strings.Split(resourceData.Id(), ":")
	if resourceData.HasChange("all_authority") {
		res = append(res, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyUserAuthority",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"all_authority": {
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							v, ok := data.GetOkExists("all_authority")
							if !ok {
								return false
							}
							return v
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["UserName"] = ids[1]
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
	return res
}

func (s *ByteplusKafkaSaslUserService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteUser",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": ids[0],
				"UserName":   ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
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

func (s *ByteplusKafkaSaslUserService) DatasourceResources(d *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "UserName",
		IdField:      "UserId",
		CollectField: "users",
		ExtraData: func(sourceData []interface{}) ([]interface{}, error) {
			var next []interface{}
			for _, i := range sourceData {
				v := i.(map[string]interface{})
				v["UserId"] = fmt.Sprintf("%s:%s", d.Get("instance_id"), v["UserName"])
				next = append(next, i)
			}
			return next, nil
		},
	}
}

func (s *ByteplusKafkaSaslUserService) ReadResourceId(id string) string {
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
