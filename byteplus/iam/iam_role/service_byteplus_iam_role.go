package iam_role

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusIamRoleService struct {
	Client *bp.SdkClient
}

func NewIamRoleService(c *bp.SdkClient) *ByteplusIamRoleService {
	return &ByteplusIamRoleService{
		Client: c,
	}
}

func (s *ByteplusIamRoleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusIamRoleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRoles"
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

		results, err = bp.ObtainSdkValue("Result.RoleMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.RoleMetadata is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusIamRoleService) ReadResource(resourceData *schema.ResourceData, roleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if roleId == "" {
		roleId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"RoleName": roleId,
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
		return data, fmt.Errorf("Role %s not exist ", roleId)
	}
	return data, err
}

func (s *ByteplusIamRoleService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusIamRoleService) WithResourceResponseHandlers(role map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		role["Id"] = role["RoleName"]
		return role, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusIamRoleService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRole",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				roleName, err := bp.ObtainSdkValue("Result.Role.RoleName", *resp)
				if err != nil {
					return err
				}
				d.SetId(roleName.(string))
				return nil
			},
		},
	}
	return []bp.Callback{createIamRoleCallback}
}

func (s *ByteplusIamRoleService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	updateRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateRole",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RoleName"] = d.Get("role_name")
				if maxSessionDuration, ok := d.GetOk("max_session_duration"); ok {
					(*call.SdkParam)["MaxSessionDuration"] = maxSessionDuration
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{updateRoleCallback}
}

func (s *ByteplusIamRoleService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deleteRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRole",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RoleName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading iam role on delete %q, %w", d.Id(), callErr))
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
	return []bp.Callback{deleteRoleCallback}
}

func (s *ByteplusIamRoleService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"RoleName": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
		NameField:    "RoleName",
		IdField:      "RoleName",
		CollectField: "roles",
	}
}

func (s *ByteplusIamRoleService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
		RegionType:  bp.Global,
	}
}
