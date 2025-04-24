package rds_mysql_account

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	database "github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_database"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlAccountService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsMysqlAccountService(c *bp.SdkClient) *ByteplusRdsMysqlAccountService {
	return &ByteplusRdsMysqlAccountService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsMysqlAccountService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlAccountService) ReadResources(m map[string]interface{}) ([]interface{}, error) {
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

func (s *ByteplusRdsMysqlAccountService) ReadResource(resourceData *schema.ResourceData, RdsMysqlAccountId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		account map[string]interface{}
		ok      bool
	)
	if RdsMysqlAccountId == "" {
		RdsMysqlAccountId = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(RdsMysqlAccountId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid rds mysql account id")
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
		return data, fmt.Errorf("RDS account %s not exist ", RdsMysqlAccountId)
	}

	return data, err
}

func (s *ByteplusRdsMysqlAccountService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusRdsMysqlAccountService) WithResourceResponseHandlers(rdsAccount map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rdsAccount, map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlAccountService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBAccount",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"account_privileges": {
					TargetField: "AccountPrivileges",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"db_name": {
							TargetField: "DBName",
						},
					},
				},
			},
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				privileges := d.Get("account_privileges")
				if privileges == nil || privileges.(*schema.Set).Len() == 0 {
					return true, nil
				}
				for _, privilege := range privileges.(*schema.Set).List() {
					privilegeMap, ok := privilege.(map[string]interface{})
					if !ok {
						return false, fmt.Errorf("account_privilege is not map")
					}
					if name, ok := privilegeMap["db_name"]; ok {
						if len(name.(string)) > 0 {
							// Global模式name可以为空，这里只检查不为空的情况
							id := fmt.Sprintf("%s:%s", d.Get("instance_id"), name.(string))
							_, err := database.NewRdsMysqlDatabaseService(client).ReadResource(d, id)
							if err != nil {
								return false, err
							}
						}
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				//创建RdsMysqlAccount
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

func (s *ByteplusRdsMysqlAccountService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := []bp.Callback{}
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
				Action:      "GrantDBAccountPrivilege",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"account_privileges": {
						TargetField: "AccountPrivileges",
						ForceGet:    true,
						NextLevelConvert: map[string]bp.RequestConvert{
							"db_name": {
								TargetField: "DBName",
							},
						},
						ConvertType: bp.ConvertJsonObjectArray,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
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

		add, remove, _, _ := bp.GetSetDifference("account_privileges", resourceData, RdsMysqlAccountPrivilegeHash, false)
		if remove != nil && remove.Len() > 0 {
			removeCallback := bp.Callback{
				Call: bp.SdkCall{
					Action:      "RevokeDBAccountPrivilege",
					ConvertMode: bp.RequestConvertIgnore,
					BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
						(*call.SdkParam)["InstanceId"] = d.Get("instance_id")
						(*call.SdkParam)["AccountName"] = d.Get("account_name")
						dbNames := make([]string, 0)
						for _, item := range remove.List() {
							m, ok := item.(map[string]interface{})
							if !ok {
								continue
							}
							removeDbName := m["db_name"].(string)
							// 过滤掉有Grant操作的db_name,Grant权限方式为覆盖，先取消原有权限，再赋新权限，此处无需再取消一次。
							if add != nil && add.Len() > 0 && hasDbNameInSet(removeDbName, add) {
								continue
							}
							dbNames = append(dbNames, removeDbName)
						}
						if len(dbNames) == 0 {
							return false, nil
						}
						(*call.SdkParam)["DBNames"] = strings.Join(dbNames, ",")
						return true, nil
					},
					ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
						logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
						return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					},
				},
			}
			callbacks = append(callbacks, removeCallback)
		}
	}
	return callbacks
}

func hasDbNameInSet(dbName string, set *schema.Set) bool {
	for _, item := range set.List() {
		if m, ok := item.(map[string]interface{}); ok {
			if v, ok := m["db_name"]; ok && v.(string) == dbName {
				if detail, ok := m["account_privilege_detail"].(string); ok {
					if len(detail) == 0 && m["account_privilege"].(string) == "Custom" {
						return false
					}
				} else {
					return true
				}
				return true
			}
		}
	}
	return false
}

func (s *ByteplusRdsMysqlAccountService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
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
				//删除RdsMysqlAccount
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading RDS account on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusRdsMysqlAccountService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		NameField:    "AccountName",
		CollectField: "accounts",
		ResponseConverts: map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
		},
	}
}

func (s *ByteplusRdsMysqlAccountService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_mysql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
