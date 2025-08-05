package allow_list

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusMongoDBAllowListService struct {
	Client *bp.SdkClient
}

func NewMongoDBAllowListService(c *bp.SdkClient) *ByteplusMongoDBAllowListService {
	return &ByteplusMongoDBAllowListService{
		Client: c,
	}
}

func (s *ByteplusMongoDBAllowListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongoDBAllowListService) readAllowListDetails(allowListId string) (allowList interface{}, err error) {
	var (
		resp *map[string]interface{}
		//ok   bool
	)
	action := "DescribeAllowListDetail"
	cond := map[string]interface{}{
		"AllowListId": allowListId,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	if err != nil {
		return allowList, err
	}
	logger.Debug(logger.RespFormat, action, resp)

	allowList, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return allowList, err
	}
	if allowList == nil {
		allowList = map[string]interface{}{}
	}
	return allowList, err
}

func (s *ByteplusMongoDBAllowListService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		IdField:      "AllowListId",
		NameField:    "AllowListName",
		CollectField: "allow_lists",
		ContentType:  bp.ContentTypeJson,
		RequestConverts: map[string]bp.RequestConvert{
			"allow_list_ids": {
				TargetField: "AllowListIds",
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"AllowListIPNum": {
				TargetField: "allow_list_ip_num",
			},
			"VPC": {
				TargetField: "vpc",
			},
		},
	}
}

func (s *ByteplusMongoDBAllowListService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp            *map[string]interface{}
		results         interface{}
		allowListIdsMap = make(map[string]bool)
		exists          bool
	)
	action := "DescribeAllowLists"
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
	results, err = bp.ObtainSdkValue("Result.AllowLists", *resp)
	if err != nil {
		logger.DebugInfo("bp.ObtainSdkValue return :%v", err)
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	allowLists, ok := results.([]interface{})
	if !ok {
		return data, fmt.Errorf("DescribeAllowLists responsed instances is not a slice")
	}

	if _, exists = condition["AllowListIds"]; exists {
		if allowListIds, ok := condition["AllowListIds"].([]interface{}); ok {
			for _, id := range allowListIds {
				allowListIdsMap[id.(string)] = true
			}
		}
	}

	for _, ele := range allowLists {
		allowList := ele.(map[string]interface{})
		id := allowList["AllowListId"].(string)

		// 如果存在 allow_list_ids，过滤掉 allow_list_ids 中未包含的 id
		if _, ok := allowListIdsMap[id]; exists && !ok {
			continue
		}

		detail, err := s.readAllowListDetails(id)
		if err != nil {
			logger.DebugInfo("read allow list %s detail failed,err:%v.", id, err)
			data = append(data, ele)
			continue
		}
		allowList["AllowList"] = detail.(map[string]interface{})["AllowList"]
		allowList["AssociatedInstances"] = detail.(map[string]interface{})["AssociatedInstances"]

		data = append(data, allowList)
	}
	return data, nil
}

func (s *ByteplusMongoDBAllowListService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"RegionId":     s.Client.Region,
		"AllowListIds": []interface{}{id},
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
		return data, fmt.Errorf("allowlist %s is not exist", id)
	}
	return data, err
}

func (s *ByteplusMongoDBAllowListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusMongoDBAllowListService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, map[string]bp.ResponseConvert{
			"AllowListIPNum": {
				TargetField: "allow_list_ip_num",
			},
			"VPC": {
				TargetField: "vpc",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBAllowListService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAllowList",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				id, _ := bp.ObtainSdkValue("Result.AllowListId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBAllowListService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyAllowList",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"allow_list_desc": {
					TargetField: "AllowListDesc",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllowListId"] = d.Id()
				(*call.SdkParam)["AllowListName"] = resourceData.Get("allow_list_name")
				if resourceData.HasChange("allow_list") {
					//describe allow list, get instance num
					var applyInstanceNum int
					detail, err := s.readAllowListDetails(d.Id())
					if err != nil {
						return false, fmt.Errorf("read allow list detail faield")
					}
					if associatedInstances, ok := detail.(map[string]interface{})["AssociatedInstances"]; !ok {
						return false, fmt.Errorf("read AssociatedInstances failed")
					} else {
						applyInstanceNum = len(associatedInstances.([]interface{}))
					}
					(*call.SdkParam)["ApplyInstanceNum"] = applyInstanceNum
					allowList, ok := resourceData.GetOk("allow_list")
					if !ok || len(allowList.(string)) == 0 {
						// 对于置空allow list做特殊处理
						old, _ := resourceData.GetChange("allow_list")
						(*call.SdkParam)["AllowList"] = old.(string)
						(*call.SdkParam)["ModifyMode"] = "Delete"
						return true, nil
					}
					(*call.SdkParam)["AllowList"] = resourceData.Get("allow_list")
					(*call.SdkParam)["ModifyMode"] = "Cover"
				}
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

func (s *ByteplusMongoDBAllowListService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAllowList",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllowListId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading mongodb allow list on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusMongoDBAllowListService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusMongoDBAllowListService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "mongodb",
		ResourceType:         "allowlist",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
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
