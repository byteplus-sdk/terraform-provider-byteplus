package alb_acl

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
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAclService(c *bp.SdkClient) *ByteplusAclService {
	return &ByteplusAclService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
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
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Acls", *resp)
		if err != nil {
			return data, err
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Acls is not Slice")
		}

		for index, ele := range data {
			acl := ele.(map[string]interface{})
			query := map[string]interface{}{
				"AclId": acl["AclId"],
			}

			logger.Debug(logger.ReqFormat, "DescribeAclAttributes", query)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeAclAttributes"), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, "DescribeAclAttributes", query, *resp)

			ls, err := bp.ObtainSdkValue("Result.Listeners", *resp)
			if err != nil {
				return data, err
			}
			ae, err := bp.ObtainSdkValue("Result.AclEntries", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["Listeners"] = ls
			data[index].(map[string]interface{})["AclEntries"] = ae
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
	return data, err
}

func (s *ByteplusAclService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusAclService) WithResourceResponseHandlers(acl map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}

}

func (s *ByteplusAclService) CreateResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
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
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(2 * time.Second)
				id, _ := bp.ObtainSdkValue("Result.AclId", *resp)
				d.SetId(id.(string))
				return nil
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
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(2 * time.Second) // 均为异步操作，无法获取 Status
				return nil
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
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"acl_name": {
					TargetField: "AclName",
				},
				"description": {
					TargetField: "Description",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) == 0 {
					return false, nil
				}
				(*call.SdkParam)["AclId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, callback)

	//规则修改
	add, remove, _, _ := bp.GetSetDifference("acl_entries", resourceData, AclEntryHash, false)

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
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(2 * time.Second)
				return nil
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
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(2 * time.Second)
				return nil
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

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "alb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
