package instance_parameter

import (
	"errors"
	"fmt"
	"strings"
	"time"

	mongodbInstance "github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/instance"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusMongoDBInstanceParameterService struct {
	Client *bp.SdkClient
}

func NewMongoDBInstanceParameterService(c *bp.SdkClient) *ByteplusMongoDBInstanceParameterService {
	return &ByteplusMongoDBInstanceParameterService{
		Client: c,
	}
}

func (s *ByteplusMongoDBInstanceParameterService) GetClient() *bp.SdkClient {
	return s.Client
}

// ReadAll 保证data兼容性
func (s *ByteplusMongoDBInstanceParameterService) ReadAll(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	action := "DescribeDBInstanceParameters"

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

	logger.Debug(logger.RespFormat, action, resp)
	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = map[string]interface{}{}
	}
	data = []interface{}{results}
	return data, err
}

func (s *ByteplusMongoDBInstanceParameterService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action := "DescribeDBInstanceParameters"

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

	logger.Debug(logger.RespFormat, action, resp)
	results, err = bp.ObtainSdkValue("Result.InstanceParameters", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.InstanceParameters is not slice")
	}
	return data, err
}

func (s *ByteplusMongoDBInstanceParameterService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	parts := strings.Split(id, ":")
	// 兼容处理 id 为 param:instanceId:parameterName 的情况
	if len(parts) != 4 && len(parts) != 3 {
		return data, fmt.Errorf("the format of import id must be 'param:instanceId:parameterName:parameterRole'")
	}
	if len(parts) == 3 {
		role := resourceData.Get("parameter_role").(string)
		if role == "" {
			return data, fmt.Errorf("the format of import id must be 'param:instanceId:parameterName:parameterRole'")
		}
		parts = append(parts, role)
	}
	req := map[string]interface{}{
		"InstanceId":     parts[1],
		"ParameterNames": parts[2],
		"ParameterRole":  parts[3],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("instance parameters %s is not exist", id)
	}
	if _, ok = data["ParameterNames"]; ok {
		data["ParameterName"] = data["ParameterNames"]
	}
	return data, nil
}

func (s *ByteplusMongoDBInstanceParameterService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusMongoDBInstanceParameterService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBInstanceParameterService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceParameters",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"parameter_name": {
					Ignore: true,
				},
				"parameter_role": {
					Ignore: true,
				},
				"parameter_value": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ParametersObject"] = map[string]interface{}{
					"ParameterName":  d.Get("parameter_name"),
					"ParameterRole":  d.Get("parameter_role"),
					"ParameterValue": d.Get("parameter_value"),
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				parameterName := d.Get("parameter_name").(string)
				parameterRole := d.Get("parameter_role").(string)
				id := fmt.Sprintf("%v:%v:%v:%v", "param", instanceId, parameterName, parameterRole)
				d.SetId(id)
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return resourceData.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBInstanceParameterService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	id := s.ReadResourceId(resourceData.Id())
	parts := strings.Split(id, ":")
	instanceId := parts[1]
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceParameters",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceId"] = instanceId
				(*call.SdkParam)["ParametersObject"] = map[string]interface{}{
					"ParameterName":  parts[2],
					"ParameterRole":  d.Get("parameter_role"),
					"ParameterValue": d.Get("parameter_value"),
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBInstanceParameterService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongoDBInstanceParameterService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "instance_parameters",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"ParameterNames": {
				TargetField: "parameter_name",
			},
		},
	}
}

func (s *ByteplusMongoDBInstanceParameterService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "mongodb",
		Action:      actionName,
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
	}
}
