package ssl_state

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

type ByteplusMongoDBSSLStateService struct {
	Client *bp.SdkClient
}

func NewMongoDBSSLStateService(c *bp.SdkClient) *ByteplusMongoDBSSLStateService {
	return &ByteplusMongoDBSSLStateService{
		Client: c,
	}
}

func (s *ByteplusMongoDBSSLStateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongoDBSSLStateService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "ssl_state",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"SSLEnable": {
				TargetField: "ssl_enable",
			},
			"SSLExpiredTime": {
				TargetField: "ssl_expired_time",
			},
		},
	}
}

func (s *ByteplusMongoDBSSLStateService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	action := "DescribeDBInstanceSSL"
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

	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		logger.DebugInfo("bp.ObtainSdkValue return :%v", err)
		return data, err
	}
	if results == nil {
		results = map[string]interface{}{}
	}
	return []interface{}{results}, nil
}

func (s *ByteplusMongoDBSSLStateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	resourceId := resourceData.Id()
	parts := strings.Split(resourceId, ":")
	instanceId := parts[1]

	req := map[string]interface{}{
		"InstanceId": instanceId,
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
		return data, fmt.Errorf("SSLState %s is not exist", id)
	}

	expiredTime := data["SSLExpiredTime"]
	flag := certificateSetPendingRenewal(expiredTime.(string))
	if flag {
		_ = resourceData.Set("ssl_action", "EarlyRenewal")
	}

	return data, err
}

func (s *ByteplusMongoDBSSLStateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusMongoDBSSLStateService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, map[string]bp.ResponseConvert{
			"SSLEnable": {
				TargetField: "ssl_enable",
			},
			"SSLExpiredTime": {
				TargetField: "ssl_expired_time",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBSSLStateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSSL",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Get("instance_id"),
				"SSLAction":  "Open",
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				d.SetId(fmt.Sprintf("ssl:%s", d.Get("instance_id")))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				logger.Debug("lock instance id:%s", instanceId, "")
				return instanceId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBSSLStateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSSL",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Get("instance_id"),
				"SSLAction":  "Update",
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				_ = d.Set("ssl_action", "Update")
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				logger.Debug("lock instance id:%s", instanceId, "")
				return instanceId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBSSLStateService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSSL",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Get("instance_id"),
				"SSLAction":  "Close",
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				logger.Debug("lock instance id:%s", instanceId, "")
				return instanceId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBSSLStateService) ReadResourceId(id string) string {
	return id
}

func certificateSetPendingRenewal(sslExpiredTime string) bool {
	expiredTime, err := time.Parse(time.RFC3339, sslExpiredTime)
	if err != nil {
		return false
	}

	// 到期前30天执行更新操作
	earlyExpiration := expiredTime.AddDate(0, 0, -30)

	return time.Now().After(earlyExpiration)
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
