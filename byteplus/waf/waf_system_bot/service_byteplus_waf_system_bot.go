package waf_system_bot

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

type ByteplusWafSystemBotService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafSystemBotService(c *bp.SdkClient) *ByteplusWafSystemBotService {
	return &ByteplusWafSystemBotService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafSystemBotService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafSystemBotService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListSystemBotConfig"

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
		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Data is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusWafSystemBotService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
		result  map[string]interface{}
	)
	req := map[string]interface{}{
		"Host": resourceData.Get("host"),
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if data["BotType"] == resourceData.Get("bot_type") {
			result = data
			break
		}
	}
	if len(result) == 0 {
		return result, fmt.Errorf("waf_system_bot %s not exist ", id)
	}

	return result, err
}

func (s *ByteplusWafSystemBotService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafSystemBotService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateSystemBotConfig",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%s", d.Get("bot_type"), d.Get("host")))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafSystemBotService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafSystemBotService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateSystemBotConfig",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"project_name": {
					TargetField: "ProjectName",
					ForceGet:    true,
				},
				"action": {
					TargetField: "Action",
					ForceGet:    true,
				},
				"enable": {
					TargetField: "Enable",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf system bot resource id is invalid,%s", d.Id())
				}
				botType := parts[0]
				host := parts[1]
				(*call.SdkParam)["BotType"] = botType
				(*call.SdkParam)["Host"] = host
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
	return []bp.Callback{callback}
}

func (s *ByteplusWafSystemBotService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateSystemBotConfig",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf system bot resource id is invalid,%s", d.Id())
				}
				botType := parts[0]
				host := parts[1]
				(*call.SdkParam)["BotType"] = botType
				(*call.SdkParam)["Host"] = host
				(*call.SdkParam)["Action"] = d.Get("action")
				(*call.SdkParam)["Enable"] = 0
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return s.checkResourceUtilRemoved(d, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading waf custom page on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusWafSystemBotService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:        "BotType",
		IdField:          "RuleTag",
		CollectField:     "data",
		ContentType:      bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{},
	}
}

func (s *ByteplusWafSystemBotService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "waf",
		Version:     "2023-12-25",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
		RegionType:  bp.Global,
	}
}

func (s *ByteplusWafSystemBotService) checkResourceUtilRemoved(d *schema.ResourceData, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		systemBotConfig, _ := s.ReadResource(d, d.Id())
		logger.Debug(logger.RespFormat, "systemBotConfig", systemBotConfig)

		// 能查询成功代表还在删除中，重试
		systemBotConfigInt, ok := systemBotConfig["Enable"].(float64)
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("enable is not float64"))
		}
		if int(systemBotConfigInt) == 1 {
			return resource.RetryableError(fmt.Errorf("resource still in removing status "))
		} else {
			if int(systemBotConfigInt) == 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("system bot status is not disable "))
			}
		}
	})
}
