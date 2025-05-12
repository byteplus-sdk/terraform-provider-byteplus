package alb_server_group

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/vpc"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbServerGroupService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbServerGroupService(c *bp.SdkClient) *ByteplusAlbServerGroupService {
	return &ByteplusAlbServerGroupService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbServerGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbServerGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeServerGroups"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.ServerGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ServerGroups is not Slice")
		}
		return data, err
	})
	if err != nil {
		return data, err
	}

	for _, value := range data {
		serverGroup, ok := value.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("Server group is not map ")
		}

		detailAction := "DescribeServerGroupAttributes"
		req := map[string]interface{}{
			"ServerGroupId": serverGroup["ServerGroupId"],
		}
		logger.Debug(logger.ReqFormat, detailAction, req)
		detailResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(detailAction), &req)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, detailAction, *detailResp)

		servers, err := bp.ObtainSdkValue("Result.Servers", *detailResp)
		if err != nil {
			return data, err
		}
		serverGroup["Servers"] = servers
	}

	return data, err
}

func (s *ByteplusAlbServerGroupService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ServerGroupIds.1": id,
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
		return data, fmt.Errorf("alb_server_group %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbServerGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("alb_server_group status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (ByteplusAlbServerGroupService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"URI": {
				TargetField: "uri",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbServerGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateServerGroup",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"health_check": {
					TargetField: "HealthCheck",
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"uri": {
							TargetField: "URI",
						},
					},
				},
				"sticky_session_config": {
					TargetField: "StickySessionConfig",
					ConvertType: bp.ConvertListUnique,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ServerGroupId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbServerGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyServerGroupAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"server_group_name": {
					TargetField: "ServerGroupName",
				},
				"description": {
					TargetField: "Description",
				},
				"scheduler": {
					TargetField: "Scheduler",
				},
				"health_check": {
					TargetField: "HealthCheck",
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"uri": {
							TargetField: "URI",
						},
					},
				},
				"sticky_session_config": {
					TargetField: "StickySessionConfig",
					ConvertType: bp.ConvertListUnique,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					if enabled := d.Get("sticky_session_config.0.sticky_session_enabled").(string); enabled == "on" {
						(*call.SdkParam)["StickySessionConfig.StickySessionEnabled"] = enabled
						(*call.SdkParam)["StickySessionConfig.StickySessionType"] = d.Get("sticky_session_config.0.sticky_session_type").(string)
						(*call.SdkParam)["StickySessionConfig.Cookie"] = d.Get("sticky_session_config.0.cookie").(string)
						(*call.SdkParam)["StickySessionConfig.CookieTimeout"] = d.Get("sticky_session_config.0.cookie_timeout").(int)
					}
					(*call.SdkParam)["ServerGroupId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbServerGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteServerGroup",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"ServerGroupId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading alb server group on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusAlbServerGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ServerGroupIds",
				ConvertType: bp.ConvertWithN,
			},
			"server_group_names": {
				TargetField: "ServerGroupNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ServerGroupName",
		IdField:      "ServerGroupId",
		CollectField: "server_groups",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ServerGroupId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"URI": {
				TargetField: "uri",
			},
		},
	}
}

func (s *ByteplusAlbServerGroupService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusAlbServerGroupService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "alb",
		ResourceType:         "servergroup",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
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
