package waf_cdn_domain

import (
	"encoding/json"
	"errors"
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusWafDomainService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafCdnDomainService(c *bp.SdkClient) *ByteplusWafDomainService {
	return &ByteplusWafDomainService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafDomainService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafDomainService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "Page", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListDomain"

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
		logger.Debug(logger.RespFormat, "Result.Data is", resp)
		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		logger.Debug(logger.RespFormat, "Result.Data is", results)

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

func (s *ByteplusWafDomainService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Domain":        id,
		"AccurateQuery": 1,
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
		return data, fmt.Errorf("waf_cdn_domain %s not exist ", id)
	}

	if data["DefenceMode"] != nil {
		data["DefenceModeComputed"] = data["DefenceMode"]
		delete(data, "DefenceMode")
	}

	return data, err
}

func (s *ByteplusWafDomainService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "3")

			if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				d, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}

				status, err = bp.ObtainSdkValue("Status", d)
				logger.Debug(logger.ReqFormat, "waf domain status is %s", status)

				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("get sdk status value error %s", err))
				}
				statusInt, ok := status.(float64)
				if !ok {
					return resource.NonRetryableError(fmt.Errorf("status is not int type %s", status))
				}
				statusString := strconv.Itoa(int(statusInt))

				for _, v := range failStates {
					if v == statusString {
						logger.Debug(logger.ReqFormat, "waf domain statusString is %s", statusString)
						return resource.NonRetryableError(fmt.Errorf("waf domain status error, status: %s", statusString))
					}
				}

				if statusString == "2" || statusString == "5" {
					return resource.RetryableError(fmt.Errorf("waf domain status is %s, retry", statusString))
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			return d, "status by retry", nil
		},
	}
}

func (s *ByteplusWafDomainService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVolcWafServicesByBytePlusCDN",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"tls_enable": {
					TargetField: "TLSEnable",
				},
				"tls_fields_config": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "TLSFieldsConfig",
					NextLevelConvert: map[string]bp.RequestConvert{
						"headers_config": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "HeadersConfig",
							NextLevelConvert: map[string]bp.RequestConvert{
								"enable": {
									TargetField: "Enable",
								},
								"excluded_key_list": {
									TargetField: "ExcludedKeyList",
								},
								"statistical_key_list": {
									TargetField: "StatisticalKeyList",
								},
							},
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Region"] = client.Region
				tLSEnable, ok := d.Get("tls_enable").(int)
				if !ok {
					return false, errors.New("TLSEnable is not int")
				}
				if tLSEnable == 0 {
					(*call.SdkParam)["TLSEnable"] = 0
				}
				domainString, ok := d.Get("domain").(string)
				if !ok {
					return false, errors.New("domain is not string")
				}
				(*call.SdkParam)["Domains"] = []string{domainString}
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				logger.Debug(logger.RespFormat, "before execute", "wxy-test")
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, "id is before", resp)
				id, _ := bp.ObtainSdkValue("Result.DomainList", *resp)
				logger.Debug(logger.RespFormat, "id is", id)
				domainId, ok := id.([]interface{})[0].(string)
				if !ok {
					return errors.New("id is not string")
				}
				logger.Debug(logger.RespFormat, "domainId is", domainId)
				d.SetId(domainId)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"status by retry"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafDomainService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafDomainService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) (callbacks []bp.Callback) {
	if resourceData.HasChanges("bot_repeat_enable", "bot_dytoken_enable", "auto_cc_enable", "bot_sequence_enable",
		"bot_sequence_default_action", "bot_frequency_enable", "waf_enable", "cc_enable", "white_enable",
		"black_ip_enable", "black_lct_enable", "waf_white_req_enable", "white_field_enable", "custom_rsp_enable",
		"system_bot_enable", "custom_bot_enable", "api_enable", "tamper_proof_enable", "dlp_enable") {
		modifyUpdateWafServiceControl := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateWafServiceControl",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"auto_cc_enable": {
						TargetField: "AutoCCEnable",
					},
					"tls_enable": {
						TargetField: "TLSEnable",
					},
					"bot_repeat_enable": {
						TargetField: "BotRepeatEnable",
					},
					"bot_dytoken_enable": {
						TargetField: "BotDytokenEnable",
					},
					"bot_sequence_enable": {
						TargetField: "BotSequenceEnable",
					},
					"bot_sequence_default_action": {
						TargetField: "BotSequenceDefaultAction",
					},
					"bot_frequency_enable": {
						TargetField: "BotFrequencyEnable",
					},
					"waf_enable": {
						TargetField: "WafEnable",
					},
					"cc_enable": {
						TargetField: "CcEnable",
					},
					"white_enable": {
						TargetField: "WhiteEnable",
					},
					"black_ip_enable": {
						TargetField: "BlackIpEnable",
					},
					"black_lct_enable": {
						TargetField: "BlackLctEnable",
					},
					"waf_white_req_enable": {
						TargetField: "WafWhiteReqEnable",
					},
					"white_field_enable": {
						TargetField: "WhiteFieldEnable",
					},
					"custom_rsp_enable": {
						TargetField: "CustomRspEnable",
					},
					"system_bot_enable": {
						TargetField: "SystemBotEnable",
					},
					"custom_bot_enable": {
						TargetField: "CustomBotEnable",
					},
					"api_enable": {
						TargetField: "ApiEnable",
					},
					"tamper_proof_enable": {
						TargetField: "TamperProofEnable",
					},
					"dlp_enable": {
						TargetField: "DlpEnable",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["Region"] = client.Region
					(*call.SdkParam)["Host"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("UpdateWafServiceControl"), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"status by retry"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, modifyUpdateWafServiceControl)
	}

	if resourceData.HasChanges("extra_defence_mode_lb_instance", "defence_mode") {

		modifyServiceDefenceMode := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyServiceDefenceMode",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"extra_defence_mode_lb_instance": {
						ConvertType: bp.ConvertJsonObjectArray,
						TargetField: "ExtraDefenceModeLBInstance",
						NextLevelConvert: map[string]bp.RequestConvert{
							"defence_mode": {
								TargetField: "DefenceMode",
							},
							"instance_id": {
								TargetField: "InstanceID",
							},
						},
					},
					"defence_mode": {
						TargetField: "DefenceMode",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["Host"] = d.Id()
					(*call.SdkParam)["Region"] = client.Region
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("ModifyServiceDefenceMode"), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"status by retry"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, modifyServiceDefenceMode)
	}
	return callbacks
}

func (s *ByteplusWafDomainService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVolcWafService",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Host": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusWafDomainService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Ids",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "Name",
		IdField:      "Id",
		CollectField: "instances",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Id": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusWafDomainService) ReadResourceId(id string) string {
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

func (s *ByteplusWafDomainService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "waf",
		ResourceType:         "domain",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
