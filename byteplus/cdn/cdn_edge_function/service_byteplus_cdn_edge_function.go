package cdn_edge_function

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnEdgeFunctionService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnEdgeFunctionService(c *bp.SdkClient) *ByteplusCdnEdgeFunctionService {
	return &ByteplusCdnEdgeFunctionService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnEdgeFunctionService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnEdgeFunctionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageOffsetQuery(m, "Limit", "Page", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListSparrow"

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

		results, err = bp.ObtainSdkValue("Result.Sparrows", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Sparrows is not Slice")
		}
		return data, err
	})

	for _, v := range data {
		function, ok := v.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf(" The Sparrow of Result is not map ")
		}
		functionId := function["FunctionId"]

		// 查询函数最新版本代码
		codeAction := "GetSourceCode"
		codeReq := map[string]interface{}{
			"FunctionId": functionId,
		}
		logger.Debug(logger.ReqFormat, codeAction, codeReq)
		codeResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(codeAction), &codeReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, codeAction, codeResp)
		sourceCode, err := bp.ObtainSdkValue("Result.SourceCode", *codeResp)
		if err != nil {
			return data, err
		}
		function["SourceCode"] = sourceCode

		// 查询函数环境变量
		envAction := "GetEnv"
		envReq := map[string]interface{}{
			"FunctionId": functionId,
		}
		logger.Debug(logger.ReqFormat, envAction, envReq)
		envResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(envAction), &envReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, envAction, envResp)
		envs, err := bp.ObtainSdkValue("Result.Envs", *envResp)
		if err != nil {
			return data, err
		}
		function["Envs"] = envs

		// 查询函数灰度发布集群
		canaryAction := "ListContinentCluster"
		canaryReq := map[string]interface{}{
			"FunctionId":  functionId,
			"ClusterType": 100,
		}
		logger.Debug(logger.ReqFormat, canaryAction, canaryReq)
		canaryResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(canaryAction), &canaryReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, canaryAction, canaryResp)
		clusters, err := bp.ObtainSdkValue("Result.ContinentCluster", *canaryResp)
		if err != nil {
			return data, err
		}
		function["ContinentCluster"] = clusters
	}

	return data, err
}

func (s *ByteplusCdnEdgeFunctionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		result interface{}
		ok     bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetSparrow"
	req := map[string]interface{}{
		"FunctionId": id,
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result.Sparrow", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = result.(map[string]interface{}); !ok {
		return data, errors.New("Value is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("edge_function %s is not exist ", id)
	}

	// 查询函数最新版本代码
	codeAction := "GetSourceCode"
	codeReq := map[string]interface{}{
		"FunctionId": id,
	}
	logger.Debug(logger.ReqFormat, codeAction, codeReq)
	codeResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(codeAction), &codeReq)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, codeAction, codeResp)
	sourceCode, err := bp.ObtainSdkValue("Result.SourceCode", *codeResp)
	if err != nil {
		return data, err
	}
	data["SourceCode"] = sourceCode

	// 查询函数环境变量
	envAction := "GetEnv"
	envReq := map[string]interface{}{
		"FunctionId": id,
	}
	logger.Debug(logger.ReqFormat, envAction, envReq)
	envResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(envAction), &envReq)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, envAction, envResp)
	envs, err := bp.ObtainSdkValue("Result.Envs", *envResp)
	if err != nil {
		return data, err
	}
	data["Envs"] = envs

	// 查询函数灰度发布集群
	canaryAction := "ListContinentCluster"
	canaryReq := map[string]interface{}{
		"FunctionId":  id,
		"ClusterType": 100,
	}
	logger.Debug(logger.ReqFormat, canaryAction, canaryReq)
	canaryResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(canaryAction), &canaryReq)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, canaryAction, canaryResp)
	clusters, err := bp.ObtainSdkValue("Result.ContinentCluster", *canaryResp)
	if err != nil {
		return data, err
	}
	clusterArr, ok := clusters.([]interface{})
	if !ok {
		return data, fmt.Errorf("Result.ContinentCluster is not slice")
	}
	var countries []interface{}
	for _, v := range clusterArr {
		cluster, ok := v.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("Result.ContinentCluster value is not map")
		}
		countries = append(countries, cluster["Country"])
	}
	data["CanaryCountries"] = countries

	return data, err
}

func (s *ByteplusCdnEdgeFunctionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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

			// 添加重试操作
			if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				d, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) || bp.AccessDeniedError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}

			return d, strconv.Itoa(int(status.(float64))), err
		},
	}
}

func (ByteplusCdnEdgeFunctionService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnEdgeFunctionService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateSparrow",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"name": {
					TargetField: "Name",
				},
				"remark": {
					TargetField: "Remark",
				},
				"project_name": {
					TargetField: "ProjectName",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.FunctionId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"100", "400", "500"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	if _, exist := resourceData.GetOk("source_code"); exist {
		codeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateSourceCode",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"source_code": {
						TargetField: "SourceCode",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["FunctionId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"100", "400", "500"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, codeCallback)
	}

	if _, exist := resourceData.GetOk("envs"); exist {
		codeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "AddEnv",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"envs": {
						TargetField: "Envs",
						ConvertType: bp.ConvertJsonObjectArray,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["FunctionId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"100", "400", "500"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, codeCallback)
	}

	if _, exist := resourceData.GetOk("canary_countries"); exist {
		codeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateCountryCluster",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"canary_countries": {
						TargetField: "Countries",
						ConvertType: bp.ConvertJsonArray,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["FunctionId"] = d.Id()
					(*call.SdkParam)["ClusterType"] = 100
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"100", "400", "500"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, codeCallback)
	}

	return callbacks
}

func (s *ByteplusCdnEdgeFunctionService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if resourceData.HasChanges("name", "remark") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateSparrow",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"name": {
						TargetField: "Name",
						ForceGet:    true,
					},
					"remark": {
						TargetField: "Remark",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["FunctionId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"100", "400", "500"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("source_code") {
		codeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateSourceCode",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"source_code": {
						TargetField: "SourceCode",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["FunctionId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"100", "400", "500"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, codeCallback)
	}

	if resourceData.HasChange("envs") {
		callbacks = s.updateEnvs(resourceData, callbacks)
	}

	if resourceData.HasChange("canary_countries") {
		callbacks = s.updateCanaryCountries(resourceData, callbacks)
	}

	return callbacks
}

func (s *ByteplusCdnEdgeFunctionService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSparrow",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"FunctionId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading edge function on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusCdnEdgeFunctionService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "Name",
		IdField:      "FunctionId",
		CollectField: "edge_functions",
		ResponseConverts: map[string]bp.ResponseConvert{
			"FunctionId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusCdnEdgeFunctionService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusCdnEdgeFunctionService) updateEnvs(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
	addedEnvs, removedEnvs, _, _ := bp.GetSetDifference("envs", resourceData, envsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteEnv",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedEnvs != nil && len(removedEnvs.List()) > 0 {
					(*call.SdkParam)["FunctionId"] = resourceData.Id()
					(*call.SdkParam)["EnvKeys"] = make([]string, 0)
					for _, v := range removedEnvs.List() {
						env, ok := v.(map[string]interface{})
						if !ok {
							return false, fmt.Errorf(" The env is not map ")
						}
						(*call.SdkParam)["EnvKeys"] = append((*call.SdkParam)["EnvKeys"].([]string), env["key"].(string))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddEnv",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedEnvs != nil && len(addedEnvs.List()) > 0 {
					(*call.SdkParam)["FunctionId"] = resourceData.Id()
					(*call.SdkParam)["Envs"] = make([]map[string]interface{}, 0)
					for _, v := range addedEnvs.List() {
						env, ok := v.(map[string]interface{})
						if !ok {
							return false, fmt.Errorf(" The env is not map ")
						}
						(*call.SdkParam)["Envs"] = append((*call.SdkParam)["Envs"].([]map[string]interface{}), env)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
}

func (s *ByteplusCdnEdgeFunctionService) updateCanaryCountries(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
	addedCountries, removedCountries, _, _ := bp.GetSetDifference("canary_countries", resourceData, schema.HashString, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCountryCluster",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedCountries != nil && len(removedCountries.List()) > 0 {
					(*call.SdkParam)["FunctionId"] = resourceData.Id()
					(*call.SdkParam)["ClusterType"] = 200
					(*call.SdkParam)["Countries"] = removedCountries.List()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCountryCluster",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedCountries != nil && len(addedCountries.List()) > 0 {
					(*call.SdkParam)["FunctionId"] = resourceData.Id()
					(*call.SdkParam)["ClusterType"] = 100
					(*call.SdkParam)["Countries"] = addedCountries.List()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
}

func (s *ByteplusCdnEdgeFunctionService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "CDN",
		ResourceType:         "function",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func getPostUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
