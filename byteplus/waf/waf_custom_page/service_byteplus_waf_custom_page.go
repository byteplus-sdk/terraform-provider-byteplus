package waf_custom_page

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusWafCustomPageService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafCustomPageService(c *bp.SdkClient) *ByteplusWafCustomPageService {
	return &ByteplusWafCustomPageService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafCustomPageService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafCustomPageService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "Page", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListCustomPage"

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

		for _, ele := range data {
			customPage, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" customPage is not Map ")
			}

			customPage["CustomPageId"] = strconv.Itoa(int(customPage["Id"].(float64)))

			logger.Debug(logger.ReqFormat, "CustomPageId", customPage["CustomPageId"])

		}

		return data, err
	})
}

func (s *ByteplusWafCustomPageService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return data, fmt.Errorf("format of waf custom page resource id is invalid,%s", id)
	}
	customPageId := parts[0]
	host := parts[1]

	customPageIdInt, err := strconv.Atoi(customPageId)
	tag := fmt.Sprintf("%012d", customPageIdInt)
	ruleTag := "D" + tag

	req := map[string]interface{}{
		"RuleTag": ruleTag,
		"Host":    host,
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
		return data, fmt.Errorf("waf_custom_page %s not exist ", id)
	}

	if code, codeExist := data["Code"]; codeExist {
		codeString, ok := code.(string)
		if !ok {
			return data, fmt.Errorf("code is not string")
		}
		codeInt, err := strconv.Atoi(codeString)
		if err != nil {
			return data, fmt.Errorf("code can not to int")
		}
		data["Code"] = codeInt
	}

	return data, err
}

func (s *ByteplusWafCustomPageService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafCustomPageService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCustomPage",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"accurate": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Accurate",
					NextLevelConvert: map[string]bp.RequestConvert{
						"accurate_rules": {
							ConvertType: bp.ConvertJsonObjectArray,
							TargetField: "AccurateRules",
							NextLevelConvert: map[string]bp.RequestConvert{
								"http_obj": {
									TargetField: "HttpObj",
								},
								"obj_type": {
									TargetField: "ObjType",
								},
								"opretar": {
									TargetField: "Opretar",
								},
								"property": {
									TargetField: "Property",
								},
								"value_string": {
									TargetField: "ValueString",
								},
							},
						},
						"logic": {
							TargetField: "Logic",
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				enable, ok := d.Get("enable").(int)
				if !ok {
					return false, errors.New("enable is not int")
				}
				if enable == 0 {
					(*call.SdkParam)["Enable"] = 0
				}

				pageMode, ok := d.Get("page_mode").(int)
				if !ok {
					return false, errors.New("pageMode is not int")
				}
				if pageMode == 0 {
					(*call.SdkParam)["PageMode"] = 0
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				host, ok := d.Get("host").(string)
				if !ok {
					return errors.New("host is not string")
				}
				d.SetId(fmt.Sprintf("%s:%s", strconv.Itoa(int(id.(float64))), host))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafCustomPageService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafCustomPageService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCustomPage",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"policy": {
					TargetField: "Policy",
					ForceGet:    true,
				},
				"redirect_url": {
					TargetField: "RedirectUrl",
					ForceGet:    true,
				},
				"advanced": {
					TargetField: "Advanced",
					ForceGet:    true,
				},
				"body": {
					TargetField: "Body",
					ForceGet:    true,
				},
				"page_mode": {
					TargetField: "PageMode",
					ForceGet:    true,
				},
				"code": {
					TargetField: "Code",
					ForceGet:    true,
				},
				"enable": {
					TargetField: "Enable",
					ForceGet:    true,
				},
				"url": {
					TargetField: "Url",
					ForceGet:    true,
				},
				"description": {
					TargetField: "Description",
					ForceGet:    true,
				},
				"name": {
					TargetField: "Name",
					ForceGet:    true,
				},
				"client_ip": {
					TargetField: "ClientIp",
					ForceGet:    true,
				},
				"project_name": {
					TargetField: "ProjectName",
					ForceGet:    true,
				},
				"accurate": {
					ConvertType: bp.ConvertJsonObject,
					ForceGet:    true,
					TargetField: "Accurate",
					NextLevelConvert: map[string]bp.RequestConvert{
						"accurate_rules": {
							ConvertType: bp.ConvertJsonObjectArray,
							ForceGet:    true,
							TargetField: "AccurateRules",
							NextLevelConvert: map[string]bp.RequestConvert{
								"http_obj": {
									TargetField: "HttpObj",
									ForceGet:    true,
								},
								"obj_type": {
									TargetField: "ObjType",
									ForceGet:    true,
								},
								"opretar": {
									TargetField: "Opretar",
									ForceGet:    true,
								},
								"property": {
									TargetField: "Property",
									ForceGet:    true,
								},
								"value_string": {
									TargetField: "ValueString",
									ForceGet:    true,
								},
							},
						},
						"logic": {
							TargetField: "Logic",
							ForceGet:    true,
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf custom page resource id is invalid,%s", d.Id())
				}
				id := parts[0]
				host := parts[1]
				customPageId, err := strconv.Atoi(id)
				if err != nil {
					return false, fmt.Errorf(" custom page id cannot convert to int ")
				}
				(*call.SdkParam)["Host"] = host
				(*call.SdkParam)["Id"] = customPageId
				logic, ok := d.Get("accurate.0.logic").(int)
				if !ok {
					return false, fmt.Errorf("accurate.0.logic cannot convert to int ")
				}

				if logic == 0 {
					delete(*call.SdkParam, "Accurate.Logic")
				}
				enable, ok := d.Get("enable").(int)
				if !ok {
					return false, errors.New("enable is not int")
				}
				if enable == 0 {
					(*call.SdkParam)["Enable"] = 0
				}

				pageMode, ok := d.Get("page_mode").(int)
				if !ok {
					return false, errors.New("pageMode is not int")
				}
				if pageMode == 0 {
					(*call.SdkParam)["PageMode"] = 0
				}
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

func (s *ByteplusWafCustomPageService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCustomPage",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("format of waf custom page resource id is invalid,%s", d.Id())
				}
				id := parts[0]
				host := parts[1]
				customPageId, err := strconv.Atoi(id)
				if err != nil {
					return false, fmt.Errorf(" custom page id cannot convert to int ")
				}
				(*call.SdkParam)["Id"] = customPageId
				(*call.SdkParam)["Host"] = host
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
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

func (s *ByteplusWafCustomPageService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:        "Name",
		IdField:          "CustomPageId",
		CollectField:     "data",
		ContentType:      bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{},
	}
}

func (s *ByteplusWafCustomPageService) ReadResourceId(id string) string {
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
