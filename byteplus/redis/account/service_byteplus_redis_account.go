package account

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisAccountService struct {
	Client *bp.SdkClient
}

func NewAccountService(c *bp.SdkClient) *ByteplusRedisAccountService {
	return &ByteplusRedisAccountService{
		Client: c,
	}
}

func (s *ByteplusRedisAccountService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisAccountService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	universalClient := s.Client.UniversalClient
	action := "ListDBAccount"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
		if err != nil {
			return data, err
		}
	}

	results, err = bp.ObtainSdkValue("Result.Accounts", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.Accounts is not Slice")
	}
	return data, err
}

func (s *ByteplusRedisAccountService) ReadResource(resourceData *schema.ResourceData, RedisAccountId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if RedisAccountId == "" {
		RedisAccountId = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(RedisAccountId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid redis account id")
	}

	instanceId := ids[0]
	accountName := ids[1]

	req := map[string]interface{}{
		"InstanceId":  instanceId,
		"AccountName": accountName,
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
		return data, fmt.Errorf("Redis account %s not exist ", RedisAccountId)
	}

	return data, err
}

func (s *ByteplusRedisAccountService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      10 * time.Second,
		Pending:    []string{},
		Target:     target,
		Timeout:    timeout,
		MinTimeout: 1 * time.Second,

		Refresh: func() (result interface{}, state string, err error) {
			var (
				instance   map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "CreateFailed", "TaskFailed")

			logger.DebugInfo("start refresh :%s", id)

			ids := strings.Split(resourceData.Id(), ":")
			if len(ids) != 2 {
				return nil, "", fmt.Errorf("invalid redis account id")
			}
			instanceId := ids[0]
			if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				status, err = s.DescribeRedisInstanceStatus(instanceId)
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
			logger.DebugInfo("Refresh instance status resp: %v", instance)
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("instance %s status error, status %s", id, status.(string))
				}
			}
			logger.DebugInfo("refresh status:%v", status)
			return instance, status.(string), err
		},
	}
}

func (ByteplusRedisAccountService) WithResourceResponseHandlers(account map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return account, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusRedisAccountService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBAccount",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("instance_id"), d.Get("account_name"))
				d.SetId(id)
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisAccountService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBAccount",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				accountId := d.Id()
				ids := strings.Split(accountId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid redis account id")
				}
				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["AccountName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisAccountService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBAccount",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				redisAccountId := d.Id()
				ids := strings.Split(redisAccountId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid redis account id")
				}

				if ids[1] == "default" {
					return false, fmt.Errorf("can not delete `default` account of redis instance")
				}

				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["AccountName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				// 不能删除 default 账号
				if strings.Contains(baseErr.Error(), "can not delete `default` account of redis instance") {
					msg := fmt.Sprintf("error: %s. msg: %s",
						baseErr.Error(),
						"If you want to remove it form terraform state, "+
							"please use `terraform state rm byteplus_redis_account.resource_name` command ")
					return fmt.Errorf(msg)
				}
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading redis account on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusRedisAccountService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AccountName",
		CollectField: "accounts",
	}
}

func (s *ByteplusRedisAccountService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "redis",
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func (s *ByteplusRedisAccountService) DescribeRedisInstanceStatus(id string) (string, error) {
	var (
		results interface{}
		data    map[string]interface{}
	)
	action := "DescribeDBInstances"
	req := map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return "", err
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	results, err = bp.ObtainSdkValue("Result.Instances", *resp)
	if err != nil {
		logger.DebugInfo("bp.ObtainSdkValue return :%v", err)
		return "", err
	}
	if results == nil {
		results = []interface{}{}
	}
	instances, ok := results.([]interface{})
	if !ok {
		return "", fmt.Errorf("DescribeDBInstances responsed instances is not a slice")
	}

	for _, v := range instances {
		if data, ok = v.(map[string]interface{}); !ok {
			return "", fmt.Errorf("Value is not map ")
		}
	}

	if len(data) == 0 {
		return "", fmt.Errorf("db instance %s not exist ", id)
	}

	status, ok := data["Status"].(string)
	if !ok {
		return "", fmt.Errorf("db instance %s status is not string ", id)
	}

	return status, nil

}
