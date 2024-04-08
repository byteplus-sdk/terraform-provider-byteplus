package acl

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAclService struct {
	Client *bp.SdkClient
}

func NewAclService(c *bp.SdkClient) *ByteplusAclService {
	return &ByteplusAclService{
		Client: c,
	}
}

func (s *ByteplusAclService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAclService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeAcls"
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

		results, err = bp.ObtainSdkValue("Result.Acls", *resp)
		if err != nil {
			return data, err
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Acls is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusAclService) ReadResource(resourceData *schema.ResourceData, aclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if aclId == "" {
		aclId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"AclIds.1": aclId,
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
		return data, fmt.Errorf("acl %s not exist ", aclId)
	}

	//查询属性
	var (
		resp *map[string]interface{}
	)
	action := "DescribeAclAttributes"
	condition := make(map[string]interface{})
	condition["AclId"] = aclId
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	entries, _ := bp.ObtainSdkValue("Result.AclEntries", *resp)
	logger.Debug(logger.ReqFormat, action, condition, entries)
	logger.Debug(logger.ReqFormat, action, condition, data)
	if entries != nil {
		data["AclEntries"] = entries
	}
	return data, err
}

func (s *ByteplusAclService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d      map[string]interface{}
				status interface{}
			)
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			return d, status.(string), err
		},
	}

}

func (ByteplusAclService) WithResourceResponseHandlers(acl map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return acl, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusAclService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAcl",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"acl_entries": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.AclId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}

	callbacks = append(callbacks, callback)
	//规则创建
	entryCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddAclEntries",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"acl_entries": {
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, entryCallback)
	return callbacks

}

func (s *ByteplusAclService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyAclAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"acl_entries": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AclId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	//规则修改
	add, remove, _, _ := bp.GetSetDifference("acl_entries", resourceData, bp.ClbAclEntryHash, false)

	entryRemoveCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveAclEntries",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if remove != nil && len(remove.List()) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					for index, entry := range remove.List() {
						(*call.SdkParam)["Entries."+strconv.Itoa(index+1)] = entry.(map[string]interface{})["entry"].(string)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//假如需要异步状态 这里需要等一下
				time.Sleep(time.Duration(5) * time.Second)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, entryRemoveCallback)

	entryAddCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddAclEntries",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if add != nil && len(add.List()) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					for index, entry := range add.List() {
						(*call.SdkParam)["AclEntries."+strconv.Itoa(index+1)+"."+"Entry"] = entry.(map[string]interface{})["entry"].(string)
						(*call.SdkParam)["AclEntries."+strconv.Itoa(index+1)+"."+"Description"] = entry.(map[string]interface{})["description"].(string)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, entryAddCallback)

	return callbacks
}

func (s *ByteplusAclService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAcl",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"AclId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading acl on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusAclService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "AclIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "AclName",
		IdField:      "AclId",
		CollectField: "acls",
		ResponseConverts: map[string]bp.ResponseConvert{
			"AclId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusAclService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusAclService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "clb",
		ResourceType:         "acl",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
