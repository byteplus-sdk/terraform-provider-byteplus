package cdn_cipher_template

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnCipherTemplateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnCipherTemplateService(c *bp.SdkClient) *ByteplusCdnCipherTemplateService {
	return &ByteplusCdnCipherTemplateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnCipherTemplateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnCipherTemplateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageNumberQuery(m, "PageSize", "PageNum", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {

		action := "DescribeTemplates"

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

		results, err = bp.ObtainSdkValue("Result.Templates", *resp)
		if err != nil {
			return data, err
		}

		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Templates is not Slice")
		}

		return data, err
	})
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		template, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("template is not a map")
		}
		tmpType := template["Type"].(string)
		if tmpType != "cipher" {
			continue
		}
		action := "DescribeCipherTemplate"
		req := map[string]interface{}{
			"TemplateId": template["TemplateId"],
		}
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)

		https, err := bp.ObtainSdkValue("Result.HTTPS", *resp)
		if err != nil {
			return data, err
		}
		template["HTTPS"] = https

		httpForcedRedirect, err := bp.ObtainSdkValue("Result.HttpForcedRedirect", *resp)
		if err != nil {
			return data, err
		}
		template["HttpForcedRedirect"] = httpForcedRedirect

		quic, err := bp.ObtainSdkValue("Result.Quic", *resp)
		if err != nil {
			return data, err
		}
		template["Quic"] = quic
	}

	return data, nil
}

func (s *ByteplusCdnCipherTemplateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	filter := map[string]interface{}{
		"Fuzzy": true,
		"Name":  "Id",
		"Value": []string{id},
	}
	req := map[string]interface{}{
		"Filters": []interface{}{filter},
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
		return data, fmt.Errorf("cdn_cipher_template %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnCipherTemplateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusCdnCipherTemplateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCipherTemplate",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"https": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "HTTPS",
					NextLevelConvert: map[string]bp.RequestConvert{
						"forced_redirect": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "ForcedRedirect",
						},
						"http2": {
							TargetField: "HTTP2",
						},
						"ocsp": {
							TargetField: "OCSP",
						},
						"tls_version": {
							TargetField: "TlsVersion",
							ConvertType: bp.ConvertJsonArray,
						},
						"hsts": {
							TargetField: "Hsts",
							ConvertType: bp.ConvertJsonObject,
						},
					},
				},
				"http_forced_redirect": {
					TargetField: "HttpForcedRedirect",
					ConvertType: bp.ConvertJsonObject,
				},
				"quic": {
					TargetField: "Quic",
					ConvertType: bp.ConvertJsonObject,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.TemplateId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusCdnCipherTemplateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"OCSP": {
				TargetField: "ocsp",
			},
			"HTTP2": {
				TargetField: "http2",
			},
			"HTTPS": {
				TargetField: "https",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnCipherTemplateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateCipherTemplate",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"https": {
					TargetField: "HTTPS",
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"disable_http": {
							TargetField: "DisableHttp",
						},
						"forced_redirect": {
							TargetField: "ForcedRedirect",
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"enable_forced_redirect": {
									TargetField: "EnableForcedRedirect",
								},
								"status_code": {
									TargetField: "StatusCode",
								},
							},
						},
						"http2": {
							TargetField: "HTTP2",
						},
						"hsts": {
							TargetField: "Hsts",
							ConvertType: bp.ConvertJsonObject,
						},
						"tls_version": {
							TargetField: "TlsVersion",
							ConvertType: bp.ConvertJsonArray,
						},
					},
				},
				"http_forced_redirect": {
					TargetField: "HttpForcedRedirect",
					ConvertType: bp.ConvertJsonObject,
				},
				"message": {
					TargetField: "Message",
				},
				"title": {
					TargetField: "Title",
				},
				"quic": {
					TargetField: "Quic",
					ConvertType: bp.ConvertJsonObject,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["TemplateId"] = d.Id()
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

func (s *ByteplusCdnCipherTemplateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteTemplate",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"TemplateId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading cdn cipher template on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusCdnCipherTemplateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"filters": {
				TargetField: "Filters",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Value",
						ConvertType: bp.ConvertJsonArray,
					},
				},
			},
		},
		NameField:    "Title",
		IdField:      "TemplateId",
		CollectField: "templates",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"OCSP": {
				TargetField: "ocsp",
			},
			"HTTP2": {
				TargetField: "http2",
			},
			"HTTPS": {
				TargetField: "https",
			},
		},
	}
}

func (s *ByteplusCdnCipherTemplateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
