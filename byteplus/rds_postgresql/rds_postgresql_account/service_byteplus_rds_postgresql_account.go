package rds_postgresql_account

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

type ByteplusRdsPostgresqlAccountService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsPostgresqlAccountService(c *bp.SdkClient) *ByteplusRdsPostgresqlAccountService {
	return &ByteplusRdsPostgresqlAccountService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsPostgresqlAccountService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsPostgresqlAccountService) ReadResources(m map[string]interface{}) ([]interface{}, error) {
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) (data []interface{}, err error) {
		var (
			resp    *map[string]interface{}
			results interface{}
			ok      bool
		)
		universalClient := s.Client.UniversalClient
		action := "DescribeDBAccounts"
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
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
	})
}

func (s *ByteplusRdsPostgresqlAccountService) ReadResource(resourceData *schema.ResourceData, accountId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		account map[string]interface{}
		ok      bool
	)
	if accountId == "" {
		accountId = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(accountId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid rds account id")
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

	for _, r := range results {
		account, ok = r.(map[string]interface{})
		if !ok {
			return data, errors.New("Value is not map ")
		}
		if accountName == account["AccountName"].(string) {
			data = account
			break
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("RDS account %s not exist ", accountId)
	}

	return data, err
}

func (s *ByteplusRdsPostgresqlAccountService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusRdsPostgresqlAccountService) WithResourceResponseHandlers(rdsAccount map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusRdsPostgresqlAccountService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBAccount",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				// 单独处理
				"account_privileges": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if d.Get("account_type").(string) == "Super" {
					if len(d.Get("account_privileges").(string)) > 0 {
						return false, fmt.Errorf(" Super account should not pass account_privileges param")
					}
				} else {
					v, ok := d.GetOkExists("account_privileges") // 没有输入使用默认值
					if ok {
						(*call.SdkParam)["AccountPrivileges"] = v
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("instance_id"), d.Get("account_name"))
				d.SetId(id)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlAccountService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	if resourceData.HasChange("account_password") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ResetDBAccount",
				ConvertMode: bp.RequestConvertIgnore,
				SdkParam: &map[string]interface{}{
					"InstanceId":      resourceData.Get("instance_id"),
					"AccountName":     resourceData.Get("account_name"),
					"AccountPassword": resourceData.Get("account_password"),
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		callbacks = append(callbacks, callback)
	}
	if resourceData.HasChange("account_privileges") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBAccountPrivilege",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"account_privileges": {
						TargetField: "AccountPrivileges",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if d.Get("account_type").(string) == "Super" {
						return false, fmt.Errorf("modification of Super account permissions is not supported")
					}
					(*call.SdkParam)["InstanceId"] = d.Get("instance_id")
					(*call.SdkParam)["AccountName"] = d.Get("account_name")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		callbacks = append(callbacks, callback)
	}
	return callbacks
}

func (s *ByteplusRdsPostgresqlAccountService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBAccount",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				rdsAccountId := d.Id()
				ids := strings.Split(rdsAccountId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid rds account id")
				}
				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["AccountName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlAccountService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AccountName",
		IdField:      "AccountName",
		CollectField: "accounts",
	}
}

func (s *ByteplusRdsPostgresqlAccountService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_postgresql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
