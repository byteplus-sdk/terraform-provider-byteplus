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

type ByteplusMongoDBAccountService struct {
	Client *bp.SdkClient
}

func NewMongoDBAccountService(c *bp.SdkClient) *ByteplusMongoDBAccountService {
	return &ByteplusMongoDBAccountService{
		Client: c,
	}
}

func (s *ByteplusMongoDBAccountService) GetClient() *bp.SdkClient {
	return s.Client
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "mongodb",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func (s *ByteplusMongoDBAccountService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 10, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		var (
			resp    *map[string]interface{}
			results interface{}
			ok      bool
			err     error
			data    []interface{}
		)

		action := "DescribeDBAccounts"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		}
		if err != nil {
			return nil, err
		}

		logger.Debug(logger.RespFormat, action, condition, *resp)

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

func (s *ByteplusMongoDBAccountService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		account map[string]interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")
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
		return data, fmt.Errorf("mongodb account %s not exist ", id)
	}

	if privileges, ok := data["AccountPrivileges"].([]interface{}); ok {
		var privilegeArr []interface{}
		tempMap := make(map[string]interface{})
		for _, v := range privileges {
			privilegeMap, ok := v.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("The value of the AccountPrivileges is not map ")
			}
			dbName := privilegeMap["DBName"].(string)
			roleName := privilegeMap["RoleName"]
			if value, exist := tempMap[dbName]; exist {
				roleNames, ok := value.([]interface{})
				if !ok {
					return data, fmt.Errorf("The value of the PrivilegeMap is not slice ")
				}
				roleNames = append(roleNames, roleName)
				tempMap[dbName] = roleNames
			} else {
				roleNames := make([]interface{}, 0)
				roleNames = append(roleNames, roleName)
				tempMap[dbName] = roleNames
			}
		}
		for k, v := range tempMap {
			privilege := make(map[string]interface{})
			privilege["DbName"] = k
			privilege["RoleNames"] = v
			privilegeArr = append(privilegeArr, privilege)
		}
		data["AccountPrivileges"] = privilegeArr
	}

	return data, err
}

func (s *ByteplusMongoDBAccountService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusMongoDBAccountService) WithResourceResponseHandlers(account map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return account, map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
			"AuthDB": {
				TargetField: "auth_db",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBAccountService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBAccount",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"auth_db": {
					TargetField: "AuthDB",
				},
				"account_privileges": {
					TargetField: "AccountDBPrivileges",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"db_name": {
							TargetField: "DBName",
						},
						"role_names": {
							TargetField: "RoleNames",
							ConvertType: bp.ConvertJsonArray,
						},
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
				instanceId := d.Get("instance_id").(string)
				accountName := d.Get("account_name").(string)
				d.SetId(instanceId + ":" + accountName)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBAccountService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	if resourceData.HasChange("account_password") {
		passwordCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ResetDBAccount",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"auth_db": {
						TargetField: "AuthDB",
						ForceGet:    true,
					},
					"account_password": {
						TargetField: "AccountPassword",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					ids := strings.Split(d.Id(), ":")
					if len(ids) != 2 {
						return false, fmt.Errorf("Invalid account id: %s ", d.Id())
					}

					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["AccountName"] = ids[1]
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
		callbacks = append(callbacks, passwordCallback)
	}

	if resourceData.HasChanges("account_privileges", "account_desc") {
		privilegeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateDBAccountPrivilege",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"account_desc": {
						TargetField: "AccountDesc",
					},
					"account_privileges": {
						Ignore: true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					ids := strings.Split(d.Id(), ":")
					if len(ids) != 2 {
						return false, fmt.Errorf("Invalid account id: %s ", d.Id())
					}

					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["AccountName"] = ids[1]
					(*call.SdkParam)["AuthDB"] = d.Get("auth_db")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					privilegeArr := make([]interface{}, 0)
					privilegeMap := make(map[string]interface{})
					add, remove, _, _ := bp.GetSetDifference("account_privileges", resourceData, MongoDBAccountPrivilegeHash, false)
					if remove != nil && remove.Len() > 0 {
						for _, v := range remove.List() {
							privilege, ok := v.(map[string]interface{})
							if !ok {
								return nil, fmt.Errorf("The value of account_privileges is not map ")
							}
							dbName := privilege["db_name"].(string)
							privilegeMap[dbName] = []interface{}{}
						}
					}
					if add != nil && add.Len() > 0 {
						for _, v := range add.List() {
							privilege, ok := v.(map[string]interface{})
							if !ok {
								return nil, fmt.Errorf("The value of account_privileges is not map ")
							}
							dbName := privilege["db_name"].(string)
							privilegeMap[dbName] = privilege["role_names"].(*schema.Set).List()
						}
					}
					for k, v := range privilegeMap {
						privilege := make(map[string]interface{})
						privilege["DbName"] = k
						privilege["RoleNames"] = v
						privilegeArr = append(privilegeArr, privilege)
					}
					(*call.SdkParam)["AccountDBPrivileges"] = privilegeArr

					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
			},
		}
		callbacks = append(callbacks, privilegeCallback)
	}

	return callbacks
}

func (s *ByteplusMongoDBAccountService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBAccount",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Id": resourceData.Id(),
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("Invalid account id: %s ", d.Id())
				}
				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["AccountName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading mongodb account on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusMongoDBAccountService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"auth_db": {
				TargetField: "AuthDB",
			},
		},
		NameField:    "AccountName",
		CollectField: "accounts",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
			"AuthDB": {
				TargetField: "auth_db",
			},
		},
	}
}

func (s *ByteplusMongoDBAccountService) ReadResourceId(id string) string {
	return id
}
